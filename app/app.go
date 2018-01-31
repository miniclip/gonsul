package app

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/hook"
	"github.com/miniclip/gonsul/once"
	"github.com/miniclip/gonsul/poll"

	"os/signal"
	"syscall"
	"fmt"
	"os"
)

var config configuration.Config // Set our Configuration as global package scope
var logger errorutil.Logger     // Set our Logger as global package scope

func Start(conf *configuration.Config, log *errorutil.Logger) {
	// Set the appropriate values for our package global variables
	config = *conf
	logger = *log

	// Create our channel for the Signal and relay Signal Notify to it
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Spin a routine to wait for a Signal
	go func() {
		// Wait for a signal through the channel
		<-sigChannel
		// Try to write to working channel (thus waiting for any in progress non interruptible work)
		config.Working <- false
		// Exit
		fmt.Print(" Interrupt received... Quitting!")
		os.Exit(0)
	}()

	// Switch our run strategy
	switch config.GetStrategy() {
	case configuration.StrategyDry, configuration.StrategyOnce:
		startOnce()

	case configuration.StrategyHook:
		startHook()

	case configuration.StrategyPoll:
		startPolling()

	}
}

func startPolling() {
	logger.PrintInfo("Starting in mode: POLL")

	poll.Start(&config, &logger)
}

func startHook() {
	logger.PrintInfo("Starting in mode: HOOK")

	hook.Start(&config, &logger)
}

func startOnce() {
	if config.GetStrategy() == configuration.StrategyDry {
		logger.PrintInfo("Starting in mode: DRYRUN")
	} else if config.GetStrategy() == configuration.StrategyOnce {
		logger.PrintInfo("Starting in mode: ONCE")
	}

	once.Start(&config, &logger)
}
