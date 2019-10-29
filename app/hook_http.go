package app

import (
	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/internal/util"

	"errors"
	"net/http"
)

// IHookHttp is our interface used for the Hook sgtartegy
type IHookHttp interface {
	Start(route string, handler func(http.ResponseWriter, *http.Request))
}

// hookHttp is our IHookHttp concrete implementation
type hookHttp struct {
	config config.IConfig
	logger util.ILogger
}

// NewHookHttp is our hookHttp constructor
func NewHookHttp(config config.IConfig, logger util.ILogger) IHookHttp {
	return &hookHttp{config: config, logger: logger}
}

// Start starts our HTTP server
func (h *hookHttp) Start(route string, handler func(http.ResponseWriter, *http.Request)) {
	// Create our routes and set handlers
	http.HandleFunc(route, handler)

	// Launch our HTTP server
	if err := http.ListenAndServe(":8000", nil); err != nil {
		util.ExitError(errors.New("Hook: "+err.Error()), util.ErrorFailedHTTPServer, h.logger)
	}
}
