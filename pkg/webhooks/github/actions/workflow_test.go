package actions

import (
	"context"
	"github.tools.sap/actions-rollout-app/pkg/clients"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestNewWorkflowAction(t *testing.T) {
	type args struct {
		logger    *zap.SugaredLogger
		client    *clients.Github
		rawConfig map[string]any
	}
	var tests []struct {
		name    string
		args    args
		want    *WorkflowAction
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewWorkflowAction(tt.args.logger, tt.args.client, tt.args.rawConfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewWorkflowAction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWorkflowAction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkflowAction_HandleWorkflow(t *testing.T) {
	type fields struct {
		logger       *zap.SugaredLogger
		client       *clients.Github
		repository   string
		organization string
		workflowId   int64
	}
	type args struct {
		ctx context.Context
		p   *WorkflowActionParams
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WorkflowAction{
				logger:       tt.fields.logger,
				client:       tt.fields.client,
				repository:   tt.fields.repository,
				organization: tt.fields.organization,
				workflowId:   tt.fields.workflowId,
			}
			if err := w.HandleWorkflow(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("HandleWorkflow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkflowAction_handleWorkflowDispatch(t *testing.T) {
	type fields struct {
		logger       *zap.SugaredLogger
		client       *clients.Github
		repository   string
		organization string
		workflowId   int64
	}
	type args struct {
		ctx context.Context
		p   *WorkflowActionParams
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WorkflowAction{
				logger:       tt.fields.logger,
				client:       tt.fields.client,
				repository:   tt.fields.repository,
				organization: tt.fields.organization,
				workflowId:   tt.fields.workflowId,
			}
			if err := w.handleWorkflowDispatch(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("handleWorkflowDispatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkflowAction_handleWorkflowJob(t *testing.T) {
	type fields struct {
		logger       *zap.SugaredLogger
		client       *clients.Github
		repository   string
		organization string
		workflowId   int64
	}
	type args struct {
		ctx context.Context
		p   *WorkflowActionParams
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WorkflowAction{
				logger:       tt.fields.logger,
				client:       tt.fields.client,
				repository:   tt.fields.repository,
				organization: tt.fields.organization,
				workflowId:   tt.fields.workflowId,
			}
			if err := w.handleWorkflowJob(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("handleWorkflowJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWorkflowAction_handleWorkflowRun(t *testing.T) {
	type fields struct {
		logger       *zap.SugaredLogger
		client       *clients.Github
		repository   string
		organization string
		workflowId   int64
	}
	type args struct {
		ctx context.Context
		p   *WorkflowActionParams
	}
	var tests []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &WorkflowAction{
				logger:       tt.fields.logger,
				client:       tt.fields.client,
				repository:   tt.fields.repository,
				organization: tt.fields.organization,
				workflowId:   tt.fields.workflowId,
			}
			if err := w.handleWorkflowRun(tt.args.ctx, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("handleWorkflowRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
