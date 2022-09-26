package app

import (
	"github.com/miniclip/gonsul/internal/config"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Application ...
type Application struct {
	config  config.IConfig
	once    Ionce
	read    Iread
	hook    Ihook
	poll    Ipoll
	sigChan chan os.Signal
}

// NewApplication ...
func NewApplication(
	config config.IConfig,
	once Ionce,
	read Iread,
	hook Ihook,
	poll Ipoll,
	sigChan chan os.Signal,
) *Application {
	return &Application{
		config:  config,
		once:    once,
		read:    read,
		hook:    hook,
		poll:    poll,
		sigChan: sigChan,
	}
}

// Start ...
func (a *Application) Start() {
	// Relay all Signals to our channel
	signal.Notify(a.sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Spin a routine to wait for a Signal
	go func() {
		// Wait for a signal through the channel
		<-a.sigChan
		fmt.Print(" Interrupt received, waiting for work to finish... ")
		// Try to write to working channel (thus waiting for any in progress non interruptible work)
		a.config.WorkingChan() <- false
		// Exit
		fmt.Print(" Quitting!")
		os.Exit(0)
	}()

	// Switch our run strategy
	switch a.config.GetStrategy() {
	case config.StrategyRead:
		a.read.RunRead()
	case config.StrategyDry, config.StrategyOnce:
		a.once.RunOnce()
	case config.StrategyHook:
		a.hook.RunHook()
	case config.StrategyPoll:
		a.poll.RunPoll()
	}
}
