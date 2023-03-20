package actions

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
)

type WorkflowActionParams struct {
	WorkflowId   int64
	Organization string
	Repository   string
	WebhookEvent string
}

type WorkflowAction struct {
	logger *zap.SugaredLogger
	client *clients.Github

	repository   string
	organization string
	workflowId   int64
}

func NewWorkflowAction(logger *zap.SugaredLogger, client *clients.Github, rawConfig map[string]any) (*WorkflowAction, error) {
	var typedConfig config.WorkflowHandleConfig
	err := mapstructure.Decode(rawConfig, &typedConfig)
	if err != nil {
		return nil, err
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
	ok := w.repository == p.Repository
	if !ok {
		w.logger.Debugw("repository does not match", "expected", w.repository, "actual", p.Repository)
		return nil
	}

	ok = w.organization == p.Organization
	if !ok {
		w.logger.Debugw("organization does not match", "expected", w.organization, "actual", p.Organization)
		return nil
	}

	ok = w.workflowId == p.WorkflowId
	if !ok {
		w.logger.Debugw("workflow id does not match", "expected", w.workflowId, "actual", p.WorkflowId)
		return nil
	}

	switch p.WebhookEvent {
	case "workflow_run":
		w.logger.Debugw("handling workflow run event")
		err := w.handleWorkflowRun(ctx, p)
		if err != nil {
			w.logger.Debugw("failed to handle workflow run event", "error", err)
		}
	case "workflow_dispatch":
		w.logger.Debugw("handling workflow dispatch event")
		err := w.handleWorkflowDispatch(ctx, p)
		if err != nil {
			w.logger.Debugw("failed to handle workflow dispatch event", "error", err)
		}
	case "workflow_job":
		w.logger.Debugw("handling workflow job event")
		err := w.handleWorkflowJob(ctx, p)
		if err != nil {
			w.logger.Debugw("failed to handle workflow job event", "error", err)
		}
	}

	return nil
}

func (w *WorkflowAction) handleWorkflowRun(ctx context.Context, p *WorkflowActionParams) error {
	return nil
}

func (w *WorkflowAction) handleWorkflowDispatch(ctx context.Context, p *WorkflowActionParams) error {
	return nil
}

func (w *WorkflowAction) handleWorkflowJob(ctx context.Context, p *WorkflowActionParams) error {
	return nil
}
