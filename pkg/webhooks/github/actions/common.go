package actions

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
	"github.tools.sap/actions-rollout-app/pkg/webhooks/constants"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	ActionIssuesHandler   string = "issue-handling"
	ActionWorkflowHandler string = "workflow-handling"
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
			return nil, fmt.Errorf("webhook action client not found: %s", spec.Client)
		}

		switch clientType := c.(type) {
		case *clients.Github:
		default:
			return nil, fmt.Errorf("action %s only supports github clients, not: %s", spec.Type, clientType)
		}

		switch t := spec.Type; t {
		case ActionIssuesHandler:
			h, err := NewIssuesAction(logger, c.(*clients.Github), spec.Args)
			if err != nil {
				return nil, err
			}
			actions.ih = append(actions.ih, h)
		default:
			return nil, fmt.Errorf("handler type not supported: %s", t)
		}

		logger.Debugw("initialized github webhook action", "name", spec.Type)
	}

	return &actions, nil
}

func (w *WebhookActions) ProcessIssueCommentEvent(payload *ghwebhooks.IssueCommentPayload) {
	ctx, cancel := context.WithTimeout(context.Background(), constants.WebhookHandleTimeout)
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

}

func extractTag(payload *ghwebhooks.PushPayload) string {
	return strings.Replace(payload.Ref, "refs/tags/", "", 1)
}
