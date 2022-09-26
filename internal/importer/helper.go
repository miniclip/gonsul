package importer

import (
	"strings"

	"github.com/miniclip/gonsul/internal/entities"
	"github.com/miniclip/gonsul/internal/util"

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
	"strconv"
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
func (i *importer) printOperations(matrix entities.OperationMatrix, printWhat string, printValue bool) {
	// Add a new line before the table
	fmt.Println()
	// Let's make sure there are any operation
	if matrix.GetTotalOps() > 0 {
		// Instantiate our table and set table header
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"", "BATCH", "OP INDEX", "OPERATION NAME", "CONSUL VERB", "PATH", "VALUE"})
		// Align our rows
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		// Initialize the batch counter
		batch := 1
		opIndex := 0

		var transactions []entities.ConsulTxn
		var newTransactions []entities.ConsulTxn

		// Loop each operation and add to table
		for _, op := range matrix.GetOperations() {

			if printWhat == entities.OperationAll || printWhat == op.GetType() {
				var TxnKV entities.ConsulTxnKV
				var warning string

				// generate the actual payload to calculate it's lenght
				verb := op.GetVerb()
				path := op.GetPath()
				if op.GetType() == entities.OperationDelete {
					warning = "!!"
					TxnKV = entities.ConsulTxnKV{Verb: &verb, Key: &path}
				} else {
					warning = ""
					val := op.GetValue()
					TxnKV = entities.ConsulTxnKV{Verb: &verb, Key: &path, Value: &val}
				}

				// add the next transaction and check payload lenght
				newTransactions = transactions
				newTransactions = append(newTransactions, entities.ConsulTxn{KV: TxnKV})
				newPayloadSize := i.getTransactionsPayloadSize(&newTransactions)

				// If the next transaction brings us over the maximum payload size,
				// or the maximum transaction per batch limit is reached, start a new batch
				if newPayloadSize > maximumPayloadSize || len(transactions) == consulTxnLimit {
					// reset transactions and add the next transaction
					transactions = []entities.ConsulTxn{}
					// start a new batch counter
					opIndex = 0
					batch++
				}

				transactions = append(transactions, entities.ConsulTxn{KV: TxnKV})
				opValue := ""
				if printValue {
					opValue = i.decodeOpValue(op.GetValue())
				}
				table.Append([]string{warning, strconv.Itoa(batch), strconv.Itoa(opIndex), op.GetType(), op.GetVerb(), op.GetPath(), opValue})

				opIndex++
			}
		}
		// Outputs ASCII table
		table.Render()
	} else {
		i.logger.PrintInfo("No operations to process, all synced")
	}
}

func (i *importer) decodeOpValue(opValue string) string {
	if opValue == "" {
		return opValue
	}
	decodedValue, err := base64.StdEncoding.DecodeString(opValue)
	if err != nil {
		util.ExitError(errors.New(err.Error()), util.ErrorRead, i.logger)
	}
	return string(decodedValue)
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

// Get the payload size for a slice of transactions
func (i *importer) getTransactionsPayloadSize(transactions *[]entities.ConsulTxn) int {
	payload, err := json.Marshal(&transactions)
	if err != nil {
		util.ExitError(errors.New("Marshal: "+err.Error()), util.ErrorFailedJsonEncode, i.logger)
	}

	return len(string(payload))
}
