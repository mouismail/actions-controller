package clients

import (
	"context"
	"reflect"
	"testing"

	"github.tools.sap/actions-rollout-app/pkg/config"

	"github.com/bradleyfalzon/ghinstallation/v2"
	v3 "github.com/google/go-github/v50/github"
	"go.uber.org/zap"
)

func TestGithub_GetV3AppClient(t *testing.T) {
	type fields struct {
		logger         *zap.SugaredLogger
		keyPath        string
		appID          int64
		installationID int64
		organizationID string
		atr            *ghinstallation.AppsTransport
		itr            *ghinstallation.Transport
		serverInfo     *config.ServerInfo
	}
	var tests []struct {
		name   string
		fields fields
		want   *v3.Client
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Github{
				logger:         tt.fields.logger,
				keyPath:        tt.fields.keyPath,
				appID:          tt.fields.appID,
				installationID: tt.fields.installationID,
				organizationID: tt.fields.organizationID,
				atr:            tt.fields.atr,
				itr:            tt.fields.itr,
				serverInfo:     tt.fields.serverInfo,
			}
			if got := a.GetV3AppClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetV3AppClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithub_GetV3Client(t *testing.T) {
	type fields struct {
		logger         *zap.SugaredLogger
		keyPath        string
		appID          int64
		installationID int64
		organizationID string
		atr            *ghinstallation.AppsTransport
		itr            *ghinstallation.Transport
		serverInfo     *config.ServerInfo
	}
	var tests []struct {
		name   string
		fields fields
		want   *v3.Client
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Github{
				logger:         tt.fields.logger,
				keyPath:        tt.fields.keyPath,
				appID:          tt.fields.appID,
				installationID: tt.fields.installationID,
				organizationID: tt.fields.organizationID,
				atr:            tt.fields.atr,
				itr:            tt.fields.itr,
				serverInfo:     tt.fields.serverInfo,
			}
			if got := a.GetV3Client(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetV3Client() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithub_GitToken(t *testing.T) {
	type fields struct {
		logger         *zap.SugaredLogger
		keyPath        string
		appID          int64
		installationID int64
		organizationID string
		atr            *ghinstallation.AppsTransport
		itr            *ghinstallation.Transport
		serverInfo     *config.ServerInfo
	}
	type args struct {
		ctx context.Context
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Github{
				logger:         tt.fields.logger,
				keyPath:        tt.fields.keyPath,
				appID:          tt.fields.appID,
				installationID: tt.fields.installationID,
				organizationID: tt.fields.organizationID,
				atr:            tt.fields.atr,
				itr:            tt.fields.itr,
				serverInfo:     tt.fields.serverInfo,
			}
			got, err := a.GitToken(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GitToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GitToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithub_Organization(t *testing.T) {
	type fields struct {
		logger         *zap.SugaredLogger
		keyPath        string
		appID          int64
		installationID int64
		organizationID string
		atr            *ghinstallation.AppsTransport
		itr            *ghinstallation.Transport
		serverInfo     *config.ServerInfo
	}
	var tests []struct {
		name   string
		fields fields
		want   string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Github{
				logger:         tt.fields.logger,
				keyPath:        tt.fields.keyPath,
				appID:          tt.fields.appID,
				installationID: tt.fields.installationID,
				organizationID: tt.fields.organizationID,
				atr:            tt.fields.atr,
				itr:            tt.fields.itr,
				serverInfo:     tt.fields.serverInfo,
			}
			if got := a.Organization(); got != tt.want {
				t.Errorf("Organization() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithub_initClients(t *testing.T) {
	type fields struct {
		logger         *zap.SugaredLogger
		keyPath        string
		appID          int64
		installationID int64
		organizationID string
		atr            *ghinstallation.AppsTransport
		itr            *ghinstallation.Transport
		serverInfo     *config.ServerInfo
	}
	var tests []struct {
		name    string
		fields  fields
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Github{
				logger:         tt.fields.logger,
				keyPath:        tt.fields.keyPath,
				appID:          tt.fields.appID,
				installationID: tt.fields.installationID,
				organizationID: tt.fields.organizationID,
				atr:            tt.fields.atr,
				itr:            tt.fields.itr,
				serverInfo:     tt.fields.serverInfo,
			}
			if err := a.initClients(); (err != nil) != tt.wantErr {
				t.Errorf("initClients() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewGithub(t *testing.T) {
	type args struct {
		logger         *zap.SugaredLogger
		organizationID string
		repository     string
		severInfo      *config.ServerInfo
		config         *config.GithubClient
	}
	var tests []struct {
		name    string
		args    args
		want    *Github
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGithub(tt.args.logger, tt.args.organizationID, tt.args.repository, tt.args.severInfo, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGithub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGithub() got = %v, want %v", got, tt.want)
			}
		})
	}
}
