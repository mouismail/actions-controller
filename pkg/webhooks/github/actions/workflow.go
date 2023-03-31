package actions

import (
	"context"
	"errors"
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"github.com/google/go-github/v50/github"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
	"go.uber.org/zap"
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

	repository   string
	organization string
	workflow     string
}

func NewWorkflowAction(logger *zap.SugaredLogger, client *clients.Github, rawConfig map[string]any) (*WorkflowAction, error) {
	repo, ok := rawConfig["repository"].(string)
	if !ok {
		return nil, errors.New("repo not found or is not a string")
	}

	org := rawConfig["organization"].(string)
	if !ok {
		return nil, errors.New("org not found or is not a string")
	}

	workflow, ok := rawConfig["workflow"].(string)
	if !ok {
		return nil, errors.New("workflowId not found or is not an int64")
	}

	typedConfig := config.WorkflowHandleConfig{
		Repo:     repo,
		Org:      org,
		Workflow: workflow,
	}

	return &WorkflowAction{
		logger:       logger,
		client:       client,
		repository:   typedConfig.Repo,
		organization: typedConfig.Org,
		workflow:     typedConfig.Workflow,
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

	if w.workflow != p.Workflow {
		w.logger.Debugw("workflow id does not match", "expected", w.workflow, "actual", p.Workflow)
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
// TODO: combine all those func to one function

func (w *WorkflowAction) handleWorkflowRun(ctx context.Context, p *WorkflowActionParams) error {

	issue, issueResp, err := w.client.GetV3Client().Issues.Create(ctx, w.client.Organization(), "sap-demos", &github.IssueRequest{
		Title:    github.String(p.Organization + "/" + p.Repository + " - " + p.Workflow),
		Body:     github.String("Workflow Run triggered by " + p.Sender + " on organization " + p.Organization + " and repository " + p.Repository + ""),
		Assignee: github.String("mouismail"),
	})
	if err != nil {
		w.logger.Errorw("error creating issue", "error", err)
	} else {
		w.logger.Infow("issue created", "issue_id", issue.ID, "response", issueResp.StatusCode)
	}

	return nil

	return nil
}

func (w *WorkflowAction) handleWorkflowDispatch(ctx context.Context, p *WorkflowActionParams) error {
	issue, issueResp, err := w.client.GetV3Client().Issues.Create(ctx, w.organization, w.repository, &github.IssueRequest{
		Title:    github.String(p.Organization + "/" + p.Repository + " - " + p.Workflow),
		Body:     github.String("Workflow dispatch triggered by " + p.Sender + " on organization " + p.Organization + " and repository " + p.Repository + ""),
		Assignee: github.String("mouismail"),
	})
	if err != nil {
		w.logger.Errorw("error creating issue", "error", err)
	} else {
		w.logger.Infow("issue created", "issue_id", issue.ID, "response", issueResp.StatusCode)
	}

	return nil
}

func (w *WorkflowAction) handleWorkflowJob(ctx context.Context, p *WorkflowActionParams) error {
	issue, issueResp, err := w.client.GetV3Client().Issues.Create(ctx, w.client.Organization(), "sap-demos", &github.IssueRequest{
		Title:    github.String(p.Organization + "/" + p.Repository + " - " + p.Workflow),
		Body:     github.String("Workflow Job triggered by " + p.Sender + " on organization " + p.Organization + " and repository " + p.Repository + ""),
		Assignee: github.String("mouismail"),
	})
	if err != nil {
		w.logger.Errorw("error creating issue", "error", err)
	} else {
		w.logger.Infow("issue created", "issue_id", issue.ID, "response", issueResp.StatusCode)
	}

	return nil
}
