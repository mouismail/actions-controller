package actions

import (
	"context"
	"errors"
	"fmt"
	"github.tools.sap/actions-rollout-app/pkg/utils"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"github.com/google/go-github/v50/github"
	"go.uber.org/zap"

	"github.tools.sap/actions-rollout-app/pkg/clients"
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
	assignees      *[]string
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
	assigneesInterface, ok := rawConfig["issue_assignees"].([]interface{})
	if !ok {
		return nil, errors.New("assignees not found or is not a slice of interface{}")
	}
	assignees := make([]string, len(assigneesInterface))
	for i, v := range assigneesInterface {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("assignee at index %d is not a string", i)
		}
		assignees[i] = str
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
		assignees:      &assignees,
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

func (w *WorkflowAction) createWorkflowIssue(ctx context.Context, title, message string, assignees, labels []string) error {
	if ctx == nil {
		ctx = context.Background()
	}

	issue, issueResp, err := w.client.GetV3Client().Issues.Create(ctx, w.organization, w.repository, &github.IssueRequest{
		Title:     github.String(title),
		Body:      github.String(message),
		Assignees: &assignees,
		Labels:    &[]string{w.organization, w.repository, w.workflow},
	})
	if err != nil {
		w.logger.Errorw("error creating issue", "error", err)
		return err
	}

	w.logger.Infow("issue created", "issue_id", issue.ID, "response", issueResp.Response.StatusCode)
	return nil
}

func (w *WorkflowAction) handleWorkflowRun(ctx context.Context, p *WorkflowActionParams) error {
	err := w.handleWorkflowEvent(ctx, p, "run")
	if err != nil {
		return err
	}
	return nil
}

func (w *WorkflowAction) handleWorkflowDispatch(ctx context.Context, p *WorkflowActionParams) error {
	// TODO: to be validated with @stoe
	//err := w.handleWorkflowEvent(ctx, p, "dispatch")
	//if err != nil {
	//	return err
	//}
	return nil
}

func (w *WorkflowAction) handleWorkflowJob(ctx context.Context, p *WorkflowActionParams) error {
	err := w.handleWorkflowEvent(ctx, p, "job")
	if err != nil {
		return err
	}
	return nil
}

func (w *WorkflowAction) handleWorkflowEvent(ctx context.Context, p *WorkflowActionParams, eventType string) error {
	title := fmt.Sprintf("%s/%s - %s", p.Organization, p.Repository, p.Workflow)

	var message string
	switch eventType {
	case "run":
		message = fmt.Sprintf(utils.WorkflowRunMessage, p.Sender, p.Organization, p.Repository)
	case "dispatch":
		message = fmt.Sprintf(utils.WorkflowDispatchMessage, p.Sender, p.Organization, p.Repository)
	case "job":
		message = fmt.Sprintf(utils.WorkflowJobMessage, p.WebhookEvent, p.Sender, p.Organization, p.Repository)
	default:
		return errors.New("unsupported event type")
	}

	repoAction := &RepoAction{
		logger:                 w.logger,
		client:                 w.client,
		validationOrganization: w.organization,
		validationRepository:   w.repository,
		filesPath:              w.filesPath,
		workerPoolSize:         w.workerPoolSize,
		assignees:              w.assignees,
	}

	repoParams := &RepoActionParams{
		ValidationOrganization: p.Organization,
		ValidationRepository:   p.Repository,
	}
	err := repoAction.HandleRepo(repoParams)
	if err != nil {
		return w.createWorkflowIssue(ctx, title, message, *w.assignees, nil)
	}
	return nil
}
