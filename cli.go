package main

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/exporter"
	"github.com/miniclip/gonsul/importer"
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
	case configuration.StrategyDry:
		startOnce(config, logger)

	case configuration.StrategyOnce:
		startOnce(config, logger)

	case configuration.StrategyHook:
		startHook(config, logger)

	case configuration.StrategyPoll:
		startPolling(config, logger)

	}

	logger.PrintInfo("Quitting... bye 😀")
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
	if conf.GetStrategy() == configuration.StrategyDry {
		log.PrintInfo("Starting in mode: DRYRUN")
	} else if conf.GetStrategy() == configuration.StrategyOnce {
		log.PrintInfo("Starting in mode: ONCE")
	}
	// Export our data
	log.PrintDebug("Starting data retrieve from GIT")
	processedData 	:= exporter.Export(conf, log)
	log.PrintDebug("Finished data retrieve from GIT")

	// Start data import to Consul
	log.PrintDebug("Starting data import to Consul")
	importer.Start(processedData, conf, log)
	log.PrintDebug("Finished data import to Consul")
}