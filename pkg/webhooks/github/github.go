package github

import (
	"errors"
	"net/http"
	"os"

	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/webhooks/github/actions"

	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"go.uber.org/zap"
)

var listenEvents = []ghwebhooks.Event{
	ghwebhooks.IssuesEvent,
	ghwebhooks.WorkflowDispatchEvent,
	ghwebhooks.WorkflowRunEvent,
	ghwebhooks.WorkflowRunEvent,
}

type Webhook struct {
	logger *zap.SugaredLogger
	cs     clients.ClientMap
	hook   *ghwebhooks.Webhook
	a      *actions.WebhookActions
}

// NewGithubWebhook returns a new webhook controller
func NewGithubWebhook(logger *zap.SugaredLogger, w config.Webhook, cs clients.ClientMap) (*Webhook, error) {
	hook, err := ghwebhooks.New(ghwebhooks.Options.Secret(os.Getenv(w.Secret)))
	if err != nil {
		return nil, err
	}

	a, err := actions.InitActions(logger, cs, w.Actions)

	if err != nil {
		return nil, err
	}

	controller := &Webhook{
		logger: logger,
		cs:     cs,
		hook:   hook,
		a:      a,
	}

	return controller, nil
}

// Handle handles GitHub webhook events
func (w *Webhook) Handle(response http.ResponseWriter, request *http.Request) {
	payload, err := w.hook.Parse(request, listenEvents...)

	if err != nil {
		if errors.Is(err, ghwebhooks.ErrEventNotFound) {
			//w.logger.Warnw("received unregistered github event", "error", err)
			response.WriteHeader(http.StatusOK)
		} else {
			w.logger.Errorw("received malformed github event", "error", err)
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	switch payload := payload.(type) {
	case ghwebhooks.WorkflowDispatchPayload:
		w.logger.Infow("received workflow dispatch event")
		go w.a.ProcessWorkflowDispatchEvent(&payload)
	case ghwebhooks.WorkflowRunPayload:
		w.logger.Infow("received workflow run event")
		go w.a.ProcessWorkflowRunEvent(&payload)
	case ghwebhooks.WorkflowJobPayload:
		w.logger.Infow("received workflow job event")
		go w.a.ProcessWorkflowJobEvent(&payload)
	case ghwebhooks.IssuesPayload:
		w.logger.Infow("received issues event")
		go w.a.ProcessIssuesEvent(&payload)
	default:
		w.logger.Warnw("missing handler", "payload", payload)
	}

	response.WriteHeader(http.StatusOK)
}
