package main

import (
	"github.com/miniclip/gonsul/app"
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"

	"os"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			var recoveredError = r.(errorutil.GonsulError)
			os.Exit(recoveredError.Code)
		}
	}()

	bootstrap()
}

func bootstrap() {
	// Build our configuration
	config, err := configuration.GetConfig()
	if err != nil {
		var logger = errorutil.NewLogger(0)
		errorutil.ExitError(err, errorutil.ErrorBadParams, logger)
	}

	// Build our logger
	logger := errorutil.NewLogger(config.GetLogLevel())

	// Start our application
	app.Start(config, logger)

	// We're still here, all went well, good bye
	logger.PrintInfo("Quitting... bye.")
}
