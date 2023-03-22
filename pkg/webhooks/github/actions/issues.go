package actions

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v50/github"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"

	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
	"github.tools.sap/actions-rollout-app/pkg/git"
	"github.tools.sap/actions-rollout-app/pkg/utils"
)

type IssuesAction struct {
	logger *zap.SugaredLogger
	client *clients.Github

	targetRepos map[string]bool
}

type IssuesActionParams struct {
	PullRequestNumber int
	AuthorAssociation string
	RepositoryName    string
	RepositoryURL     string
	Comment           string
	CommentID         int64
}

type IssueParams struct {
	RepositoryName string
	RepositoryURL  string
	Title          string
	Body           string
	assignees      *[]string
	labels         *[]string
}

func NewIssuesAction(logger *zap.SugaredLogger, client *clients.Github, rawConfig map[string]any) (*IssuesAction, error) {
	var typedConfig config.IssuesCommentHandlerConfig
	err := mapstructure.Decode(rawConfig, &typedConfig)
	if err != nil {
		return nil, err
	}

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

func (i *IssuesAction) CreateIssue(ctx context.Context, p *IssueParams) error {
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
	switch utils.IssueCreateCommand(issueTitle) {
	case utils.IssueCreateFork:
		return i.createForkIssue(ctx, p)
	default:
		i.logger.Debugw("skip creating issue, message does not contain a valid command", "source-repo", p.RepositoryName)
	}
	return nil
}

func (i *IssuesAction) HandleIssueComment(ctx context.Context, p *IssuesActionParams) error {
	_, ok := i.targetRepos[p.RepositoryName]
	if !ok {
		i.logger.Debugw("skip handling issues comment action, not in list of target repositories", "source-repo", p.RepositoryName)
		return nil
	}

	_, ok = utils.AllowedAuthorAssociations[p.AuthorAssociation]
	if !ok {
		i.logger.Debugw("skip handling issues comment action, author is not allowed", "source-repo", p.RepositoryName, "association", p.AuthorAssociation)
		return nil
	}

	comment := strings.TrimSpace(p.Comment)

	_, ok = utils.IssueCommentCommands[utils.IssueCommentCommand(comment)]
	if !ok {
		i.logger.Debugw("skip handling issues comment action, message does not contain a valid command", "source-repo", p.RepositoryName)
		return nil
	}

	switch utils.IssueCommentCommand(comment) {
	case utils.IssueCommentBuildFork:
		return i.buildForkPR(ctx, p)
	default:
		i.logger.Debugw("skip handling issues comment action, message does not contain a valid command", "source-repo", p.RepositoryName)
		return nil
	}
}

func (i *IssuesAction) buildForkPR(ctx context.Context, p *IssuesActionParams) error {
	pullRequest, _, err := i.client.GetV3Client().PullRequests.Get(ctx, i.client.Organization(), p.RepositoryName, p.PullRequestNumber)
	if err != nil {
		return fmt.Errorf("error finding issue related pull request %w", err)
	}

	if pullRequest.Head.Repo.Fork == nil || !*pullRequest.Head.Repo.Fork {
		i.logger.Debugw("skip handling issues comment action, pull request is not from a fork", "source-repo", p.RepositoryName)
		return nil
	}

	token, err := i.client.GitToken(ctx)
	if err != nil {
		return fmt.Errorf("error creating git token %w", err)
	}

	targetRepoURL, err := url.Parse(p.RepositoryURL)
	if err != nil {
		return err
	}
	targetRepoURL.User = url.UserPassword("x-access-token", token)

	headRef := *pullRequest.Head.Ref
	commitMessage := "Triggering fork build approved by maintainer"
	err = git.PushToRemote(*pullRequest.Head.Repo.CloneURL, headRef, targetRepoURL.String(), "fork-build/"+headRef, commitMessage)
	if err != nil {
		return fmt.Errorf("error pushing to target remote repository %w", err)
	}

	i.logger.Infow("triggered fork build action by pushing to fork-build branch", "source-repo", p.RepositoryName, "branch", headRef)

	_, _, err = i.client.GetV3Client().Reactions.CreateIssueCommentReaction(ctx, i.client.Organization(), p.RepositoryName, p.CommentID, "rocket")
	if err != nil {
		return fmt.Errorf("error creating issue comment reaction %w", err)
	}

	err = git.DeleteBranch(targetRepoURL.String(), "fork-build/"+headRef)
	if err != nil {
		return err
	}

	return nil
}

func (i *IssuesAction) createForkIssue(ctx context.Context, p *IssueParams) error {
	_, _, err := i.client.GetV3Client().Issues.Create(ctx, i.client.Organization(), p.RepositoryName, &github.IssueRequest{
		Title:     github.String(p.Title),
		Body:      github.String(p.Body),
		Assignees: p.assignees,
		Labels:    p.labels,
	})
	if err != nil {
		i.logger.Errorw("error creating issue", "error", err)
	}
	return nil
}
