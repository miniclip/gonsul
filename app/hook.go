package app

import (
	"github.com/miniclip/gonsul/util"

	"github.com/gorilla/mux"

	"errors"
	"fmt"
	"net/http"
	"strings"
)

// hook ...
func (a *Application) hook() {
	// User information
	a.logger.PrintInfo("Starting in mode: HOOK")

	// Start our router and HTTP server
	router := mux.NewRouter()
	router.HandleFunc("/v1/run", a.hookHandler).Methods("GET")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		util.ExitError(errors.New("Hook: "+err.Error()), util.ErrorFailedHTTPServer, a.logger)
	}
}

// run ...
func (a *Application) hookHandler(response http.ResponseWriter, request *http.Request) {
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
	a.once()

	// If here, process ran smooth, return HTTP 200
	response.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(response, "Done")
}
