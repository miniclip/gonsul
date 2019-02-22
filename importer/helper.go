package importer

import (
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/structs"

	"github.com/cbroglie/mustache"
	"github.com/olekukonko/tablewriter"

	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func createOperationMatrix(liveData map[string]string, localData map[string]string) structs.OperationMatrix {
	// Set local error variable
	var err error

	// Create our Operations array
	var operations = structs.NewOperationsMatrix()

	// Check for updates or inserts
	for localKey, localVal := range localData {

		// Make sure we do not have an empty value (Consul KV will not have it)
		if localVal == "" {
			continue
		}

		// Shall we run secret replacement
		if config.DoSecrets() {
			localVal, err = mustache.Render(localVal, config.GetSecretsMap())
		}
		if err != nil {
			errorutil.ExitError(errors.New("MustacheRender: "+err.Error()), errorutil.ErrorFailedMustache, &logger)
		}

		// Base64 encode local value
		localValB64 := base64.StdEncoding.EncodeToString([]byte(localVal))

		// Does the current local KV key (path) exists in live?
		if liveVal, ok := liveData[localKey]; ok {
			// it does, is it different value?
			if localValB64 != liveVal {
				// Gentleman we have an update
				operations.AddUpdate(structs.Entry{KVPath: localKey, Value: localValB64})
			}
		} else {
			// Current key does not exist in live data, we have an insert
			operations.AddInsert(structs.Entry{KVPath: localKey, Value: localValB64})
		}
	}
	// Now check for deletes
	// Check for deletes
	for liveKey := range liveData {
		if _, ok := localData[liveKey]; !ok {
			// Not found in local - DELETE
			operations.AddDelete(structs.Entry{KVPath: liveKey, Value: ""})
		}
	}

	return operations
}

func createLiveData(client *http.Client) map[string]string {
	// Create some local variables
	var liveData map[string]string

	// Create our URL
	consulUrl := config.GetConsulURL() + "/v1/kv/" + config.GetConsulbasePath() + "/?recurse=true"
	// build our request
	req, err := http.NewRequest("GET", consulUrl, nil)
	if err != nil {
		errorutil.ExitError(errors.New("NewRequestGET: "+err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Set ACL token (if given)
	if config.GetConsulACL() != "" {
		req.Header.Set("X-Consul-Token", config.GetConsulACL())
	}

	// Send the request via a client, Do sends an HTTP request and returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		errorutil.ExitError(errors.New("DoGET: "+err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Clean response after function ends
	defer resp.Body.Close()

	// Invalid response, path is empty then, fresh import
	if resp.StatusCode == 404 {
		return nil
	}

	if resp.StatusCode >= 400 {
		errorutil.ExitError(errors.New("Invalid response from consul: "+resp.Status), errorutil.ErrorFailedConsulConnection, &logger)
  }

	// Read response from HTTP Response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorutil.ExitError(errors.New("ReadGetResponse: "+err.Error()), errorutil.ErrorFailedReadingResponse, &logger)
	}
	// Create a structure for our response, basically an array of
	// Consul results because we're doing a recurse call
	var bodyStruct []structs.ConsulResult
	// Convert response to a string and then parse it to our struct
	bodyString := string(bodyBytes)
	err = json.Unmarshal([]byte(bodyString), &bodyStruct)
	if err != nil {
		errorutil.ExitError(errors.New("Unmarshal: "+err.Error()), errorutil.ErrorFailedJsonDecode, &logger)
	}

	// All good so far, instantiate our map
	liveData = map[string]string{}

	// Loop each entry on our Consul response
	for _, v := range bodyStruct {
		// Add to our map
		liveData[v.Key] = v.Value
	}

	return liveData
}

func printOperations(matrix structs.OperationMatrix, printWhat string) {
	// Add a new line before the table
	fmt.Println()
	// Let's make sure there are any operation
	if matrix.GetTotalOps() > 0 {
		// Instantiate our table and set table header
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"", "OPERATION NAME", "CONSUL VERB", "PATH"})
		// Align our rows
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		// Loop each operation and add to table
		for _, op := range matrix.GetOperations() {
			if printWhat == structs.OperationAll || printWhat == op.GetType() {
				if op.GetType() == structs.OperationDelete {
					table.Append([]string{"!!", op.GetType(), op.GetVerb(), op.GetPath()})
				} else {
					table.Append([]string{"", op.GetType(), op.GetVerb(), op.GetPath()})
				}
			}
		}
		// Outputs ASCII table
		table.Render()
	} else {
		logger.PrintInfo("No operations to process, all synced")
	}
}

func setDeletesToLogger(matrix structs.OperationMatrix) {
	// Let's make sure there are any operation
	if matrix.GetTotalOps() > 0 {
		// Loop each operation and add to table
		for _, op := range matrix.GetOperations() {
			if op.GetType() == structs.OperationDelete {
				logger.AddMessage(op.GetPath())
			}
		}
	}
}
