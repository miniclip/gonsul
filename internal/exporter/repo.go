package exporter

import (
	"github.com/miniclip/gonsul/internal/util"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"errors"
	"fmt"
)

// downloadRepo ...
func (e *exporter) downloadRepo() {
	// Get some variables
	var (
		fileSystemPath = e.config.GetRepoRootDir()
		url            = e.config.GetRepoURL()
		sshUser        = e.config.GetRepoSSHUser()
		sshKey         = e.config.GetRepoSSHKey()
		auth           ssh.AuthMethod
	)

	// Check if SSH Key path was given
	if sshUser != "" && sshKey != "" {
		auth, _ = ssh.NewPublicKeysFromFile(sshUser, sshKey, "")
	}

	// Clone given repository
	repo, err := git.PlainClone(fileSystemPath, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
	})

	if err != nil {
		e.logger.PrintDebug(fmt.Sprintf("REPO: failed clone (%s), trying to open directory", err.Error()))

		// Cloning failed, most probably due to directory already cloned, moving to Open Dir
		repo, err = git.PlainOpen(e.config.GetRepoRootDir())

		if err != nil {
			util.ExitError(
				errors.New("REPO: failed clone and directory is not a git repo, try cleaning dir"),
				util.ErrorFailedCloning,
				e.logger,
			)
		}

		e.logger.PrintDebug(fmt.Sprintf("REPO: git directory opened: %s", e.config.GetRepoRootDir()))
	}

	// We're still here, let's try to checkout required branch
	e.tryCheckout(repo, &auth)
}

// tryCheckout ...
func (e *exporter) tryCheckout(repo *git.Repository, auth *ssh.AuthMethod) {
	// Initiate our worktree
	workTree, err := repo.Worktree()
	e.checkRepoError(err)

	// Get remotes, to check if current GIT is ours
	remotes, err := repo.Remotes()
	e.checkRepoError(err)

	// Check if remote is valid (the same as ours
	if !e.checkIfRemoteValid(remotes) {
		util.ExitError(
			errors.New(fmt.Sprintf("REPO: remote url is not equal to provided: %s", e.config.GetRepoURL())),
			util.ErrorFailedCloning,
			e.logger,
		)
	}

	e.logger.PrintDebug(fmt.Sprintf("REPO: pulling changes: %s", e.config.GetRepoBranch()))
	// We shall ignore error here, as Pull return messages such as "non-fast-forward update" as an error
	err = workTree.Pull(&git.PullOptions{
		RemoteName: e.config.GetRepoRemoteName(),
		Auth:       *auth,
	})
	// TODO: Even though the comment just above is true, we should handle this cases in a better way
	if err != nil {
		e.logger.PrintDebug(fmt.Sprintf("REPO: pull complete: %s", err.Error()))
	} else {
		e.logger.PrintDebug("REPO: pull complete")
	}

	e.logger.PrintDebug(fmt.Sprintf("REPO: checking out: %s", e.config.GetRepoBranch()))
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/remotes/%s/%s", e.config.GetRepoRemoteName(), e.config.GetRepoBranch())),
		Create: false,
		Force:  true,
	})
	e.checkRepoError(err)
}

// checkIfRemoteValid ...
func (e *exporter) checkIfRemoteValid(remotes []*git.Remote) bool {
	// Iterate over remotes
	for _, remote := range remotes {
		// Iterate over URLs
		for _, url := range remote.Config().URLs {
			// Compare current url with ours
			if url == e.config.GetRepoURL() {
				return true
			}
		}
	}

	return false
}

// checkRepoError ...
func (e *exporter) checkRepoError(err error) {
	if err != nil {
		util.ExitError(errors.New("REPO: "+err.Error()), util.ErrorFailedCloning, e.logger)
	}
}
