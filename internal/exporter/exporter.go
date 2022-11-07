package exporter

import (
	"errors"
	"fmt"
	"os"

	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/internal/util"

	"path"
)

// IExporter ...
type IExporter interface {
	Start() map[string]string
}

// exporter ...
type exporter struct {
	config config.IConfig
	logger util.ILogger
}

// NewExporter ...
func NewExporter(config config.IConfig, logger util.ILogger) IExporter {
	return &exporter{config: config, logger: logger}
}

// Start ...
func (e *exporter) Start() map[string]string {
	// Instantiate our local data map
	var localData = map[string]string{}

	// Should we clone the repo, or is it already done via 3rd party
	if e.config.IsCloning() {
		e.logger.PrintInfo("EXPORTER: Git cloning from configured remote repository")
		e.downloadRepo()
	} else {
		e.logger.PrintInfo("EXPORTER: Skipping Git clone, using local path: " + e.config.GetRepoRootDir())
	}

	// Set the path where Gonsul should start traversing files to add to Consul
	repoDir := path.Join(e.config.GetRepoRootDir(), e.config.GetRepoBasePath())
	// Traverse our repo directory, filling up the data.EntryCollection structure

	for _, fileName := range e.config.GetPaths() {
		filePath := path.Join(repoDir, fileName)
		fmt.Printf("Importing: %s\n", filePath)
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			util.ExitError(errors.New(err.Error()), util.ErrorRead, e.logger)
		}
		if fileInfo.IsDir() {
			e.parseDir(filePath, localData)
		} else {
			e.loadDictFile(filePath, localData, false)
		}
	}

	// Return our final data.EntryCollection structure
	return localData
}
