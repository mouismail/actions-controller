package actions

import (
	"context"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"go.uber.org/zap"

	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
)

type WorkflowActionParams struct {
	WorkflowId   int64
	Organization string
	Repository   string
	WebhookEvent ghwebhooks.Event
	Sender       string
}

type WorkflowAction struct {
	logger *zap.SugaredLogger
	client *clients.Github

	repository   string
	organization string
	workflowId   int64
}

func NewWorkflowAction(logger *zap.SugaredLogger, client *clients.Github, rawConfig map[string]any) (*WorkflowAction, error) {
	typedConfig := config.WorkflowHandleConfig{
		Repo:       rawConfig["repo"].(string),
		Org:        rawConfig["org"].(string),
		WorkflowId: rawConfig["workflowId"].(int64),
	}

	return &WorkflowAction{
		logger:       logger,
		client:       client,
		repository:   typedConfig.Repo,
		organization: typedConfig.Org,
		workflowId:   typedConfig.WorkflowId,
	}, nil
}

func (w *WorkflowAction) HandleWorkflow(ctx context.Context, p *WorkflowActionParams) error {
	if w.repository != p.Repository {
		w.logger.Debugw("repository does not match", "expected", w.repository, "actual", p.Repository)
		return nil
	}

	if w.organization != p.Organization {
		w.logger.Debugw("organization does not match", "expected", w.organization, "actual", p.Organization)
		return nil
	}

	if w.workflowId != p.WorkflowId {
		w.logger.Debugw("workflow id does not match", "expected", w.workflowId, "actual", p.WorkflowId)
		return nil
	}

	var err error
	if p.WebhookEvent == "workflow_run" {
		w.logger.Debugw("handling workflow run event")
		err = w.handleWorkflowRun(ctx, p)
	} else if p.WebhookEvent == "workflow_dispatch" {
		w.logger.Debugw("handling workflow dispatch event")
		err = w.handleWorkflowDispatch(ctx, p)
	} else if p.WebhookEvent == "workflow_job" {
		w.logger.Debugw("handling workflow job event")
		err = w.handleWorkflowJob(ctx, p)
	}

	if err != nil {
		w.logger.Debugw("failed to handle event", "event", p.WebhookEvent, "error", err)
	}

	return nil
}

// TODO: https://github.com/githubcustomers/SAP/issues/1002#issuecomment-1477526984
func (w *WorkflowAction) handleWorkflowRun(ctx context.Context, p *WorkflowActionParams) error {
	w.logger.Infow("Organization", "organization", p.Organization)
	w.logger.Infow("Repository", "repository", p.Repository)
	w.logger.Infow("Sender", "sender", p.WebhookEvent)
	w.logger.Infow("Workflow", "workflow", p.WorkflowId)
	return nil
}

func (w *WorkflowAction) handleWorkflowDispatch(ctx context.Context, p *WorkflowActionParams) error {
	w.logger.Infow("Organization", "organization", p.Organization)
	w.logger.Infow("Repository", "repository", p.Repository)
	w.logger.Infow("Sender", "sender", p.WebhookEvent)
	w.logger.Infow("Workflow", "workflow", p.WorkflowId)
	return nil
}

func (w *WorkflowAction) handleWorkflowJob(ctx context.Context, p *WorkflowActionParams) error {
	w.logger.Infow("Organization", "organization", p.Organization)
	w.logger.Infow("Repository", "repository", p.Repository)
	w.logger.Infow("Sender", "sender", p.WebhookEvent)
	w.logger.Infow("Workflow", "workflow", p.WorkflowId)
	return nil
}
