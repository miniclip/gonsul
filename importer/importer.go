package importer

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/structs"

	"errors"
	"fmt"
	"net/http"
	"time"
)

var config configuration.Config // Set our Configuration as global package scope
var logger errorutil.Logger     // Set our Logger as global package scope

func Start(localData map[string]string, conf *configuration.Config, log *errorutil.Logger) {

	// Create some local variables
	var ops structs.OperationMatrix
	var liveData map[string]string

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

	// Populate our Consul live data
	liveData = createLiveData(client)

	// Create our operations Matrix
	ops = createOperationMatrix(liveData, localData)

	// Check if it's a dry run
	if conf.GetStrategy() == configuration.StrategyDry {
		// Print matrix and exit
		printOperations(ops, structs.OperationAll)

		return
	}

	// Process our operations matrix
	processOperations(ops, client)

	// Print result summary
	logger.PrintInfo(fmt.Sprintf("Finished: %d Inserts, %d Updates %d Deletes", ops.GetTotalInserts(), ops.GetTotalUpdates(), ops.GetTotalDeletes()))
}

func processOperations(matrix structs.OperationMatrix, client *http.Client) {
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

	for _, op := range matrix.GetOperations() {
		// We need to get the values to use pointers for our structure
		// so we can clearly identify nil values, as in https://willnorris.com/2014/05/go-rest-apis-and-pointers
		verb := op.GetVerb()
		path := op.GetPath()
		if op.GetType() == structs.OperationDelete {
			TxnKV := structs.ConsulTxnKV{Verb: &verb, Key: &path}
			transactions = append(transactions, structs.ConsulTxn{KV: TxnKV})
		} else {
			val := op.GetValue()
			TxnKV := structs.ConsulTxnKV{Verb: &verb, Key: &path, Value: &val}
			transactions = append(transactions, structs.ConsulTxn{KV: TxnKV})
		}

		if len(transactions) == ConsulTxnLimit {
			// Flush current transactions because we hit max operation per transaction
			// One day Consul will release an API endpoint from where we can get this limit
			// so we do can stop hardcoding this constant
			processConsulTransaction(transactions, client)
			// Reset our transaction array
			transactions = []structs.ConsulTxn{}
		}
	}

	// Do we have transactions to process
	if len(transactions) > 0 {
		processConsulTransaction(transactions, client)
	}
}
