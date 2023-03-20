package clients

import (
	"fmt"

	"github.tools.sap/actions-rollout-app/pkg/config"

	"go.uber.org/zap"
)

type ClientMap map[string]Client

type Client interface {
	VCS() config.VCSType
	Organization() string
}

func InitClients(logger *zap.SugaredLogger, config []config.Client) (ClientMap, error) {
	cs := ClientMap{}
	for _, clientConfig := range config {
		ghConfig := clientConfig.GithubAuthConfig

		if ghConfig == nil {
			return nil, fmt.Errorf("github client config must be provided for client %q", clientConfig.Name)
		}

		if ghConfig != nil {
			client, err := NewGithub(logger.Named(clientConfig.Name), clientConfig.OrganizationName, ghConfig)
			if err != nil {
				return nil, err
			}

			cs[clientConfig.Name] = client
		}
	}
	return cs, nil
}
