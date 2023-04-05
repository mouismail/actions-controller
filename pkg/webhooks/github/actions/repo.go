package actions

import (
	"context"
	"errors"
	"fmt"
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
	}, nil
}

func (r *RepoAction) HandleRepo(params *RepoActionParams) error {
	r.logger.Infof("Validating repository %s/%s", params.ValidationOrganization, params.ValidationRepository)
	err := r.handleRepoConfig(params)
	if err != nil {
		return err
	}
	r.logger.Infof("Repository %s/%s is valid, validating...", params.ValidationOrganization, r.filesPath)

	return nil
}

func (r *RepoAction) handleRepoConfig(params *RepoActionParams) error {
	if r.filesPath == nil {
		r.logger.Infof("No files to validate, skipping")
		return nil
	}

	r.logger.Infof("Checking paths")
	if err := r.handleRepoConfigFile(params); err != nil {
		r.logger.Errorf("Error validating repository: %v", err)
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
				r.logger.Infow("Checking  file %s in %s/%s", path, params.ValidationOrganization, params.ValidationRepository)
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
					r.logger.Infow("Processing file %s in %s/%s", file.GetName(), params.ValidationOrganization, params.ValidationRepository)
					err := r.downloadRawData(params, fmt.Sprintf("%s/%s", path, file.GetName()))
					if err != nil {
						r.logger.Errorf("Error downloading file content: %s", err.Error())
						errCh <- err
						continue
					}
					//r.HandleRepoConfigFileContent(params, contentString)
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
	rawContents, resp, err := r.client.GetV3Client().Repositories.DownloadContents(
		context.Background(),
		r.client.Organization(),
		r.client.Repository(),
		filePath,
		&github.RepositoryContentGetOptions{Ref: "main"},
	)
	if err != nil {
		fmt.Printf("Error downloading file content: %s\n", err.Error())
		return err
	}
	if resp.StatusCode != http.StatusOK {

	}
	defer rawContents.Close()

	bytes, err := io.ReadAll(rawContents)
	if err != nil {
		r.logger.Errorw("Error reading file content", "error", err)
		return err
	}

	err = r.handleRepoConfigFileContent(params, bytes)
	if err != nil {
		fmt.Printf("Error validating file content: %s\n", err.Error())
		return err
	}
	return nil
}

func (r *RepoAction) handleRepoConfigFileContent(params *RepoActionParams, content []byte) error {
	if content == nil {
		return errors.New("invalid params")
	}

	var validation ValidatorData

	if err := yaml.Unmarshal(content, &validation); err != nil {
		return err
	}

	if validation.URL != "https://octodemo.com/"+params.ValidationOrganization {
		return errors.New("invalid url")
	}
	if validation.ContactEmail != "" {
		return errors.New("invalid contact email")
	}
	if validation.UseCase != params.ValidationOrganization {
		return errors.New("invalid use case")
	}
	if len(validation.Repos) < 0 {
		return errors.New("invalid repos")
	}
	if len(validation.Repos) != 0 {
		for _, repo := range validation.Repos {
			if repo != params.ValidationRepository {
				return errors.New("invalid repo")
			}
		}
	}
	r.logger.Infof("Repository %s/%s is valid", params.ValidationOrganization, params.ValidationRepository)

	return nil
}
