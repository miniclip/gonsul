package poll

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/once"
	"time"
)

var config 		configuration.Config 		// Set our Configuration as global package scope
var logger 		errorutil.Logger     		// Set our Logger as global package scope

func Start(conf *configuration.Config, log *errorutil.Logger) {
	// Set the appropriate values for our package global variables
	config 		= *conf
	logger 		= *log

	loop()
}

func loop() {
	// Forever
	count := 1
	for {
		logger.PrintDebug("POLL: performing iteration - " + string(count))
		// Run our once step
		once.Start(&config, &logger)

		// Sleep for the amount of time in Config
		time.Sleep(time.Second * time.Duration(config.GetPollInterval()))

		// Increment our iteration counter
		count++
	}
}