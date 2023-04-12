package clients

import (
	"fmt"
	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/utils"

	"go.uber.org/zap"
)

type ClientMap map[string]Client

type Client interface {
	Organization() string
	Repository() string
}

func InitClients(logger *zap.SugaredLogger, clientConfigs []config.Client) (ClientMap, error) {
	clients := make(ClientMap)

	for _, config := range clientConfigs {
		if config.GithubAuthConfig == nil {
			return nil, fmt.Errorf(utils.ErrMissingClientConfig, config.Name)
		}

		logger := logger.Named(config.Name)
		client, err := NewGithub(logger, config.OrganizationName, config.RepositoryName, config.ServerInfo, config.GithubAuthConfig)
		if err != nil {
			return nil, err
		}

		clients[config.Name] = client
	}

	return clients, nil
}
