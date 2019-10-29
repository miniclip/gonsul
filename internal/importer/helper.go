package importer

import (
	"github.com/miniclip/gonsul/internal/entities"
	"github.com/miniclip/gonsul/internal/util"
	"strings"

	"github.com/cbroglie/mustache"
	"github.com/olekukonko/tablewriter"

	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

// createOperationMatrix ...
func (i *importer) createOperationMatrix(liveData map[string]string, localData map[string]string) entities.OperationMatrix {
	// Set local error variable
	var err error
	// Create our Operations array
	var operations = entities.NewOperationsMatrix()

	// Check for updates or inserts
	for localKey, localVal := range localData {
		// Make sure we do not have an empty value (Consul KV will not have it)
		if localVal == "" {
			continue
		}

		// Shall we run secret replacement
		if i.config.DoSecrets() {
			localVal, err = mustache.Render(localVal, i.config.GetSecretsMap())
		}
		if err != nil {
			util.ExitError(errors.New("MustacheRender: "+err.Error()), util.ErrorFailedMustache, i.logger)
		}

		// Base64 encode local value
		localValB64 := base64.StdEncoding.EncodeToString([]byte(localVal))

		// Does the current local KV key (path) exists in live?
		if liveVal, ok := liveData[localKey]; ok {
			// it does, is it different value?
			if localValB64 != liveVal {
				// Gentleman we have an update
				operations.AddUpdate(entities.Entry{KVPath: localKey, Value: localValB64})
			}
		} else {
			// Current key does not exist in live data, we have an insert
			operations.AddInsert(entities.Entry{KVPath: localKey, Value: localValB64})
		}
	}

	// Now check for deletes
	// Check for deletes
	for liveKey := range liveData {
		if _, ok := localData[liveKey]; !ok && i.config.AllowDeletes() != "skip" {
			// Not found in local - DELETE
			operations.AddDelete(entities.Entry{KVPath: liveKey, Value: ""})
		}
	}

	return operations
}

// createLiveData ...
func (i *importer) createLiveData() map[string]string {
	// Create some local variables
	var liveData map[string]string

	// Create our URL
	consulBasePath := strings.TrimSuffix(i.config.GetConsulBasePath(), "/")
	fullUrl := path.Join("v1", "kv", consulBasePath)
	hostname := strings.TrimSuffix(i.config.GetConsulURL(), "/")
	consulUrl := hostname + "/" + fullUrl + "/?recurse=true"
	// build our request
	req, err := http.NewRequest("GET", consulUrl, nil)
	if err != nil {
		util.ExitError(errors.New("NewRequestGET: "+err.Error()), util.ErrorFailedConsulConnection, i.logger)
	}

	// Set ACL token (if given)
	if i.config.GetConsulACL() != "" {
		req.Header.Set("X-Consul-Token", i.config.GetConsulACL())
	}

	// Send the request via a client, Do sends an HTTP request and returns an HTTP response
	resp, err := i.client.Do(req)
	if err != nil {
		util.ExitError(errors.New("DoGET: "+err.Error()), util.ErrorFailedConsulConnection, i.logger)
	}

	// Clean response after function ends
	defer func() {
		if err := resp.Body.Close(); err != nil {
			i.logger.PrintError("Could not close Consul http body")
		}
	}()

	// Invalid response, path is empty then, fresh import
	if resp.StatusCode == 404 {
		return nil
	}

	if resp.StatusCode >= 400 {
		util.ExitError(errors.New("Invalid response from consul: "+resp.Status), util.ErrorFailedConsulConnection, i.logger)
	}

	// Read response from HTTP Response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.ExitError(errors.New("ReadGetResponse: "+err.Error()), util.ErrorFailedReadingResponse, i.logger)
	}
	// Create a structure for our response, basically an array of
	// Consul results because we're doing a recurse call
	var bodyStruct []entities.ConsulResult
	// Convert response to a string and then parse it to our struct
	bodyString := string(bodyBytes)
	err = json.Unmarshal([]byte(bodyString), &bodyStruct)
	if err != nil {
		util.ExitError(errors.New("Unmarshal: "+err.Error()), util.ErrorFailedJsonDecode, i.logger)
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

// printOperations ...
func (i *importer) printOperations(matrix entities.OperationMatrix, printWhat string) {
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
			if printWhat == entities.OperationAll || printWhat == op.GetType() {
				if op.GetType() == entities.OperationDelete {
					table.Append([]string{"!!", op.GetType(), op.GetVerb(), op.GetPath()})
				} else {
					table.Append([]string{"", op.GetType(), op.GetVerb(), op.GetPath()})
				}
			}
		}
		// Outputs ASCII table
		table.Render()
	} else {
		i.logger.PrintInfo("No operations to process, all synced")
	}
}

// setDeletesToLogger ...
func (i *importer) setDeletesToLogger(matrix entities.OperationMatrix) {
	// Let's make sure there are any operation
	if matrix.GetTotalOps() > 0 {
		// Loop each operation and add to table
		for _, op := range matrix.GetOperations() {
			if op.GetType() == entities.OperationDelete {
				i.logger.AddMessage(op.GetPath())
			}
		}
	}
}
