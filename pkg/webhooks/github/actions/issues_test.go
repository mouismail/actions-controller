package actions

import (
	"context"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestIssuesAction_HandleIssueComment(t *testing.T) {
	type fields struct {
		logger      *zap.SugaredLogger
		client      *clients.Github
		targetRepos map[string]bool
	}
	type args struct {
		ctx context.Context
		p   *IssuesActionParams
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &IssuesAction{
				logger:      tt.fields.logger,
				client:      tt.fields.client,
				targetRepos: tt.fields.targetRepos,
			}
			if err := r.HandleIssueComment(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("HandleIssueComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIssuesAction_buildForkPR(t *testing.T) {
	type fields struct {
		logger      *zap.SugaredLogger
		client      *clients.Github
		targetRepos map[string]bool
	}
	type args struct {
		ctx context.Context
		p   *IssuesActionParams
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &IssuesAction{
				logger:      tt.fields.logger,
				client:      tt.fields.client,
				targetRepos: tt.fields.targetRepos,
			}
			if err := r.buildForkPR(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("buildForkPR() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewIssuesAction(t *testing.T) {
	type args struct {
		logger    *zap.SugaredLogger
		client    *clients.Github
		rawConfig map[string]any
	}
	var tests []struct {
		name    string
		args    args
		want    *IssuesAction
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewIssuesAction(tt.args.logger, tt.args.client, tt.args.rawConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewIssuesAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIssuesAction() got = %v, want %v", got, tt.want)
			}
		})
	}
}
