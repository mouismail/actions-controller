package utils

import "time"

type IssueCommentCommand string
type IssueCreateCommand string

const (
	// WebhookHandleTimeout is the duration after which a webhook handle function context times out
	// the entire thing is asynchronous anyway, so the VCS will get an immediate response, this is just
	// that we do not have processing of events hanging internally
	WebhookHandleTimeout                          = 120 * time.Second
	ErrClientNotFound                             = "webhook action client not found: %s"
	ErrUnsupportedType                            = "handler type not supported: %s"
	ErrInvalidClient                              = "action %s only supports github clients, not: %s"
	ErrMissingClintConfig                         = "client config missing for action %s"
	ActionIssuesHandler       string              = "issue-handling"
	ActionWorkflowHandler     string              = "workflow-handling"
	IssueCommentCommandPrefix                     = "/"
	IssueCommentBuildFork     IssueCommentCommand = IssueCommentCommandPrefix + "ok-to-build"
	IssueCreateCommandPrefix                      = "/"
	IssueCreateFork           IssueCreateCommand  = IssueCreateCommandPrefix + "ok-to-create"
	DefaultLocalRef                               = "refs/heads"
	DefaultAuthor                                 = "actions-controller"
	DefaultAuthorMail                             = "actions-controller@sap.com"
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
