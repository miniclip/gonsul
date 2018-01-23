package importer

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/data"
	"github.com/cbroglie/mustache"
	"net/http"
	"errors"
	"io/ioutil"
	"bytes"
	"time"
	"encoding/json"
	"encoding/base64"
)

var config 			configuration.Config		// Set our Configuration as global package scope
var logger 			errorutil.Logger			// Set our Logger as global package scope
var consulLiveData	map[string]string			// Create a map with Consul recursive GET, so it's easier to search

func Start(data data.EntryCollection, conf *configuration.Config, log *errorutil.Logger) {

	// Set the appropriate values for our package global variables
	config = *conf
	logger = *log

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Populate our Consul live data, to compare before writes
	populateConsulLiveData(client)

	// Iterate over our import data
	for _, entry := range data.Entries {
		insert(entry.KVPath, entry.Value, client)
	}
}

func populateConsulLiveData(client *http.Client) {
	// Create our URL
	consulUrl := config.GetConsulURL() + "/v1/kv/" + config.GetConsulbasePath() + "?recurse=true"
	// build our request
	req, err := http.NewRequest("GET", consulUrl, nil)
	if err != nil {
		errorutil.ExitError(errors.New("NewRequestGET: " + err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Set ACL token
	req.Header.Set("X-Consul-Token", config.GetConsulACL())

	// Send the request via a client, Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		errorutil.ExitError(errors.New("DoGET: " + err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Clean response after function ends
	defer resp.Body.Close()

	// Invalid response, path is empty then, fresh import
	if resp.StatusCode == 404 {
		return
	}

	// Read response from HTTP Response
	bodyBytes, err 	:= ioutil.ReadAll(resp.Body)
	if err != nil {
		errorutil.ExitError(errors.New("ReadGetResponse: " + err.Error()), errorutil.ErrorFailedReadingResponse, &logger)
	}
	// Create a structure for our response, basically an array of
	// Consul results because we're doing a recurse call
	var bodyStruct	[]data.ConsulResult
	// Convert response to a string and then parse it to our struct
	bodyString := string(bodyBytes)
	err = json.Unmarshal([]byte(bodyString), &bodyStruct)
	if err != nil {
		errorutil.ExitError(errors.New("Unmarshal: " + err.Error()), errorutil.ErrorFailedJsonDecode, &logger)
	}

	// All good so far, instantiate our map
	consulLiveData = map[string]string{}

	// Loop each entry on our Consul response
	for _, v := range bodyStruct {
		// Add to our map
		consulLiveData[v.Key] = v.Value
	}
}

func insert(path string, value string, client *http.Client) {
	var err error

	// Create our URL
	consulUrl := config.GetConsulURL() + "/v1/kv/" + path
	logger.PrintDebug("CONSUL: Importing to URL: " + consulUrl)

	// Shall we run secret replacement
	if config.DoSecrets() {
		value, err = mustache.Render(value, config.GetSecretsMap())
	}

	// Check if we should insert the value, to save writes on Consul cluster
	if !shouldInsert(path, value) {
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

	logger.PrintInfo("IMPORTING - " + path + " -> " + bodyString)
}


func shouldInsert(path string, value string) bool {

	// Set our values (the original base64 response + convert actual value to base64)
	respValB64		:= consulLiveData[path]
	currValB64		:= base64.StdEncoding.EncodeToString([]byte(value))

	// If values are equal return false so we do not write value
	if respValB64 == currValB64 {
		return false
	}

	// Values are different, we should let caller know that value must be written to Consul
	return true
}