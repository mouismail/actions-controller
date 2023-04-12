package webhooks

import (
	"net/http"

	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/webhooks/github"

	"go.uber.org/zap"
)

func InitWebhooks(logger *zap.SugaredLogger, cs clients.ClientMap, c *config.Configuration) error {
	for _, w := range c.Webhooks {
		controller, err := github.NewGithubWebhook(logger.Named("github-webhook"), w, cs)
		if err != nil {
			return err
		}
		http.HandleFunc(w.ServePath, controller.Handle)
		logger.Infow("initialized github webhook", "serve-path", w.ServePath)
	}
	return nil
}
