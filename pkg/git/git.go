package git

// TODO: remove not part of the requirements

import (
	"fmt"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.tools.sap/actions-rollout-app/pkg/utils"
)

var NoChangesError = fmt.Errorf("no changes")

func PushToRemote(remoteURL, remoteBranch, targetURL, targetBranch, msg string) error {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		RemoteName:    "remote-repo",
		URL:           remoteURL,
		ReferenceName: plumbing.ReferenceName(utils.DefaultLocalRef + "/" + remoteBranch),
	})
	if err != nil {
		return fmt.Errorf("error cloning git repo %w", err)
	}

	remote, err := r.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{targetURL},
	})
	if err != nil {
		return fmt.Errorf("error creating remote %w", err)
	}

	err = remote.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(utils.DefaultLocalRef + "/" + remoteBranch + ":" + utils.DefaultLocalRef + "/" + targetBranch),
		},
		Force: true, // when the contributor does a force push, this will make it work anyway
	})
	if err != nil {
		return fmt.Errorf("error pushing to repo %w", err)
	}

	return nil
}

func DeleteBranch(repoURL, branch string) error {
	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:   repoURL,
		Depth: 1,
	})
	if err != nil {
		return fmt.Errorf("error cloning git repo %w", err)
	}

	err = r.Storer.RemoveReference(plumbing.NewBranchReferenceName(branch))
	if err != nil {
		return fmt.Errorf("error deleting branch in git repo %w", err)
	}

	return nil
}

//func ReadConfigFile(repoURL, branch string) error {
//	r, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
//		URL:   repoURL,
//		Depth: 1,
//	})
//	if err != nil {
//		return fmt.Errorf("error cloning git repo %w", err)
//	}
//
//	_, err = r.Worktree()
//
//	if err != nil {
//		return fmt.Errorf("error checking out branch %w", err)
//	}
//
//	return nil
//}
