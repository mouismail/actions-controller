package actions

import (
	"context"
	"fmt"
	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/utils"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type WebhookActions struct {
	logger          *zap.SugaredLogger
	workflowActions []*WorkflowAction
	repoActions     []*RepoAction
}

func InitActions(logger *zap.SugaredLogger, cs clients.ClientMap, config config.WebhookActions) (*WebhookActions, error) {
	actions := WebhookActions{
		logger: logger,
	}

	for _, spec := range config {
		c, ok := cs[spec.Client]
		if !ok {
			return nil, fmt.Errorf(utils.ErrClientNotFound, spec.Client)
		}

		switch clientType := c.(type) {
		case *clients.Github:
		default:
			return nil, fmt.Errorf(utils.ErrInvalidClient, spec.Type, clientType)
		}

		switch t := spec.Type; t {
		//case utils.ActionIssuesHandler:
		//	h, err := NewIssuesAction(logger, c.(*clients.Github))
		//	if err != nil {
		//		return nil, err
		//	}
		//	actions.issueActions = append(actions.issueActions, h)
		case utils.ActionWorkflowHandler:
			h, err := NewWorkflowAction(logger, c.(*clients.Github), spec.Args)
			if err != nil {
				return nil, err
			}
			actions.workflowActions = append(actions.workflowActions, h)
		case utils.ActionRepoHandler:
			h, err := NewRepoAction(logger, c.(*clients.Github), spec.Args)
			if err != nil {
				return nil, err
			}
			actions.repoActions = append(actions.repoActions, h)
		default:
			return nil, fmt.Errorf(utils.ErrUnsupportedType, t)
		}

		logger.Debugw(utils.LoggerDebugInitWebhookAction, "name", spec.Type)
	}

	return &actions, nil
}

func (w *WebhookActions) ProcessWorkflowDispatchEvent(payload *ghwebhooks.WorkflowDispatchPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.WebhookHandleTimeout)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	for _, wa := range w.workflowActions {
		wa := wa
		group.Go(func() error {
			params := &WorkflowActionParams{
				Repository:   payload.Repository.Name,
				Organization: payload.Organization.Login,
				WorkflowName: payload.Workflow,
				WorkflowID:   0,
				WebhookEvent: ghwebhooks.WorkflowDispatchEvent,
			}

			err := wa.handleWorkflowDispatch(ctx, params)
			if err != nil {
				w.logger.Errorw(utils.LoggerErrorCreatingWorkflowDispatch, "source-repo", params.Repository, "error", err)
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		w.logger.Errorw(utils.LoggerErrorProcessingEvent, "error", err)
	}
}

func (w *WebhookActions) ProcessWorkflowJobEvent(payload *ghwebhooks.WorkflowJobPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.WebhookHandleTimeout)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	for _, wa := range w.workflowActions {
		wa := wa
		group.Go(func() error {
			params := &WorkflowActionParams{
				Repository:   payload.Repository.Name,
				Organization: payload.Organization.Login,
				WorkflowName: payload.WorkflowJob.Name,
				WorkflowID:   payload.WorkflowJob.ID,
				WebhookEvent: ghwebhooks.WorkflowDispatchEvent,
				Sender:       payload.Sender.Login,
			}
			err := wa.handleWorkflowJob(ctx, params)
			if err != nil {
				w.logger.Errorw(utils.LoggerErrorCreatingWorkflowJob, "source-repo", params.Repository, "error", err)
				return err
			}
			return nil
		})
	}

	if err := group.Wait(); err != nil {
		w.logger.Errorw(utils.LoggerErrorProcessingEvent, "error", err)
	}
}

func (w *WebhookActions) ProcessWorkflowRunEvent(payload *ghwebhooks.WorkflowRunPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.WebhookHandleTimeout)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	for _, wa := range w.workflowActions {
		wa := wa
		group.Go(func() error {
			params := &WorkflowActionParams{
				Repository:   payload.Repository.Name,
				Organization: payload.Organization.Login,
				WorkflowName: payload.Workflow.Name,
				WorkflowID:   payload.Workflow.ID,
				WebhookEvent: ghwebhooks.WorkflowRunEvent,
				Sender:       payload.Sender.Login,
			}

			err := wa.handleWorkflowRun(ctx, params)
			if err != nil {
				w.logger.Errorw(utils.LoggerErrorCreatingWorkflowRun, "source-repo", params.Repository, "error", err)
				return err
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		w.logger.Errorw(utils.LoggerErrorProcessingEvent, "error", err)
	}
}
