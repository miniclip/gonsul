package exporter

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
)

var config configuration.Config
var logger errorutil.Logger

func Start(conf *configuration.Config, log *errorutil.Logger) map[string]string {
	// Set the appropriate values for our package global variables
	config = *conf
	logger = *log

	// Instantiate our local data map
	var localData = map[string]string{}

	// Set the path where Gonsul should start traversing files to add to Consul
	repoDir := config.GetRepoRootDir() + "/" + config.GetRepoBasePath()

	// Should we clone the repo, or is it already done via 3rd party
	if config.IsCloning() {
		logger.PrintInfo("REPO: GIT cloning from: " + config.GetRepoURL())
		downloadRepo(config.GetRepoRootDir(), config.GetRepoURL())
	} else {
		logger.PrintInfo("REPO: Skipping GIT clone, using local path: " + config.GetRepoRootDir())
	}
	// Traverse our repo directory, filling up the data.EntryCollection structure
	processDir(repoDir, localData)

	// Return our final data.EntryCollection structure
	return localData
}
