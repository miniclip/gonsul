package main

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/exporter"
	"github.com/miniclip/gonsul/importer"
	"github.com/miniclip/gonsul/data"
	"errors"
)

func start() {
	// Build our configuration
	config, err 	:= configuration.GetConfig()
	if err != nil {
		var logger = errorutil.NewLogger(0)
		errorutil.ExitError(err, errorutil.ErrorBadParams, logger)
	}
	logger 	:= errorutil.NewLogger(config.GetLogLevel())

	switch config.GetStrategy() {
	case configuration.StrategyOnce:
		startOnce(config, logger)

	case configuration.StrategyHook:
		startHook(config, logger)

	case configuration.StrategyPoll:
		startPolling(config, logger)

	}

	logger.PrintDebug("Quitting... bye 😀")
}

func startPolling(conf *configuration.Config, log *errorutil.Logger)  {
	/* TODO */
	errorutil.ExitError(errors.New("POLLING: NOT IMPLEMENTED YET"), 100, log)
}

func startHook(conf *configuration.Config, log *errorutil.Logger)  {
	/* TODO */
	errorutil.ExitError(errors.New("HOOK: NOT IMPLEMENTED YET"), 100, log)
}

func startOnce(conf *configuration.Config, log *errorutil.Logger) {
	log.PrintDebug("Starting in mode: ONCE")
	// Export our data
	datum := exportData(conf, log)
	// Start data import to Consul
	importData(datum, conf, log)
}

func exportData(conf *configuration.Config, log *errorutil.Logger) data.EntryCollection {
	log.PrintDebug("Starting data retrieve from GIT")
	processedData 	:= exporter.Export(conf, log)
	log.PrintDebug("Finished data retrieve from GIT")

	return processedData
}

func importData(data data.EntryCollection, conf *configuration.Config, log *errorutil.Logger) {
	log.PrintDebug("Starting data import to Consul")
	importer.Start(data, conf, log)
	log.PrintDebug("Finished data import to Consul")
}