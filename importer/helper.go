package importer

import (
	"github.com/miniclip/gonsul/data"
	"encoding/json"
	"net/http"
	"github.com/miniclip/gonsul/errorutil"
	"errors"
	"io/ioutil"
	"encoding/base64"
)

const IsSkipping 	= 0
const IsInserting 	= 1
const IsUpdating 	= 2

func populateLiveData(client *http.Client) {
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
	liveData = map[string]string{}

	// Loop each entry on our Consul response
	for _, v := range bodyStruct {
		// Add to our map
		liveData[v.Key] = v.Value
	}
}

func shouldInsert(path string, value string) int {

	// Set our values (the original base64 response + convert actual value to base64)
	respValB64		:= liveData[path]
	currValB64		:= base64.StdEncoding.EncodeToString([]byte(value))

	// If values are equal return false so we do not write value
	if respValB64 == currValB64 {
		return IsSkipping
	}

	// Values are different, is it going to be an update or insert
	if respValB64 != "" {
		return IsUpdating
	}

	return IsInserting
}