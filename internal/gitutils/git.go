package gitutils

import (
	"errors"
	"fmt"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

func InitGitRepo(path string, remoteURL string, branch string) error {
	repo, err := git.PlainInit(path, false)
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return err
	}

	if remoteURL != "" {
		_, err = repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{remoteURL},
		})
		if err != nil && !errors.Is(err, git.ErrRemoteExists) {
			return err
		}
	}

	// Optional: create initial commit if repo is empty?
	return nil
}

func OpenRepo(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func GetCurrentCommit(path string) (string, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return "", err
	}

	ref, err := repo.Head()
	if err != nil {
		return "", err
	}

	return ref.Hash().String(), nil
}

func GetRemoteCommit(path, remoteName, branch string) (string, error) {
	repo, err := OpenRepo(path)
	if err != nil {
		return "", err
	}

	// Fetch from remote
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		Progress:   nil,
		Force:      true,
		Tags:       git.AllTags,
		// Remove error if already up-to-date
		// https://github.com/go-git/go-git/issues/287
	})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) && !errors.Is(err, transport.ErrAuthenticationRequired) {
		return "", fmt.Errorf("fetch error: %w", err)
	}

	// Resolve remote reference
	remoteRef := plumbing.NewRemoteReferenceName(remoteName, branch)
	ref, err := repo.Reference(remoteRef, true)
	if err != nil {
		return "", fmt.Errorf("remote ref not found: %w", err)
	}

	return ref.Hash().String(), nil
}

func ResetToRemote(path, remoteName, branch string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	remoteRef := plumbing.NewRemoteReferenceName(remoteName, branch)
	ref, err := repo.Reference(remoteRef, true)
	if err != nil {
		return err
	}

	// Reset worktree to remote HEAD
	return worktree.Reset(&git.ResetOptions{
		Mode:   git.HardReset,
		Commit: ref.Hash(),
	})
}

func CommitAndPush(path, message string) error {
	repo, err := OpenRepo(path)
	if err != nil {
		return err
	}

	wt, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Add all changes
	err = wt.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}

	_, err = wt.Commit(message, &git.CommitOptions{})
	if err != nil {
		return err
	}

	// Push to origin
	err = repo.Push(&git.PushOptions{})
	if err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return err
	}

	return nil
}
