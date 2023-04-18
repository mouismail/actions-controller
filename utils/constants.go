package utils

import "time"

const (
	WebhookHandleTimeout                     = 240 * time.Second
	ErrClientNotFound                        = "webhook action client not found: %s"
	ErrUnsupportedType                       = "handler type not supported: %s"
	ErrInvalidClient                         = "action %s only supports github clients, not: %s"
	ErrMissingClientConfig                   = "client config missing for action %s"
	ErrMissingClient                         = "error creating github app client %w"
	ErrMissingEnterpriseClient               = "error creating github enterprise client %w"
	ErrCreatingInstallationToken             = "error creating Installation id token %w"
	ErrFindingOrgInstallations               = "error finding organization app installations %w"
	ErrInvalidConfigOrganization             = "invalid url for organization or the url is not the same as the organization name where the workflow is triggered"
	ErrInvalidConfigRepository               = "invalid url for repository or the url is not the same as the repository name where the workflow is triggered"
	ErrInvalidConfigOrganizationRepositories = "invalid url for organization or the url is not the same as the organization name where the workflow is triggered"
	ErrInvalidContactEmail                   = "invalid contact email or empty"
	ErrInvalidUseCase                        = "invalid use case or empty"
	ErrValidationEmptyContent                = "content is empty or nil"
	ActionWorkflowHandler                    = "workflow-handling"
	ActionRepoHandler                        = "repo-handling"
	DefaultLocalRef                          = "refs/heads"
	LoggerDebugInitWebhookAction             = "initialized github webhook action"
	LoggerErrorCreatingWorkflowDispatch      = "error in workflow dispatch handler action"
	LoggerErrorProcessingEvent               = "error processing event"
	LoggerErrorCreatingWorkflowJob           = "error in workflow Job handler action"
	LoggerErrorCreatingWorkflowRun           = "error in workflow Run handler action"
	WorkflowRunMessage                       = `
<a href='link' target="_blank"><img alt='Workflow status' src='https://img.shields.io/badge/Workflow_status - Disabled-100000?style=flat-square&logo=Workflow status&logoColor=white&labelColor=CD1111&color=E73B3B'/></a>

## Actions Controller

<a href='#' target="#"><img alt='Workflow status' src='https://img.shields.io/badge/Workflow_status - Disabled-100000?style=for-the-badge&logo=Workflow status&logoColor=white&labelColor=CD1111&color=E73B3B'/></a>

:rocket: Actions Controller bot has detected changes in [%s/%s](%s/%s/%s) on Workflow Run event

### :red_circle: Details
| Sender         | Organization       | Repository            | Workflow Name | Workflow ID |
| ---------------|---------------------|-----------------------|---------------|-------------|
| @%s      | %s     | %s    | %s            | %d         |`
	WorkflowDispatchMessage = `
<a href='link' target="_blank"><img alt='Workflow status' src='https://img.shields.io/badge/Workflow_status - Disabled-100000?style=flat-square&logo=Workflow status&logoColor=white&labelColor=CD1111&color=E73B3B'/></a>

## Actions Controller

:rocket: Actions Controller bot has detected changes in [%s/%s](%s/%s/%s)  on Workflow Dispatch event

### :red_circle: Details
| Sender         | Organization       | Repository            | Workflow Name | Workflow ID |
| ---------------|---------------------|-----------------------|---------------|-------------|
| @%s      | %s     | %s    | %s            | %d         |`
	WorkflowJobMessage = `
<a href='link' target="_blank"><img alt='Workflow status' src='https://img.shields.io/badge/Workflow_status - Disabled-100000?style=flat-square&logo=Workflow status&logoColor=white&labelColor=CD1111&color=E73B3B'/></a>

## Actions Controller

<a href='link' target="_blank"><img alt='Workflow status' src='https://img.shields.io/badge/Workflow_status - Disabled-100000?style=flat-square&logo=Workflow status&logoColor=white&labelColor=CD1111&color=E73B3B'/></a>

:rocket: Actions Controller bot has detected changes in [%s/%s](%s/%s/%s)  on Workflow Job event

### :red_circle: Details
| Sender         | Organization       | Repository            | Workflow Name | Workflow ID |
| ---------------|---------------------|-----------------------|---------------|-------------|
| @%s      | %s     | %s    | %s            | %d         |`
)
