package app

import (
	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/internal/exporter"
	"github.com/miniclip/gonsul/internal/importer"
	"github.com/miniclip/gonsul/internal/util"
)

type Ionce interface {
	RunOnce()
}

type once struct {
	config   config.IConfig
	logger   util.ILogger
	exporter exporter.IExporter
	importer importer.IImporter
}

func NewOnce(config config.IConfig, logger util.ILogger, exporter exporter.IExporter, importer importer.IImporter) Ionce {
	return &once{
		config:   config,
		logger:   logger,
		exporter: exporter,
		importer: importer,
	}
}

// RunOnce is our entry point function for the Once Application mode
func (a *once) RunOnce() {
	strategy := a.config.GetStrategy()
	// Check strategy
	if strategy == config.StrategyDry {
		a.logger.PrintInfo("Starting in mode: DRYRUN")
	} else if strategy == config.StrategyOnce {
		a.logger.PrintInfo("Starting in mode: ONCE")
	}

	// Start our data export
	a.logger.PrintDebug("Starting data retrieve from GIT")
	exportedData := a.exporter.Start()
	a.logger.PrintDebug("Finished data retrieve from GIT")

	// Start data import to Consul
	a.logger.PrintDebug("Starting data import to Consul")
	a.importer.Start(exportedData)
	a.logger.PrintDebug("Finished data import to Consul")
}
