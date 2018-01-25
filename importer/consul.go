package importer

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/cbroglie/mustache"
	"net/http"
	"errors"
	"io/ioutil"
	"bytes"
	"time"
	"fmt"
)

var config 		configuration.Config 		// Set our Configuration as global package scope
var logger 		errorutil.Logger     		// Set our Logger as global package scope
var localData	map[string]string       	// Our map that will hold our processed local data
var liveData	map[string]string       	// Our map that will hold our live Consul data
var insertCount	int							// A simple insert counter
var updateCount	int							// A simple update counter
var deleteCount	int							// A simple delete counter

func Start(data map[string]string, conf *configuration.Config, log *errorutil.Logger) {

	// Set the appropriate values for our package global variables
	localData 	= data
	config 		= *conf
	logger 		= *log
	insertCount, updateCount, deleteCount = 0, 0, 0

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Populate our Consul live data, to compare before writes
	populateLiveData(client)

	// Create our Deletion checking
	checkDeletes(client)

	// Iterate over our import data
	for k, v := range localData {
		insertToConsul(k, v, client)
	}

	logger.PrintInfo(fmt.Sprintf("Inserts: %d records", insertCount))
	logger.PrintInfo(fmt.Sprintf("Updates: %d records", updateCount))
	logger.PrintInfo(fmt.Sprintf("Deletes: %d records", deleteCount))
}

func checkDeletes(client *http.Client) {
	var deletes []string

	// Check for deletes
	for k := range liveData {
		if _, ok := localData[k]; !ok {
			// Not found in local - DELETE
			deletes = append(deletes, k)
		}
	}

	// Did we got any deletes and are we allowed to delete them?
	if !config.AllowDeletes() && len(deletes) != 0 {
		// We're not supposed to trigger Consul deletes, output report and exit with error
		logger.PrintError("We're stopping as there are deletes and Gonsul is running without delete permission")
		logger.PrintError("Below is all the Consul KV paths that should be deleted")
		for _, keyForDeletion := range deletes {
			logger.PrintError("- " + keyForDeletion)
		}
		errorutil.ExitError(errors.New(""), errorutil.ErrorDeleteNotAllowed, &logger)
	} else if len(deletes) != 0 {
		// We found some deletes to do, and we're allowed to do it. Loop each one triggering the DELETE request
		for _, keyForDeletion := range deletes {
			deleteFromConsul(keyForDeletion, client)
		}
	}
}

func insertToConsul(path string, value string, client *http.Client) {
	var err error

	// Create our URL
	consulUrl := config.GetConsulURL() + "/v1/kv/" + path
	logger.PrintDebug("CONSUL: Importing to URL: " + consulUrl)

	// Shall we run secret replacement
	if config.DoSecrets() {
		value, err = mustache.Render(value, config.GetSecretsMap())
	}

	insertOrUpdate := shouldInsert(path, value)

	// Check if we should insert the value, to save writes on Consul cluster
	if insertOrUpdate == IsSkipping {
		//logger.PrintInfo("IMPORTING - " + path + " -> Skip")
		logger.PrintDebug("CONSUL: skipping as consul and repo data are equal")
		return
	}

	// build our request
	logger.PrintDebug("CONSUL: creating PUT request")
	req, err := http.NewRequest("PUT", consulUrl, bytes.NewBufferString(value))
	if err != nil {
		errorutil.ExitError(errors.New("NewRequestPUT: " + err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Set ACL token
	req.Header.Set("X-Consul-Token", config.GetConsulACL())

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	logger.PrintDebug("CONSUL: calling PUT request")
	resp, err := client.Do(req)
	if err != nil {
		errorutil.ExitError(errors.New("DoPUT: " + err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Clean response after function ends
	defer resp.Body.Close()

	// Read the response body
	logger.PrintDebug("CONSUL: reading PUT response")
	bodyBytes, err 	:= ioutil.ReadAll(resp.Body)
	if err != nil {
		errorutil.ExitError(errors.New("ReadPutResponse: " + err.Error()), errorutil.ErrorFailedReadingResponse, &logger)
	}
	// Cast response to string
	bodyString 		:= string(bodyBytes)

	if insertOrUpdate == IsUpdating {
		logger.PrintInfo("UPDATING  - " + path + " -> " + bodyString)
		updateCount++
	} else {
		logger.PrintInfo("INSERTING - " + path + " -> " + bodyString)
		insertCount++
	}

}

func deleteFromConsul(path string, client *http.Client) {
	var err error

	// Create our URL
	consulUrl := config.GetConsulURL() + "/v1/kv/" + path
	logger.PrintDebug("CONSUL: Deleting URL: " + consulUrl)

	// build our request
	logger.PrintDebug("CONSUL: creating DELETE request")
	req, err := http.NewRequest("DELETE", consulUrl, nil)
	if err != nil {
		errorutil.ExitError(errors.New("NewRequestDELETE: " + err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Set ACL token
	req.Header.Set("X-Consul-Token", config.GetConsulACL())

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	logger.PrintDebug("CONSUL: calling DELETE request")
	resp, err := client.Do(req)
	if err != nil {
		errorutil.ExitError(errors.New("DoDELETE: " + err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Clean response after function ends
	defer resp.Body.Close()

	// Read the response body
	logger.PrintDebug("CONSUL: reading DELETE response")
	bodyBytes, err 	:= ioutil.ReadAll(resp.Body)
	if err != nil {
		errorutil.ExitError(errors.New("ReadDeleteResponse: " + err.Error()), errorutil.ErrorFailedReadingResponse, &logger)
	}
	// Cast response to string
	bodyString 		:= string(bodyBytes)

	logger.PrintInfo("DELETING  - " + path + " -> " + bodyString)
	deleteCount++
}