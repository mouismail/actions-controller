package actions

import (
	"context"
	"errors"
	"fmt"
	"github.tools.sap/actions-rollout-app/pkg/utils"
	"io"
	"net/http"
	"sync"

	"github.com/google/go-github/v50/github"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"

	"github.tools.sap/actions-rollout-app/pkg/clients"
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
		workerPoolSize:         rawConfig["workers"].(float64),
		assignees:              rawConfig["assignees"].(*[]string),
	}, nil
}

func (r *RepoAction) HandleRepo(params *RepoActionParams) error {
	r.logger.Infof("Validating repository %s/%s", params.ValidationOrganization, params.ValidationRepository)
	err := r.handleRepoConfig(params)
	if err != nil {
		return err
	}
	r.logger.Infof("Repository %s/%s is valid.", params.ValidationOrganization, params.ValidationRepository)

	return nil
}

func (r *RepoAction) handleRepoConfig(params *RepoActionParams) error {
	if r.filesPath == nil {
		return errors.New("no files to validate")
	}

	r.logger.Infof("Checking paths")
	if err := r.handleRepoConfigFile(params); err != nil {
		return err
	}

	return nil
}

func (r *RepoAction) handleRepoConfigFile(params *RepoActionParams) error {
	// Create a channel to receive files to process
	if params == nil || r.filesPath == nil {
		return errors.New("invalid params")
	}

	filesCh := make(chan string, len(*r.filesPath))
	errCh := make(chan error, len(*r.filesPath))
	for _, path := range *r.filesPath {
		filesCh <- path
	}
	close(filesCh)

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// Create a map to cache the content of the repository files
	contentCache := make(map[string][]*github.RepositoryContent)
	var mu sync.Mutex // Protects access to contentCache

	// Start worker pool
	for i := 0; i < int(r.workerPoolSize); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range filesCh {
				r.logger.Infof("Checking  file %s in %s/%s", path, params.ValidationOrganization, params.ValidationRepository)
				// Check if content is already cached
				mu.Lock()
				content, ok := contentCache[r.client.Organization()+"/"+r.client.Repository()+"/"+path]
				mu.Unlock()

				if !ok {

					// Retrieve content from GitHub API
					_, dirContent, resp, err := r.client.GetV3Client().Repositories.GetContents(context.Background(), r.client.Organization(), r.client.Repository(), path, &github.RepositoryContentGetOptions{Ref: "main"})
					if err != nil {
						r.logger.Errorf("Error retrieving content for %s/%d/%s: %v", params.ValidationOrganization, r.filesPath, path, err)
						errCh <- err
						continue
					}

					if resp.StatusCode != http.StatusOK {
						r.logger.Errorf("Unexpected status code %d while retrieving content for %s/%d/%s", resp.StatusCode, params.ValidationOrganization, r.filesPath, path)
						errCh <- err
						continue
					}
					content = dirContent

					// Cache content
					mu.Lock()
					contentCache[path] = content
					mu.Unlock()

				}
				for _, file := range content {
					//r.logger.Infof("Processing file %s in %s/%s", file.GetName(), params.ValidationOrganization, params.ValidationRepository)
					err := r.downloadRawData(params, fmt.Sprintf("%s/%s", path, file.GetName()))
					if err != nil {
						//r.logger.Errorf("Error downloading file content: %s", err)
						errCh <- err
						continue
					}
				}
			}
		}()
	}

	// Wait for all workers to finish
	//wg.Wait()
	go func() {
		wg.Wait()
		close(errCh) // Close error channel when all workers are done
	}()

	for err := range errCh {
		if err != nil {
			return err // Return the first error encountered
		}
	}
	return nil
}

func (r *RepoAction) downloadRawData(params *RepoActionParams, filePath string) error {
	rawContents, _, err := r.client.GetV3Client().Repositories.DownloadContents(
		context.Background(),
		r.client.Organization(),
		r.client.Repository(),
		filePath,
		&github.RepositoryContentGetOptions{Ref: "main"},
	)
	if err != nil {
		r.logger.Errorw("Error downloading the raw content", "error", err)
		return err
	}

	defer rawContents.Close()

	bytes, err := io.ReadAll(rawContents)
	if err != nil {
		return err
	}

	err = r.handleRepoConfigFileContent(params, bytes)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoAction) handleRepoConfigFileContent(params *RepoActionParams, content []byte) error {
	if content == nil {
		return errors.New(utils.ErrValidationEmptyContent)
	}

	var validation ValidatorData

	if err := yaml.Unmarshal(content, &validation); err != nil {
		return err
	}

	if !(validation.URL == fmt.Sprintf("https://octodemo.com/%s", params.ValidationOrganization)) {
		r.logger.Warnw("Invalid URL", "URL", validation.URL, "expected", fmt.Sprintf("https://octodemo.com/%s", params.ValidationOrganization))
		return errors.New(utils.ErrInvalidConfigOrganization)
	}
	if validation.ContactEmail != "" {
		r.logger.Warnw("Invalid Contact Email", "ContactEmail", validation.ContactEmail)
		return errors.New(utils.ErrInvalidContactEmail)
	}
	if validation.UseCase != params.ValidationOrganization {
		r.logger.Warnw("Invalid Use Case", "UseCase", validation.UseCase, "expected", params.ValidationOrganization)
		return errors.New(utils.ErrInvalidUseCase)
	}
	if len(validation.Repos) < 0 {
		r.logger.Warnw("Invalid Repos/Repo", "Repos", validation.Repos)
		return errors.New(utils.ErrInvalidConfigRepository)
	}
	if len(validation.Repos) != 0 {
		for _, repo := range validation.Repos {
			if repo != params.ValidationRepository {
				r.logger.Warnw("Invalid Repos/Repo", "Repos", validation.Repos, "expected", params.ValidationRepository)
				return errors.New(utils.ErrInvalidConfigRepository)
			}
		}
	}
	r.logger.Infof("Repository %s/%s is valid", params.ValidationOrganization, params.ValidationRepository)

	return nil
}
