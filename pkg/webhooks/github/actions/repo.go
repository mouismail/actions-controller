package actions

import (
	"context"
	"errors"
	"fmt"

	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/go-github/v50/github"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/utils"
)

type RepoActionParams struct {
	ValidationOrganization string
	ValidationRepository   string
	ConfigFileName         string
	FilesPath              *[]string
}

type RepoAction struct {
	logger *zap.SugaredLogger
	client *clients.Github

	validationOrganization string
	validationRepository   string
	filesPath              *[]string
	workerPoolSize         float64
	assignees              *[]string
}

type ValidatorData struct {
	URL          string   `yaml:"url"`
	ContactEmail string   `yaml:"contactEmail"`
	UseCase      string   `yaml:"useCase"`
	Repos        []string `yaml:"repos,omitempty"`
}

func NewRepoAction(logger *zap.SugaredLogger, client *clients.Github, rawConfig map[string]interface{}) (*RepoAction, error) {
	validationOrganization, ok := rawConfig["validationOrganization"].(string)
	if !ok {
		return nil, errors.New("validationOrganization not found or is not a string %w")
	}

	validationRepository, ok := rawConfig["validationRepositories"].(string)
	if !ok {
		return nil, errors.New("validationRepositories not found or is not a string slice")
	}

	return &RepoAction{
		logger:                 logger,
		client:                 client,
		validationOrganization: validationOrganization,
		validationRepository:   validationRepository,
		filesPath:              rawConfig["filesPath"].(*[]string),
		assignees:              rawConfig["assignees"].(*[]string),
	}, nil
}

func (r *RepoAction) HandleRepo(ctx context.Context, params *RepoActionParams) error {
	r.logger.Infof("validating repository %s/%s", params.ValidationOrganization, params.ValidationRepository)
	err := r.handleRepoConfig(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoAction) handleRepoConfig(ctx context.Context, params *RepoActionParams) error {
	if r.filesPath == nil {
		return errors.New("no files to validate")
	}

	r.logger.Infof("checking paths")
	if err := r.handleRepoConfigFile(ctx, params); err != nil {
		return err
	}

	return nil
}

func (r *RepoAction) handleRepoConfigFile(ctx context.Context, params *RepoActionParams) error {
	if params == nil || r.filesPath == nil {
		return errors.New("invalid params")
	}

	filesCh := make(chan string, len(*r.filesPath))
	errCh := make(chan error, len(*r.filesPath))
	isValidCh := make(chan bool, 1)

	for _, path := range *r.filesPath {
		filesCh <- path
	}

	var wg sync.WaitGroup

	contentCache := make(map[string][]*github.RepositoryContent)
	var mu sync.Mutex // Protects access to contentCache
	wg.Add(len(*r.filesPath))

	for _, path := range *r.filesPath {
		go func(path string) {
			defer wg.Done()

			r.logger.Infof("checking  file %s in %s/%s", path, params.ValidationOrganization, params.ValidationRepository)
			mu.Lock()
			content, ok := contentCache[r.client.Organization()+"/"+r.client.Repository()+"/"+path]
			mu.Unlock()

			if !ok {
				dirContent, err := r.getContents(ctx, path)
				if err != nil {
					r.logger.Infof("could not get contents for %s/%s on path %s", params.ValidationOrganization, params.ValidationRepository, path)
					return
				}
				content = dirContent

				mu.Lock()
				contentCache[path] = content
				mu.Unlock()

			}
			for _, file := range content {
				r.logger.Infof("checking config file %s on path %s for workflow event from %s/%s", file.GetName(), path, params.ValidationOrganization, params.ValidationRepository)

				select {
				case <-ctx.Done():
					return
				default:
				}

				isValid, err := r.isValidFile(ctx, params, path, file)
				if isValid {
					isValidCh <- true
					return
				}
				if err != nil {
					r.logger.Infof("could not validate %s/%s for file %s/%s", params.ValidationOrganization, params.ValidationRepository, path, file.GetName())
					continue
				}
			}

		}(path)
	}

	go func() {
		wg.Wait()
		select {
		case isValid := <-isValidCh:
			if !isValid {
				r.logger.Infof("repository %s/%s is not valid", params.ValidationOrganization, params.ValidationRepository)
				errCh <- fmt.Errorf("repository %s/%s is not valid", params.ValidationOrganization, params.ValidationRepository)
				return
			}
		default:
			errCh <- fmt.Errorf("no valid files found in repository %s/%s", params.ValidationOrganization, params.ValidationRepository)
			return
		}
		close(errCh)
		close(isValidCh)
		close(filesCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoAction) downloadRawData(ctx context.Context, params *RepoActionParams, filePath string) (bool, error) {
	rawContents, _, err := r.client.GetV3Client().Repositories.DownloadContents(
		ctx,
		r.client.Organization(),
		r.client.Repository(),
		filePath,
		&github.RepositoryContentGetOptions{Ref: "main"},
	)
	if err != nil {
		r.logger.Errorw("Error downloading the raw content", "error", err)
		return false, err
	}

	defer func(rawContents io.ReadCloser) {
		rawErr := rawContents.Close()
		if rawErr != nil {
			r.logger.Errorw("Error closing raw contents", "error", rawErr)
		}
	}(rawContents)

	bytes, err := io.ReadAll(rawContents)
	if err != nil {
		return false, err
	}

	err = r.handleRepoConfigFileContent(params, bytes)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *RepoAction) handleRepoConfigFileContent(params *RepoActionParams, content []byte) error {
	// TODO: replace
	if content == nil {
		return errors.New(utils.ErrValidationEmptyContent)
	}

	var validation ValidatorData

	if err := yaml.Unmarshal(content, &validation); err != nil {
		return err
	}

	if !(validation.URL == fmt.Sprintf("%s/%s", r.client.ServerInfo().EnterpriseURL, params.ValidationOrganization)) {
		r.logger.Warnw(utils.ErrInvalidConfigOrganization, "URL", validation.URL, "expected", fmt.Sprintf("%s/%s", r.client.ServerInfo().EnterpriseURL, params.ValidationOrganization))
		return fmt.Errorf("%s: %s", utils.ErrInvalidConfigOrganization, validation.URL)
	}
	if validation.ContactEmail == "" {
		r.logger.Warnw(utils.ErrInvalidContactEmail, "ContactEmail", validation.ContactEmail)
		return fmt.Errorf("%s: %s", utils.ErrInvalidContactEmail, validation.ContactEmail)
	}
	if validation.UseCase != params.ValidationOrganization {
		r.logger.Warnw(utils.ErrInvalidUseCase, "UseCase", validation.UseCase, "expected", params.ValidationOrganization)
		return fmt.Errorf("%s: %s", utils.ErrInvalidUseCase, validation.UseCase)
	}
	if len(validation.Repos) < 0 {
		r.logger.Warnw(utils.ErrInvalidConfigRepository, "Repos", validation.Repos)
		return errors.New(utils.ErrInvalidConfigRepository)
	}
	if len(validation.Repos) != 0 {
		for _, repo := range validation.Repos {
			if repo != fmt.Sprintf("%s/%s/%s", r.client.ServerInfo().EnterpriseURL, params.ValidationOrganization, params.ValidationRepository) {
				r.logger.Warnw(utils.ErrInvalidConfigRepository, "Repos", repo, "expected", fmt.Sprintf("%s/%s/%s", r.client.ServerInfo().EnterpriseURL, params.ValidationOrganization, params.ValidationRepository))
				return fmt.Errorf("%s: %s", utils.ErrInvalidConfigRepository, repo)
			}
		}
	}
	r.logger.Infof("Repository %s/%s is valid", params.ValidationOrganization, params.ValidationRepository)

	return nil
}

func (r *RepoAction) isValidFile(ctx context.Context, params *RepoActionParams, path string, file *github.RepositoryContent) (bool, error) {
	if !r.isFileValid(file) {
		return false, nil
	}

	return r.downloadRawData(ctx, params, fmt.Sprintf("%s/%s", path, file.GetName()))
}

func (r *RepoAction) isFileValid(file *github.RepositoryContent) bool {
	return file.GetType() == "file" && strings.HasSuffix(file.GetName(), ".yml")
}

func (r *RepoAction) getContents(ctx context.Context, path string) ([]*github.RepositoryContent, error) {
	_, dirContent, resp, err := r.client.GetV3Client().Repositories.GetContents(ctx, r.client.Organization(), r.client.Repository(), path, &github.RepositoryContentGetOptions{Ref: "main"})
	if err != nil {
		r.logger.Errorf("Error retrieving content for %s/%d/%s: %v", r.client.Organization(), r.filesPath, path, err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		r.logger.Errorf("Unexpected status code %d while retrieving content for %s/%d/%s", resp.StatusCode, r.client.Organization(), r.filesPath, path)
		return nil, err
	}

	return dirContent, nil
}

func (r *RepoAction) GetDisableType() (string, error) {
	return "", nil
}
