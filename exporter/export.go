package exporter

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/data"
)

var config configuration.Config
var logger errorutil.Logger

func Export(conf *configuration.Config, log *errorutil.Logger) data.EntryCollection {
	// Set the appropriate values for our package global variables
	config = *conf
	logger = *log

	// Instantiate our import data structure
	var processedData data.EntryCollection
	// Set the path where Gonsul should start traversing files to add to Consul
	repoDir 	:= config.GetRepoRootDir() + "/" + config.GetRepoBasePath()

	// Should we clone the repo, or is it already done via 3rd party
	if config.IsCloning() {
		logger.PrintInfo("REPO: GIT cloning from: " + config.GetRepoURL())
		downloadRepo(config.GetRepoRootDir(), config.GetRepoURL())
	} else {
		logger.PrintInfo("REPO: Skipping GIT clone, using local path: " + config.GetRepoRootDir())
	}
	// Traverse our repo directory, filling up the data.EntryCollection structure
	processDir(repoDir, &processedData)

	// Return our final data.EntryCollection structure
	return processedData
}