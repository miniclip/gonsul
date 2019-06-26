package app

import (
	"fmt"
	"time"
)

func (a *Application) poll() {
	a.logger.PrintInfo("Starting in mode: POLL")

	// Loop forever
	count := 1
	for {
		a.logger.PrintDebug(fmt.Sprintf("POLL: performing iteration %d", count))
		// Run our once step
		a.once()

		// Sleep for the amount of time in config
		time.Sleep(time.Second * time.Duration(a.config.GetPollInterval()))

		// Increment our iteration counter
		count++
	}
}
