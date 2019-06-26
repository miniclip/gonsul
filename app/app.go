package app

import (
	"github.com/miniclip/gonsul/internal/configuration"
	"github.com/miniclip/gonsul/internal/exporter"
	"github.com/miniclip/gonsul/internal/importer"
	"github.com/miniclip/gonsul/internal/util"
	"sync"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Application ...
type Application struct {
	mutex    *sync.Mutex
	config   configuration.IConfig
	logger   util.ILogger
	importer importer.IImporter
	exporter exporter.IExporter
}

// NewApplication ...
func NewApplication(config configuration.IConfig, logger util.ILogger, im importer.IImporter, ex exporter.IExporter) *Application {
	return &Application{
		mutex:    &sync.Mutex{},
		config:   config,
		logger:   logger,
		importer: im,
		exporter: ex,
	}
}

// Start ...
func (a *Application) Start() {
	// Create our channel for the Signal and relay Signal Notify to it
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	// Spin a routine to wait for a Signal
	go func() {
		// Wait for a signal through the channel
		<-sigChannel
		// Try to write to working channel (thus waiting for any in progress non interruptible work)
		a.config.WorkingChan() <- false
		// Exit
		fmt.Print(" Interrupt received... Quitting!")
		os.Exit(0)
	}()

	// Switch our run strategy
	switch a.config.GetStrategy() {
	case configuration.StrategyDry, configuration.StrategyOnce:
		a.once()
	case configuration.StrategyHook:
		a.hook()
	case configuration.StrategyPoll:
		a.poll()
	}
}
