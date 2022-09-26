package app

import (
	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/internal/importer"
	"github.com/miniclip/gonsul/internal/util"
)

type Iread interface {
	RunRead()
}

type read struct {
	config   config.IConfig
	logger   util.ILogger
	importer importer.IImporter
}

func NewRead(config config.IConfig, logger util.ILogger, importer importer.IImporter) Iread {
	return &read{
		config:   config,
		logger:   logger,
		importer: importer,
	}
}

// RunOnce is our entry point function for the Once Application mode
func (a *read) RunRead() {
	strategy := a.config.GetStrategy()
	// Check strategy
	if strategy == config.StrategyRead {
		a.logger.PrintInfo("Starting in mode: READ")
	} else {
		a.logger.PrintError("Bug on strategy READ")
		return
	}

	var localData = map[string]string{}
	a.logger.PrintDebug("Starting reading data fromConsul")
	a.importer.Start(localData)
	a.logger.PrintDebug("Finished reading data from Consul")
}
