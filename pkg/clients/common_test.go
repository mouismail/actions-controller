package clients

import (
	"github.tools.sap/actions-rollout-app/pkg/config"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestInitClients(t *testing.T) {
	type args struct {
		logger *zap.SugaredLogger
		config []config.Client
	}
	var tests []struct {
		name    string
		args    args
		want    ClientMap
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InitClients(tt.args.logger, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitClients() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitClients() got = %v, want %v", got, tt.want)
			}
		})
	}
}
