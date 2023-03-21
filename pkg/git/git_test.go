package git

import "testing"

func TestDeleteBranch(t *testing.T) {
	type args struct {
		repoURL string
		branch  string
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteBranch(tt.args.repoURL, tt.args.branch); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPushToRemote(t *testing.T) {
	type args struct {
		remoteURL    string
		remoteBranch string
		targetURL    string
		targetBranch string
		msg          string
	}
	var tests []struct {
		name    string
		args    args
		wantErr bool
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PushToRemote(tt.args.remoteURL, tt.args.remoteBranch, tt.args.targetURL, tt.args.targetBranch, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("PushToRemote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
