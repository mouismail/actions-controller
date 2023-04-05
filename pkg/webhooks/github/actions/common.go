package actions

import (
	"context"
	"fmt"
	"strings"

	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
	"github.tools.sap/actions-rollout-app/pkg/utils"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type WebhookActions struct {
	logger *zap.SugaredLogger
	ih     []*IssuesAction
	wa     []*WorkflowAction
	ra     []*RepoAction
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
		case utils.ActionIssuesHandler:
			h, err := NewIssuesAction(logger, c.(*clients.Github))
			if err != nil {
				return nil, err
			}
			actions.ih = append(actions.ih, h)
		case utils.ActionWorkflowHandler:
			h, err := NewWorkflowAction(logger, c.(*clients.Github), spec.Args)
			if err != nil {
				return nil, err
			}
			actions.wa = append(actions.wa, h)
		case utils.ActionRepoHandler:
			h, err := NewRepoAction(logger, c.(*clients.Github), spec.Args)
			if err != nil {
				return nil, err
			}
			actions.ra = append(actions.ra, h)
		default:
			return nil, fmt.Errorf(utils.ErrUnsupportedType, t)
		}

		logger.Debugw("initialized github webhook action", "name", spec.Type)
	}

	return &actions, nil
}

func (w *WebhookActions) ProcessWorkflowDispatchEvent(payload *ghwebhooks.WorkflowDispatchPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.WebhookHandleTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	for _, wa := range w.wa {
		wa := wa
		g.Go(func() error {
			params := &WorkflowActionParams{
				Repository:   payload.Repository.Name,
				Organization: payload.Organization.Login,
				// TODO: we need to get the run ID from the workflow dispatch event
				Workflow:     payload.Workflow,
				WebhookEvent: ghwebhooks.WorkflowDispatchEvent,
			}

			err := wa.handleWorkflowDispatch(ctx, params)
			if err != nil {
				w.logger.Errorw("error in workflow dispatch handler action", "source-repo", params.Repository, "error", err)
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		w.logger.Errorw("errors processing event", "error", err)
	}
}

func (w *WebhookActions) ProcessWorkflowJobEvent(payload *ghwebhooks.WorkflowJobPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.WebhookHandleTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	for _, wa := range w.wa {
		wa := wa
		g.Go(func() error {
			params := &WorkflowActionParams{
				Repository:   payload.Repository.Name,
				Organization: payload.Organization.Login,
				Workflow:     payload.WorkflowJob.Name,
				WebhookEvent: ghwebhooks.WorkflowDispatchEvent,
			}
			err := wa.handleWorkflowJob(ctx, params)
			if err != nil {
				w.logger.Errorw("error in workflow dispatch handler action", "source-repo", params.Repository, "error", err)
				return err
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		w.logger.Errorw("errors processing event", "error", err)
	}
}

func (w *WebhookActions) ProcessWorkflowRunEvent(payload *ghwebhooks.WorkflowRunPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.WebhookHandleTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	for _, wa := range w.wa {
		wa := wa
		g.Go(func() error {
			params := &WorkflowActionParams{
				Repository:   payload.Repository.Name,
				Organization: payload.Organization.Login,
				Workflow:     payload.Workflow.Name,
				WebhookEvent: ghwebhooks.WorkflowRunEvent,
			}

			err := wa.handleWorkflowRun(ctx, params)
			if err != nil {
				w.logger.Errorw("error in workflow dispatch handler action", "source-repo", params.Repository, "error", err)
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		w.logger.Errorw("errors processing event", "error", err)
	}
}

func (w *WebhookActions) ProcessIssueCreate(repo, org, title, body string) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.WebhookHandleTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	for _, i := range w.ih {
		i := i
		g.Go(func() error {
			params := &IssueCreateParams{
				RepositoryName: repo,
				Organization:   org,
				Title:          title,
				Body:           body,
				Assignees:      nil,
				Labels:         nil,
			}

			err := i.CreateIssue(ctx, params)
			if err != nil {
				w.logger.Errorw("error in issue comment handler action", "source-repo", params.RepositoryName, "error", err)
				return err
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		w.logger.Errorw("errors processing event", "error", err)
	}
}

func extractTag(payload *ghwebhooks.PushPayload) string {
	return strings.Replace(payload.Ref, "refs/tags/", "", 1)
}
