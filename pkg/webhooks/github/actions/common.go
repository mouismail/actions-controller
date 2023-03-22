package actions

import (
	"context"
	"fmt"

	"reflect"
	"strconv"
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

		ghc, ok := c.(*clients.Github)
		if !ok {
			return nil, fmt.Errorf(utils.ErrInvalidClient, spec.Type, reflect.TypeOf(c))
		}

		switch t := spec.Type; t {
		case utils.ActionIssuesHandler:
			h, err := NewIssuesAction(logger, ghc, spec.Args)
			if err != nil {
				return nil, err
			}
			actions.ih = append(actions.ih, h)
		case utils.ActionWorkflowHandler:
			h, err := NewWorkflowAction(logger, ghc, spec.Args)
			if err != nil {
				return nil, err
			}
			actions.wa = append(actions.wa, h)
		default:
			return nil, fmt.Errorf(utils.ErrUnsupportedType, t)
		}

		logger.Debugw("initialized github webhook action", "name", spec.Type)
	}

	return &actions, nil
}

func (w *WebhookActions) ProcessIssueCommentEvent(payload *ghwebhooks.IssueCommentPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), utils.WebhookHandleTimeout)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	for _, i := range w.ih {
		i := i
		g.Go(func() error {
			if payload.Action != "created" {
				return nil
			}
			if payload.Issue.PullRequest == nil {
				return nil
			}

			parts := strings.Split(payload.Issue.PullRequest.URL, "/")
			pullRequestNumberString := parts[len(parts)-1]
			pullRequestNumber, err := strconv.ParseInt(pullRequestNumberString, 10, 64)
			if err != nil {
				return err
			}

			params := &IssuesActionParams{
				RepositoryName:    payload.Repository.Name,
				RepositoryURL:     payload.Repository.CloneURL,
				Comment:           payload.Comment.Body,
				CommentID:         payload.Comment.ID,
				AuthorAssociation: payload.Comment.AuthorAssociation,
				PullRequestNumber: int(pullRequestNumber),
			}

			err = i.HandleIssueComment(ctx, params)
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
				// TODO: we need to get the run Id from the workflow dispatch event
				WorkflowId:   123,
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
				WorkflowId:   payload.WorkflowJob.ID,
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
				WorkflowId:   payload.Workflow.ID,
				WebhookEvent: ghwebhooks.WorkflowRunEvent,
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

func extractTag(payload *ghwebhooks.PushPayload) string {
	return strings.Replace(payload.Ref, "refs/tags/", "", 1)
}
