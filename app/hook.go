package app

import (
	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/internal/util"
	"sync"

	"fmt"
	"net/http"
	"strings"
)

type Ihook interface {
	RunHook()
}

type hook struct {
	mutex  sync.Mutex
	http   IHookHttp
	config config.IConfig
	logger util.ILogger
	once   Ionce
}

func NewHook(http IHookHttp, config config.IConfig, logger util.ILogger, once Ionce) Ihook {
	return &hook{
		mutex:  sync.Mutex{},
		http:   http,
		config: config,
		logger: logger,
		once:   once,
	}
}

// hook ...
func (a *hook) RunHook() {
	// User information
	a.logger.PrintInfo("Starting in mode: HOOK")

	// Start our HTTP Server
	a.http.Start("/v1/run", a.httpHandler)
}

// run ...
func (a *hook) httpHandler(response http.ResponseWriter, request *http.Request) {
	// Make sure this is a GET request
	if request.Method != http.MethodGet {
		response.WriteHeader(http.StatusNotFound)
		_, _ = response.Write([]byte("400 - ups, page not found!"))
		return
	}

	// Defer our recover, so we can properly send an HTTP error
	// response and carry on serving subsequent requests
	defer func(logger util.ILogger) {
		if r := recover(); r != nil {
			var recoveredError = r.(util.GonsulError)
			response.WriteHeader(503)
			response.Header().Add("X-Gonsul-Error", string(util.ErrorDeleteNotAllowed))
			// Add delete paths(they wre added to logger messages) as comma separated string to the Header
			response.Header().Add("X-Gonsul-Delete-Paths", strings.Join(logger.GetMessages(), ","))
			_, _ = fmt.Fprintf(response, "Error: %d\n", recoveredError.Code)
		}
	}(a.logger)

	// Let's try to get a lock and defer the unlock
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.logger.PrintInfo("HTTP Incoming connection from: " + request.RemoteAddr)

	// On every request, run Once as usual business
	a.once.RunOnce()

	// If here, process ran smooth, return HTTP 200
	response.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(response, "Done")
}
