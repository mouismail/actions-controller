package utils

import "time"

type IssueCommentCommand string
type IssueCreateCommand string

const (
	WebhookHandleTimeout                             = 240 * time.Second
	ErrClientNotFound                                = "webhook action client not found: %s"
	ErrUnsupportedType                               = "handler type not supported: %s"
	ErrInvalidClient                                 = "action %s only supports github clients, not: %s"
	ErrMissingClientConfig                           = "client config missing for action %s"
	ErrMissingClient                                 = "error creating github app client %w"
	ErrMissingEnterpriseClient                       = "error creating github enterprise client %w"
	ErrCreatingInstallationToken                     = "error creating Installation id token %w"
	ErrFindingOrgInstallations                       = "error finding organization app installations %w"
	ErrInvalidConfigOrganization                     = "invalid url for organization or the url is not the same as the organization name where the workflow is triggered"
	ErrInvalidConfigRepository                       = "invalid url for repository or the url is not the same as the repository name where the workflow is triggered"
	ErrInvalidContactEmail                           = "invalid contact email or empty"
	ErrInvalidUseCase                                = "invalid use case or empty"
	ErrValidationEmptyContent                        = "content is empty or nil"
	ActionIssuesHandler          string              = "issue-handling"
	ActionWorkflowHandler        string              = "workflow-handling"
	ActionRepoHandler            string              = "repo-handling"
	IssueCommentCommandPrefix                        = "/"
	IssueCommentBuildFork        IssueCommentCommand = IssueCommentCommandPrefix + "ok-to-build"
	IssueCreateCommandPrefix                         = "/"
	IssueCreateFork              IssueCreateCommand  = IssueCreateCommandPrefix + "ok-to-create"
	DefaultLocalRef                                  = "refs/heads"
	WorkflowRunMessage                               = "Workflow Run triggered by %s on organization %s and repository %s"
	WorkflowJobMessage                               = "%s triggered by %s on organization %s and repository %s"
	WorkflowDispatchMessage                          = "Workflow dispatch triggered by %s on organization %s and repository %s"
)

var (
	IssueCommentCommands = map[IssueCommentCommand]bool{
		IssueCommentBuildFork: true,
	}

	IssueCreateCommands = map[IssueCreateCommand]bool{
		IssueCreateFork: true,
	}

	AllowedAuthorAssociations = map[string]bool{
		"COLLABORATOR": true,
		"MEMBER":       true,
		"OWNER":        true,
	}
)
