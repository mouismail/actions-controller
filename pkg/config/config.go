package config

import (
	"os"
	"strings"

	"sigs.k8s.io/yaml"
)

type Configuration struct {
	Clients  []Client  `json:"clients" description:"client configurations"`
	Webhooks []Webhook `json:"webhooks" description:"webhook configurations"`
	Repos    []Repo    `json:"repos" description:"repository configurations"`
	Raw      []byte
	Global   []GlobalConfig `json:"global" description:"global configurations"`
}

type GlobalConfig struct {
	Workers float64 `json:"workers" description:"the number of workers that will be used to process the configuration files"`
}

type Repo struct {
	Organization   string    `json:"organization" description:"the organization where the repository is located"`
	Repository     string    `json:"repository" description:"the repository where the configuration files are located"`
	FilesPath      *[]string `json:"files_path" description:"the path to the configuration files"`
	Branch         string    `json:"branch" description:"the branch where the configuration files are located"`
	WorkerPoolSize int       `json:"worker_pool_size" description:"the number of workers that will be used to process the configuration files"`
}

type Client struct {
	GithubAuthConfig *GithubClient `json:"github" description:"auth config a github client"`
	Name             string        `json:"name" description:"name of the client, used for referencing in webhook config"`
	OrganizationName string        `json:"organization" description:"name of the organization that this client will act on"`
	RepositoryName   string        `json:"repository" description:"the repository where the configuration files are located"`
	ServerInfo       *ServerInfo   `json:"server_info" description:"GitHub Enterprise server info for BaseURL and UploadURL"`
}

type GithubClient struct {
	AppID              int64  `json:"app-id" description:"application id of github app"`
	PrivateKeyCertPath string `json:"key-path" description:"private key pem path of github app"`
}

type Webhook struct {
	ServePath string         `json:"serve-path" description:"path of the webhook to serve on"`
	Secret    string         `json:"secret" description:"the webhook secret"`
	Actions   WebhookActions `json:"actions" description:"webhook actions"`
}

type ServerInfo struct {
	BaseURL       string `json:"base_url"`
	UploadURL     string `json:"upload_url"`
	EnterpriseURL string `json:"enterprise_url"`
}

type IssueCreatedHandlerConfig struct {
	IssueTitle string   `json:"issue_title" description:"the title of the issue"`
	IssueBody  string   `json:"issue_body" description:"the body of the issue"`
	Assignees  []string `json:"assignees" description:"the assignees for the issue"`
	Labels     []string `json:"labels" description:"the labels for the issue"`
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
