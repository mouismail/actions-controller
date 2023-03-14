package repos_test

import (
	"testing"

	"github.tools.sap/actions-rollout-app/repos"
)

func TestLoadNoRepoConfig(t *testing.T) {
	want := repos.ActionConfig{
		URL:          "https://github.example.com/octo",
		ContactEmail: "octodemo@example.com",
		UseCase:      "octodemo",
	}

	got, err := repos.LoadConfig("./fixtures/none.yml")

	if err != nil {
		t.Errorf("got %v, wanted %v", err, nil)
	}

	if got.URL != want.URL {
		t.Errorf("got %v, wanted %v", got.URL, want.URL)
	}

	if got.ContactEmail != want.ContactEmail {
		t.Errorf("got %v, wanted %v", got.ContactEmail, want.ContactEmail)
	}

	if got.UseCase != want.UseCase {
		t.Errorf("got %v, wanted %v", got.UseCase, want.UseCase)
	}

	if len(got.Repos) != 0 {
		t.Errorf("got %v repos, wanted %v", len(got.Repos), 0)
	}
}

func TestLoadOneRepoConfig(t *testing.T) {
	want := repos.ActionConfig{
		URL:          "https://github.example.com/octo",
		ContactEmail: "octodemo@example.com",
		UseCase:      "octodemo",
		Repos: []string{
			"https://github.tools.sap/octo/demo1",
		},
	}

	got, err := repos.LoadConfig("./fixtures/one.yml")

	if err != nil {
		t.Errorf("got %v, wanted %v", err, nil)
	}

	if len(got.Repos) != 1 {
		t.Errorf("got %v repo, wanted %v", len(got.Repos), 1)
	}

	if got.Repos[0] != want.Repos[0] {
		t.Errorf("got %v, wanted %v", got.Repos[0], want.Repos[0])
	}
}

func TestLoadMultiRepoConfig(t *testing.T) {
	want := repos.ActionConfig{
		URL:          "https://github.example.com/octo",
		ContactEmail: "octodemo@example.com",
		UseCase:      "octodemo",
		Repos: []string{
			"https://github.tools.sap/octo/demo1",
			"https://github.tools.sap/octo/demo2",
			"https://github.tools.sap/octo/demo3",
			"https://github.tools.sap/octo/demo4",
		},
	}

	got, err := repos.LoadConfig("./fixtures/multiple.yml")

	if err != nil {
		t.Errorf("got %v, wanted %v", err, nil)
	}

	if len(got.Repos) != 4 {
		t.Errorf("got %v repo, wanted %v", len(got.Repos), 4)
	}

	for i, repo := range got.Repos {
		if repo != want.Repos[i] {
			t.Errorf("got %v, wanted %v", repo, want.Repos[i])
		}
	}
}
