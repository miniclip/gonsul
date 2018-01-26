package importer

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/structs"

	"net/http"
	"errors"
	"time"
	"fmt"
	"encoding/json"
)

var config 		configuration.Config 		// Set our Configuration as global package scope
var logger 		errorutil.Logger     		// Set our Logger as global package scope

func Start(localData map[string]string, conf *configuration.Config, log *errorutil.Logger) {

	// Create some local variables
	var operations 	structs.OperationMatrix
	var liveData 	map[string]string

	// Set the appropriate values for our package global variables
	config 		= *conf
	logger 		= *log

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	// Populate our Consul live data
	liveData 	= createLiveData(client)

	// Create our operations Matrix
	operations 	= createOperationMatrix(liveData, localData)

	// Check if it's a dry run
	if conf.GetStrategy() == configuration.StrategyDry {
		// Print matrix and exit
		printOperations(operations, structs.OperationAll)

		return
	}

	// Process our operations matrix
	processOperations(operations)

	// Print result summary
	logger.PrintInfo(fmt.Sprintf("Inserts: %d records", operations.GetTotalInserts()))
	logger.PrintInfo(fmt.Sprintf("Updates: %d records", operations.GetTotalUpdates()))
	logger.PrintInfo(fmt.Sprintf("Deletes: %d records", operations.GetTotalDeletes()))
}

func processOperations(matrix structs.OperationMatrix) {
	// Did we got any deletes and are we allowed to delete them?
	if !config.AllowDeletes() && matrix.HasDeletes() {
		// We're not supposed to trigger Consul deletes, output report and exit with error
		logger.PrintError("We're stopping as there are deletes and Gonsul is running without delete permission")
		logger.PrintError("Below is all the Consul KV paths that should be deleted")

		// Print matrix and exit
		printOperations(matrix, structs.OperationDelete)
		errorutil.ExitError(errors.New(""), errorutil.ErrorDeleteNotAllowed, &logger)
	}

	var transactions []structs.ConsulTxn

	for _, op := range matrix.GetOperations()  {
		// We need to get the values to use pointers for our structure
		// so we can clearly identify nil values, as in https://willnorris.com/2014/05/go-rest-apis-and-pointers
		verb 			:= op.GetVerb()
		path 			:= op.GetPath()
		if op.GetType() == structs.OperationDelete {
			TxnKV 			:= structs.ConsulTxnKV{Verb: &verb, Key: &path}
			transactions 	= append(transactions, structs.ConsulTxn{KV: TxnKV})
		} else {
			val 			:= op.GetValue()
			TxnKV 			:= structs.ConsulTxnKV{Verb: &verb, Key: &path, Value: &val}
			transactions = append(transactions, structs.ConsulTxn{KV: TxnKV})
		}
	}

	json, _ := json.MarshalIndent(transactions, "", "  ")
	fmt.Println(string(json))
}

func processConsulTransaction(transactions []structs.ConsulTxn) {

}

func insertToConsul(path string, value string, client *http.Client) {
	/*
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
	*/
}