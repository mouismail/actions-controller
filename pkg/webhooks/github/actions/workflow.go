package actions

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"github.com/google/go-github/v50/github"
	"go.uber.org/zap"

	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/utils"
)

type WorkflowActionParams struct {
	WorkflowName string
	WorkflowID   int64
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
	workerPoolSize float64
	filesPath      *[]string
	assignees      *[]string
}

// TODO: retest this

func NewWorkflowAction(logger *zap.SugaredLogger, client *clients.Github, rawConfig map[string]any) (*WorkflowAction, error) {
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
		logger:       logger,
		client:       client,
		repository:   client.Repository(),
		organization: client.Organization(),
		filesPath:    &files,
		assignees:    &assignees,
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

func (w *WorkflowAction) disableWorkflow(ctx context.Context, p *WorkflowActionParams, workflowID int64) error {
	c := config.Client{
		GithubAuthConfig: w.client.GetConfig(),
		Name:             "disable-workflow",
		OrganizationName: p.Organization,
		RepositoryName:   p.Repository,
		ServerInfo:       w.client.ServerInfo(),
	}
	configClients := []config.Client{c}
	workflowClients, err := clients.InitClients(w.logger, configClients)
	if err != nil {
		return err
	}

	for _, workflowClient := range workflowClients {
		resp, workflowErr := workflowClient.GetV3Client().Actions.DisableWorkflowByID(ctx, p.Organization, p.Repository, workflowID)
		if workflowErr != nil {
			return workflowErr
		}

		if resp.StatusCode != http.StatusNoContent {
			return errors.New(strconv.Itoa(resp.StatusCode))
		}
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
		Labels:    &labels,
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
	// TODO: to be validated with @stoe
	//err := w.handleWorkflowEvent(ctx, p, "job")
	//if err != nil {
	//	return err
	//}
	return nil
}

func (w *WorkflowAction) handleWorkflowEvent(ctx context.Context, p *WorkflowActionParams, eventType string) error {
	title := fmt.Sprintf("[%d] - %s/%s", p.WorkflowID, p.Organization, p.Repository)
	message, err := w.generateWorkflowMessage(eventType, p)

	if err != nil {
		return err
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

	err = repoAction.HandleRepo(ctx, repoParams)
	if err != nil {

		disableErr := w.disableWorkflow(ctx, p, p.WorkflowID)
		if disableErr != nil {
			return disableErr
		}
		w.logger.Infow("workflow disabled", "workflow_id", p.WorkflowID)
		return w.createWorkflowIssue(ctx, title, message, *w.assignees, []string{fmt.Sprintf("%s/%s", p.Organization, p.Repository), "not-valid"})
	}
	return nil
}

func (w *WorkflowAction) generateWorkflowMessage(eventType string, p *WorkflowActionParams) (string, error) {
	var message string
	switch eventType {
	case "run":
		message = fmt.Sprintf(utils.WorkflowRunMessage,
			p.Organization,
			p.Repository,
			w.client.ServerInfo().EnterpriseURL,
			p.Organization,
			p.Repository,
			p.Sender,
			p.Organization,
			p.Repository,
			p.WorkflowName,
			p.WorkflowID)

	case "dispatch":
		message = fmt.Sprintf(utils.WorkflowDispatchMessage,
			p.Organization,
			p.Repository,
			w.client.ServerInfo().EnterpriseURL,
			p.Organization,
			p.Repository,
			p.Sender,
			p.Organization,
			p.Repository,
			p.WorkflowName,
			p.WorkflowID)
	case "job":
		message = fmt.Sprintf(utils.WorkflowJobMessage,
			p.Organization,
			p.Repository,
			w.client.ServerInfo().EnterpriseURL,
			p.Organization,
			p.Repository,
			p.Sender,
			p.Organization,
			p.Repository,
			p.WorkflowName,
			p.WorkflowID)
	default:
		return "", errors.New("unsupported event type")
	}
	return message, nil
}

func (w *WorkflowAction) disableWorkflowByOrganization(ctx context.Context, p *WorkflowActionParams, repoID int64) error {
	enabledRepositories := "none"

	c := config.Client{
		GithubAuthConfig: w.client.GetConfig(),
		Name:             "disable-workflow",
		OrganizationName: p.Organization,
		RepositoryName:   p.Repository,
		ServerInfo:       w.client.ServerInfo(),
	}
	configClients := []config.Client{c}
	workflowClients, err := clients.InitClients(w.logger, configClients)
	if err != nil {
		return err
	}

	for _, workflowClient := range workflowClients {
		_, resp, workflowErr := workflowClient.GetV3Client().Organizations.EditActionsPermissions(ctx, p.Organization, github.ActionsPermissions{
			EnabledRepositories: &enabledRepositories,
		})
		if workflowErr != nil {
			return workflowErr
		}

		if resp.StatusCode != http.StatusNoContent {
			return errors.New(strconv.Itoa(resp.StatusCode))
		}
	}

	return nil
}

func (w *WorkflowAction) disableWorkflowForRepo(ctx context.Context, p *WorkflowActionParams, repoID int64) error {
	c := config.Client{
		GithubAuthConfig: w.client.GetConfig(),
		Name:             "disable-workflow",
		OrganizationName: p.Organization,
		RepositoryName:   p.Repository,
		ServerInfo:       w.client.ServerInfo(),
	}
	configClients := []config.Client{c}
	workflowClients, err := clients.InitClients(w.logger, configClients)
	if err != nil {
		return err
	}

	for _, workflowClient := range workflowClients {
		resp, workflowErr := workflowClient.GetV3Client().Actions.RemoveEnabledRepoInOrg(ctx, p.Organization, repoID)
		if workflowErr != nil {
			return workflowErr
		}

		if resp.StatusCode != http.StatusNoContent {
			return errors.New(strconv.Itoa(resp.StatusCode))
		}
	}

	return nil
}
