package main

import (
	"github.com/miniclip/gonsul/app"
	"github.com/miniclip/gonsul/internal/configuration"
	"github.com/miniclip/gonsul/internal/exporter"
	"github.com/miniclip/gonsul/internal/importer"
	"github.com/miniclip/gonsul/internal/util"
	"net/http"
	"os"
	"time"
)

var AppVersion = ""

func main() {
	defer func() {
		if r := recover(); r != nil {
			var recoveredError = r.(util.GonsulError)
			os.Exit(recoveredError.Code)
		}
	}()

	start()
}

func start() {
	// Build our configuration
	config, err := configuration.NewConfig()
	if err != nil {
		util.ExitError(err, util.ErrorBadParams, util.NewLogger(0))
	}

	// Build our logger
	logger := util.NewLogger(config.GetLogLevel())

	// Build our application and all it's dependencies
	httpClient := &http.Client{Timeout: time.Second * time.Duration(config.GetTimeout())}
	imp := importer.NewImporter(config, logger, httpClient)
	exp := exporter.NewExporter(config, logger)
	application := app.NewApplication(config, logger, imp, exp)

	// Start our application
	application.Start()

	// We're still here, all went well, good bye
	logger.PrintInfo("Quitting... bye.")
}
