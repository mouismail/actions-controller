package github

import (
	"errors"
	"net/http"

	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
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
	hook, err := ghwebhooks.New(ghwebhooks.Options.Secret(w.Secret))
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
			w.logger.Warnw("received unregistered github event", "error", err)
			response.WriteHeader(http.StatusOK)
		} else {
			w.logger.Errorw("received malformed github event", "error", err)
			response.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	switch payload := payload.(type) {
	// TODO: IssueEvent
	case ghwebhooks.IssuesPayload:
		w.logger.Infow("received issues event")
	case ghwebhooks.WorkflowDispatchPayload:
		w.logger.Infow("received workflow dispatch event")
		w.logger.Infow("Organization", "organization", payload.Organization.Login)
		w.logger.Infow("Repository", "repository", payload.Repository.Name)
		w.logger.Infow("Sender", "sender", payload.Sender.Login)
		w.logger.Infow("Workflow", "workflow", payload.Workflow)
		go w.a.ProcessWorkflowDispatchEvent(&payload)
	case ghwebhooks.WorkflowRunPayload:
		w.logger.Infow("received workflow run event")
		w.logger.Infow("Organization", "organization", payload.Organization.Login)
		w.logger.Infow("Repository", "repository", payload.Repository.Name)
		w.logger.Infow("Sender", "sender", payload.Sender.Login)
		w.logger.Infow("Workflow", "workflow", payload.Workflow.Name)
		go w.a.ProcessWorkflowRunEvent(&payload)
	case ghwebhooks.WorkflowJobPayload:
		w.logger.Infow("received workflow job event")
		w.logger.Infow("Organization", "organization", payload.Organization.Login)
		w.logger.Infow("Repository", "repository", payload.Repository.Name)
		w.logger.Infow("Sender", "sender", payload.Sender.Login)
		w.logger.Infow("Workflow Job", "workflow_job", payload.WorkflowJob)
		go w.a.ProcessWorkflowJobEvent(&payload)
	default:
		w.logger.Warnw("missing handler", "payload", payload)
	}

	response.WriteHeader(http.StatusOK)
}
