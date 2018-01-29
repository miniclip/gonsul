package once

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/importer"
	"github.com/miniclip/gonsul/exporter"
)

var config 		configuration.Config 		// Set our Configuration as global package scope
var logger 		errorutil.Logger     		// Set our Logger as global package scope

func Start(conf *configuration.Config, log *errorutil.Logger) {
	// Set the appropriate values for our package global variables
	config 		= *conf
	logger 		= *log

	// Start our data
	logger.PrintDebug("Starting data retrieve from GIT")
	processedData 	:= exporter.Start(&config, &logger)
	logger.PrintDebug("Finished data retrieve from GIT")

	// Start data import to Consul
	logger.PrintDebug("Starting data import to Consul")
	importer.Start(processedData, &config, &logger)
	logger.PrintDebug("Finished data import to Consul")
}