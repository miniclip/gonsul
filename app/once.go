package app

import (
	"github.com/miniclip/gonsul/internal/configuration"
)

func (a *Application) once() {
	// Check strategy
	if a.config.GetStrategy() == configuration.StrategyDry {
		a.logger.PrintInfo("Starting in mode: DRYRUN")
	} else if a.config.GetStrategy() == configuration.StrategyOnce {
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
