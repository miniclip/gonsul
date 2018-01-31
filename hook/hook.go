package hook

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/once"

	"github.com/gorilla/mux"

	"errors"
	"fmt"
	"net/http"
	"sync"
	"gopkg.in/src-d/go-git.v4/utils/merkletrie/noder"
	"strings"
)

var mutex *sync.Mutex
var config configuration.Config // Set our Configuration as global package scope
var logger errorutil.Logger     // Set our Logger as global package scope

func Start(conf *configuration.Config, log *errorutil.Logger) {
	// Set the appropriate values for our package global variables
	config = *conf
	logger = *log

	// Initialize our mutex
	mutex = &sync.Mutex{}
	// Start our router and HTTP server
	router := mux.NewRouter()
	router.HandleFunc("/v1/run", hookHandler).Methods("GET")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		errorutil.ExitError(errors.New("Hook: "+err.Error()), errorutil.ErrorFailedHTTPServer, &logger)
	}
}

func hookHandler(response http.ResponseWriter, request *http.Request) {
	// Defer our recover, so we can properly send an HTTP error
	// response and carry on serving subsequent requests
	defer func(logger errorutil.Logger) {
		if r := recover(); r != nil {
			var recoveredError = r.(errorutil.GonsulError)
			response.WriteHeader(503)
			response.Header().Add("X-Gonsul-Error", string(errorutil.ErrorDeleteNotAllowed))
			// Add delete paths(they wre added to logger messages) as comma separated string to the Header
			response.Header().Add("X-Gonsul-Delete-Paths", strings.Join(logger.GetMessages(), ","))
			fmt.Fprintf(response, "Error: %d\n", recoveredError.Code)
		}
	}(logger)

	// Let's try to get a lock and defer the unlock
	mutex.Lock()
	defer mutex.Unlock()

	logger.PrintInfo("HTTP Incoming connection from: " + request.RemoteAddr)

	// On every request, run once as usual business
	once.Start(&config, &logger)

	// If here, process ran smooth, return HTTP 200
	response.WriteHeader(http.StatusOK)
	fmt.Fprint(response, "Done")
}
