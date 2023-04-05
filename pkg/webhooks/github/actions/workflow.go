package actions

import (
	"context"
	"errors"
	"fmt"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"github.com/google/go-github/v50/github"
	"go.uber.org/zap"

	"github.tools.sap/actions-rollout-app/pkg/clients"
)

const (
	workflowRunTitle        = "Workflow run %s %s"
	workflowRunMessage      = "Workflow %s completed with status %s"
	workflowJobTitle        = "Workflow job %s %s"
	workflowJobMessage      = "Workflow job %s completed with status %s"
	workflowDispatchTitle   = "Workflow dispatch %s"
	workflowDispatchMessage = "Workflow %s dispatched"
)

type WorkflowActionParams struct {
	Workflow     string
	Organization string
	Repository   string
	WebhookEvent ghwebhooks.Event
	Sender       string
}

type WorkflowAction struct {
	logger *zap.SugaredLogger
	client *clients.Github

	repository     string
	organization   string
	workflow       string
	workerPoolSize float64
	filesPath      *[]string
}

// TODO: retest this

func NewWorkflowAction(logger *zap.SugaredLogger, client *clients.Github, rawConfig map[string]any) (*WorkflowAction, error) {
	// Validate input
	repo, ok := rawConfig["repository"].(string)
	if !ok {
		return nil, errors.New("repo not found or is not a string")
	}
	org, ok := rawConfig["organization"].(string)
	if !ok {
		return nil, errors.New("org not found or is not a string")
	}
	workflow, ok := rawConfig["workflow"].(string)
	if !ok {
		return nil, errors.New("workflow not found or is not a string")
	}
	filesInterface, ok := rawConfig["files_path"].([]interface{})
	if !ok {
		return nil, errors.New("filesPath not found or is not a slice of interface{}")
	}
	files := make([]string, len(filesInterface))
	for i, v := range filesInterface {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("file path at index %d is not a string", i)
		}
		files[i] = str
	}

	workers, ok := rawConfig["workers"].(float64)
	if !ok {
		return nil, errors.New("workers not found or is not an int")
	}

	// Create WorkflowAction object using struct initialization
	return &WorkflowAction{
		logger:         logger,
		client:         client,
		repository:     repo,
		organization:   org,
		workflow:       workflow,
		filesPath:      &files,
		workerPoolSize: workers,
	}, nil
}

func (w *WorkflowAction) HandleWorkflow(ctx context.Context, p *WorkflowActionParams) error {
	if w.repository != p.Repository {
		w.logger.Errorw("repository does not match", "expected", w.repository, "actual", p.Repository)
		return nil
	}

	if w.organization != p.Organization {
		w.logger.Errorw("organization does not match", "expected", w.organization, "actual", p.Organization)
		return nil
	}

	if w.workflow != p.Workflow {
		w.logger.Errorw("workflow id does not match", "expected", w.workflow, "actual", p.Workflow)
		return nil
	}

	handlerMap := map[ghwebhooks.Event]func(context.Context, *WorkflowActionParams) error{
		"workflow_run":      w.handleWorkflowRun,
		"workflow_dispatch": w.handleWorkflowDispatch,
		"workflow_job":      w.handleWorkflowJob,
	}

	handler, ok := handlerMap[p.WebhookEvent]
	if !ok {
		w.logger.Errorw("unknown webhook event type", "event", p.WebhookEvent)
		return nil
	}

	if err := handler(ctx, p); err != nil {
		w.logger.Debugw("failed to handle event", "event", p.WebhookEvent, "error", err)
	}

	return nil
}

// TODO: https://github.com/githubcustomers/SAP/issues/1002#issuecomment-1477526984

func (w *WorkflowAction) createWorkflowIssue(ctx context.Context, title, message string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	issue, issueResp, err := w.client.GetV3Client().Issues.Create(ctx, w.organization, w.repository, &github.IssueRequest{
		Title:    github.String(title),
		Body:     github.String(message),
		Assignee: github.String("mouismail"),
	})
	if err != nil {
		w.logger.Errorw("error creating issue", "error", err)
		return err
	}

	w.logger.Infow("issue created", "issue_id", issue.ID, "response", issueResp.Response.StatusCode)
	return nil
}

func (w *WorkflowAction) handleWorkflowRun(ctx context.Context, p *WorkflowActionParams) error {
	title := fmt.Sprintf("%s/%s - %s", p.Organization, p.Repository, p.Workflow)
	message := fmt.Sprintf("Workflow Run triggered by %s on organization %s and repository %s", p.Sender, p.Organization, p.Repository)
	repoAction := &RepoAction{
		logger:                 w.logger,
		client:                 w.client,
		validationOrganization: w.organization,
		validationRepository:   w.repository,
		filesPath:              w.filesPath,
		workerPoolSize:         w.workerPoolSize,
	}

	repoParams := &RepoActionParams{
		ValidationOrganization: p.Organization,
		ValidationRepository:   p.Repository,
	}
	err := repoAction.HandleRepo(repoParams)
	if err != nil {
		return err
	}
	return w.createWorkflowIssue(ctx, title, message)
}

func (w *WorkflowAction) handleWorkflowDispatch(ctx context.Context, p *WorkflowActionParams) error {
	title := fmt.Sprintf("%s/%s - %s", p.Organization, p.Repository, p.Workflow)
	message := fmt.Sprintf("Workflow dispatch triggered by %s on organization %s and repository %s", p.Sender, p.Organization, p.Repository)

	repoAction := &RepoAction{
		logger:                 w.logger,
		client:                 w.client,
		validationOrganization: w.organization,
		validationRepository:   w.repository,
		filesPath:              w.filesPath,
		workerPoolSize:         w.workerPoolSize,
	}

	repoParams := &RepoActionParams{
		ValidationOrganization: p.Organization,
		ValidationRepository:   p.Repository,
	}
	err := repoAction.HandleRepo(repoParams)
	if err != nil {
		return err
	}
	return w.createWorkflowIssue(ctx, title, message)
}

func (w *WorkflowAction) handleWorkflowJob(ctx context.Context, p *WorkflowActionParams) error {
	title := fmt.Sprintf("%s/%s - %s", p.Organization, p.Repository, p.Workflow)
	message := fmt.Sprintf("%s triggered by %s on organization %s and repository %s", p.WebhookEvent, p.Sender, p.Organization, p.Repository)
	repoAction := &RepoAction{
		logger:                 w.logger,
		client:                 w.client,
		validationOrganization: w.organization,
		validationRepository:   w.repository,
		filesPath:              w.filesPath,
		workerPoolSize:         w.workerPoolSize,
	}

	repoParams := &RepoActionParams{
		ValidationOrganization: p.Organization,
		ValidationRepository:   p.Repository,
	}
	err := repoAction.HandleRepo(repoParams)
	if err != nil {
		return err
	}
	return w.createWorkflowIssue(ctx, title, message)
}
