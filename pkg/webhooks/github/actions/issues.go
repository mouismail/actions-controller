package actions

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v50/github"
	"go.uber.org/zap"

	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
	"github.tools.sap/actions-rollout-app/pkg/utils"
)

type IssuesAction struct {
	logger *zap.SugaredLogger
	client *clients.Github

	targetRepos map[string]bool
}

type IssueCreateParams struct {
	RepositoryName string
	Organization   string
	Title          string
	Body           string
	Assignees      *[]string
	Labels         *[]string
}

func NewIssuesAction(logger *zap.SugaredLogger, client *clients.Github) (*IssuesAction, error) {
	var typedConfig config.IssuesCommentHandlerConfig

	targetRepos := make(map[string]bool)
	for name := range typedConfig.TargetRepos {
		targetRepos[name] = true
	}

	return &IssuesAction{
		logger:      logger,
		client:      client,
		targetRepos: targetRepos,
	}, nil
}

func (i *IssuesAction) CreateIssue(ctx context.Context, p *IssueCreateParams) error {
	_, ok := i.targetRepos[p.RepositoryName]
	if !ok {
		i.logger.Debugw("skip creating issue, not in list of target repositories", "source-repo", p.RepositoryName)
		return nil
	}

	issueTitle := strings.TrimSpace(p.Title)

	_, ok = utils.IssueCreateCommands[utils.IssueCreateCommand(issueTitle)]
	if !ok {
		i.logger.Debugw("skip creating issue, message does not contain a valid command", "source-repo", p.RepositoryName)
		return nil
	}
	err := i.createForkIssue(ctx, p)
	if err != nil {
		i.logger.Errorf("error creating issue: %v", err)
	}
	return nil
}

func (i *IssuesAction) createForkIssue(ctx context.Context, p *IssueCreateParams) (err error) {

	issue, _, err := i.client.GetV3Client().Issues.Create(ctx, i.client.Organization(), p.RepositoryName, &github.IssueRequest{
		Title:     github.String(p.Title),
		Body:      github.String(p.Body),
		Assignees: p.Assignees,
		Labels:    p.Labels,
	})
	if err != nil {
		i.logger.Errorw("error creating issue", "error", err)
		return fmt.Errorf("error creating issue: %w", err)
	}

	i.logger.Infow("issue created for workflow ", "issue", issue.ID, p.Title)

	return nil
}
