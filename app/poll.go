package app

import (
	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/internal/util"

	"fmt"
	"time"
)

type Ipoll interface {
	RunPoll()
}

type poll struct {
	config     config.IConfig
	logger     util.ILogger
	once       Ionce
	iterations int
}

func NewPoll(config config.IConfig, logger util.ILogger, once Ionce, it int) Ipoll {
	return &poll{
		config:     config,
		logger:     logger,
		once:       once,
		iterations: it,
	}
}

func (a *poll) RunPoll() {
	a.logger.PrintInfo("Starting in mode: POLL")

	// Loop forever
	count := 1
	for {
		a.logger.PrintDebug(fmt.Sprintf("POLL: performing iteration %d", count))
		// Run our once step
		a.once.RunOnce()

		// Sleep for the amount of time in config
		time.Sleep(time.Second * time.Duration(a.config.GetPollInterval()))

		// Make sure we respect the give max iterations (zero means infinite loop)
		// NOTE: This is only useful for testing purposes
		if a.iterations > 0 && a.iterations == count {
			break
		}

		// Increment our iteration counter
		count++
	}
}
