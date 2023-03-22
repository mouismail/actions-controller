package config

import (
	"os"
	"strings"

	"sigs.k8s.io/yaml"
)

// TODO: VCSType is not needed to be removed

type VCSType string

const (
	Github VCSType = "github"
)

type Configuration struct {
	Clients  []Client  `json:"clients" description:"client configurations"`
	Webhooks []Webhook `json:"webhooks" description:"webhook configurations"`
	Raw      []byte
}

type Client struct {
	GithubAuthConfig *GithubClient `json:"github" description:"auth config a github client"`
	Name             string        `json:"name" description:"name of the client, used for referencing in webhook config"`
	OrganizationName string        `json:"organization" description:"name of the organization that this client will act on"`
	ServerInfo       *ServerInfo   `json:"server_info" description:"GitHub Enterprise server info for BaseURL and UploadURL"`
}

type GithubClient struct {
	AppID              int64  `json:"app-id" description:"application id of github app"`
	PrivateKeyCertPath string `json:"key-path" description:"private key pem path of github app"`
}

type GitlabClient struct {
	Token string `json:"token" description:"auth token for gitlab client"`
}

type Webhook struct {
	VCS       VCSType        `json:"vcs" description:"type of the vcs"`
	ServePath string         `json:"serve-path" description:"path of the webhook to serve on"`
	Secret    string         `json:"secret" description:"the webhook secret"`
	Actions   WebhookActions `json:"actions" description:"webhook actions"`
}

type ServerInfo struct {
	BaseURL   string `json:"base_url"`
	UploadURL string `json:"upload_url"`
}

type IssueCreatedHandlerConfig struct {
	IssueTitle string   `mapstructure:"issue-title" description:"the title of the issue"`
	IssueBody  string   `mapstructure:"issue-body" description:"the body of the issue"`
	Assignees  []string `mapstructure:"assignees" description:"the assignees for the issue"`
	Labels     []string `mapstructure:"labels" description:"the labels for the issue"`
}

func New(configPath string) (*Configuration, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Configuration{Raw: data}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (w WebhookActions) String() string {
	var actions []string
	for _, h := range w {
		actions = append(actions, h.Type)
	}
	return strings.Join(actions, ", ")
}
