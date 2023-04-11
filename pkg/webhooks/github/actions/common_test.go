package actions

import (
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/config"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestInitActions(t *testing.T) {
	type args struct {
		logger *zap.SugaredLogger
		cs     clients.ClientMap
		config config.WebhookActions
	}
	var tests []struct {
		name    string
		args    args
		want    *WebhookActions
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitActions(tt.args.logger, tt.args.cs, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitActions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitActions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebhookActions_ProcessWorkflowDispatchEvent(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		wa     []*WorkflowAction
	}
	type args struct {
		payload *ghwebhooks.WorkflowDispatchPayload
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebhookActions{
				logger:          tt.fields.logger,
				workflowActions: tt.fields.wa,
			}
			w.ProcessWorkflowDispatchEvent(tt.args.payload)
		})
	}
}

func TestWebhookActions_ProcessWorkflowJobEvent(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		wa     []*WorkflowAction
	}
	type args struct {
		payload *ghwebhooks.WorkflowJobPayload
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebhookActions{
				logger:          tt.fields.logger,
				workflowActions: tt.fields.wa,
			}
			w.ProcessWorkflowJobEvent(tt.args.payload)
		})
	}
}

func TestWebhookActions_ProcessWorkflowRunEvent(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		wa     []*WorkflowAction
	}
	type args struct {
		payload *ghwebhooks.WorkflowRunPayload
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebhookActions{
				logger:          tt.fields.logger,
				workflowActions: tt.fields.wa,
			}
			w.ProcessWorkflowRunEvent(tt.args.payload)
		})
	}
}
