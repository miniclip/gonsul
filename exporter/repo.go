package exporter

import (
	"github.com/miniclip/gonsul/errorutil"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"

	"errors"
	"fmt"
)

func downloadRepo(fileSystemPath string, url string) {
	// Check if SSH Key path was given
	var auth ssh.AuthMethod
	var sshUser = config.GetRepoSSHUser()
	var sshKey = config.GetRepoSSHKey()

	if sshUser != "" && sshKey != "" {
		auth, _ = ssh.NewPublicKeysFromFile(sshUser, sshKey, "")
	}

	repo, err := git.PlainClone(fileSystemPath, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              auth,
	})

	if err != nil {
		logger.PrintDebug(fmt.Sprintf("REPO: failed clone (%s), trying to open directory", err.Error()))

		// Cloning failed, most probably due to directory already cloned, moving to Open Dir
		repo, err = git.PlainOpen(config.GetRepoRootDir())

		if err != nil {
			errorutil.ExitError(
				errors.New("REPO: failed clone and directory is not a git repo, try cleaning dir"),
				errorutil.ErrorFailedCloning,
				&logger,
			)
		}

		logger.PrintDebug(fmt.Sprintf("REPO: git directory openned: %s", config.GetRepoRootDir()))
	}

	// We're still here, let's try to checkout required branch
	tryCheckout(repo, &auth)
}

func tryCheckout(repo *git.Repository, auth *ssh.AuthMethod) {
	// Initiate our worktree
	workTree, err := repo.Worktree()
	checkRepoError(err)

	// Get remotes, to check if current GIT is ours
	remotes, err := repo.Remotes()
	checkRepoError(err)

	// Check if remote is valid (the same as ours
	if !checkIfRemoteValid(remotes) {
		errorutil.ExitError(errors.New(fmt.Sprintf("REPO: remote url is not equal to provided: %s", config.GetRepoURL())), errorutil.ErrorFailedCloning, &logger)
	}

	logger.PrintDebug(fmt.Sprintf("REPO: pulling changes: %s", config.GetRepoBranch()))
	// We shall ignore error here, as Pull return messages such as "non-fast-forward update" as an error
	err = workTree.Pull(&git.PullOptions{
		RemoteName: config.GetRepoRemoteName(),
		Auth:       *auth,
	})
	// TODO: Even though the comment just above is true, we should handle this cases in a better way
	if err != nil {
		logger.PrintDebug(fmt.Sprintf("REPO: pull complete: %s", err.Error()))
	} else {
		logger.PrintDebug("REPO: pull complete")
	}

	logger.PrintDebug(fmt.Sprintf("REPO: checking out: %s", config.GetRepoBranch()))
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("refs/remotes/%s/%s", config.GetRepoRemoteName(), config.GetRepoBranch())),
		Create: false,
		Force:  true,
	})
	checkRepoError(err)
}

func checkIfRemoteValid(remotes []*git.Remote) bool {
	// Iterate over remotes
	for _, remote := range remotes {
		// Iterate over URLs
		for _, url := range remote.Config().URLs {
			// Compare current url with ours
			if url == config.GetRepoURL() {
				return true
			}
		}
	}

	return false
}

func checkRepoError(err error) {
	if err != nil {
		errorutil.ExitError(errors.New("REPO: "+err.Error()), errorutil.ErrorFailedCloning, &logger)
	}
}
