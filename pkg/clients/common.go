package clients

import (
	"fmt"

	v3 "github.com/google/go-github/v50/github"
	"go.uber.org/zap"

	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/utils"
)

type ClientMap map[string]Client

type Client interface {
	Organization() string
	Repository() string
	ServerInfo() *config.ServerInfo
	GetConfig() *config.GithubClient
	GetV3Client() *v3.Client
}

func InitClients(logger *zap.SugaredLogger, clientConfigs []config.Client) (ClientMap, error) {
	clients := make(ClientMap)

	for _, clientConfig := range clientConfigs {
		if clientConfig.GithubAuthConfig == nil {
			return nil, fmt.Errorf(utils.ErrMissingClientConfig, clientConfig.Name)
		}

		logger := logger.Named(clientConfig.Name)
		client, err := NewGithub(logger, clientConfig.OrganizationName, clientConfig.RepositoryName, clientConfig.ServerInfo, clientConfig.GithubAuthConfig)
		if err != nil {
			return nil, err
		}

		clients[clientConfig.Name] = client
	}

	return clients, nil
}
