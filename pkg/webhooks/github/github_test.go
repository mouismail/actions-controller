package github

import (
	ghwebhooks "github.com/go-playground/webhooks/v6/github"
	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"github.tools.sap/actions-rollout-app/pkg/webhooks/github/actions"
	"go.uber.org/zap"
	"net/http"
	"reflect"
	"testing"
)

func TestNewGithubWebhook(t *testing.T) {
	type args struct {
		logger *zap.SugaredLogger
		w      config.Webhook
		cs     clients.ClientMap
	}
	var tests []struct {
		name    string
		args    args
		want    *Webhook
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGithubWebhook(tt.args.logger, tt.args.w, tt.args.cs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGithubWebhook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGithubWebhook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebhook_Handle(t *testing.T) {
	type fields struct {
		logger *zap.SugaredLogger
		cs     clients.ClientMap
		hook   *ghwebhooks.Webhook
		a      *actions.WebhookActions
	}
	type args struct {
		response http.ResponseWriter
		request  *http.Request
	}
	var tests []struct {
		name   string
		fields fields
		args   args
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Webhook{
				logger: tt.fields.logger,
				cs:     tt.fields.cs,
				hook:   tt.fields.hook,
				a:      tt.fields.a,
			}
			w.Handle(tt.args.response, tt.args.request)
		})
	}
}
