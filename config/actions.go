package config

type WebhookActions []WebhookAction

type WebhookAction struct {
	Type   string         `json:"type" description:"name of the webhook action"`
	Client string         `json:"client" description:"client that this webhook action uses"`
	Args   map[string]any `json:"args" description:"action configuration"`
}

type IssuesCommentHandlerConfig struct {
	TargetRepos map[string]any `mapstructure:"repos" description:"the repositories for which issue comment handling will be applied"`
}

type WorkflowHandleConfig struct {
	Repo     string `mapstructure:"repo" description:"the repository where that a workflow got triggered"`
	Org      string `mapstructure:"org" description:"the organization where that a workflow got triggered"`
	Workflow string `mapstructure:"workflow" description:"the id of the workflow that got triggered"`
}
