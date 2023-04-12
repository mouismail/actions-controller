package webhooks

import (
	"github.tools.sap/actions-rollout-app/config"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"go.uber.org/zap"
	"testing"
)

func TestInitWebhooks(t *testing.T) {
	type args struct {
		logger *zap.SugaredLogger
		cs     clients.ClientMap
		c      *config.Configuration
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitWebhooks(tt.args.logger, tt.args.cs, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("InitWebhooks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
