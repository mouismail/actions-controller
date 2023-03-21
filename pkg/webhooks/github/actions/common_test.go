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

func TestWebhookActions_ProcessIssueCommentEvent(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		ih     []*IssuesAction
		wa     []*WorkflowAction
	}
	type args struct {
		payload *ghwebhooks.IssueCommentPayload
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WebhookActions{
				logger: tt.fields.logger,
				ih:     tt.fields.ih,
				wa:     tt.fields.wa,
			}
			w.ProcessIssueCommentEvent(tt.args.payload)
		})
	}
}

func TestWebhookActions_ProcessWorkflowDispatchEvent(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		ih     []*IssuesAction
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
				logger: tt.fields.logger,
				ih:     tt.fields.ih,
				wa:     tt.fields.wa,
			}
			w.ProcessWorkflowDispatchEvent(tt.args.payload)
		})
	}
}

func TestWebhookActions_ProcessWorkflowJobEvent(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		ih     []*IssuesAction
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
				logger: tt.fields.logger,
				ih:     tt.fields.ih,
				wa:     tt.fields.wa,
			}
			w.ProcessWorkflowJobEvent(tt.args.payload)
		})
	}
}

func TestWebhookActions_ProcessWorkflowRunEvent(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		ih     []*IssuesAction
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
				logger: tt.fields.logger,
				ih:     tt.fields.ih,
				wa:     tt.fields.wa,
			}
			w.ProcessWorkflowRunEvent(tt.args.payload)
		})
	}
}

func Test_extractTag(t *testing.T) {
	type args struct {
		payload *ghwebhooks.PushPayload
	}
	var tests []struct {
		name string
		args args
		want string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractTag(tt.args.payload); got != tt.want {
				t.Errorf("extractTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
