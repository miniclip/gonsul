package importer

import (
	"github.com/miniclip/gonsul/configuration"
	"github.com/miniclip/gonsul/structs"
	"github.com/miniclip/gonsul/util"

	"errors"
	"fmt"
	"net/http"
)

// IImporter ...
type IImporter interface {
	Start(localData map[string]string)
}

// importer ...
type importer struct {
	config configuration.IConfig
	logger util.ILogger
	client *http.Client
}

// NewImporter
func NewImporter(config configuration.IConfig, logger util.ILogger, client *http.Client) IImporter {
	return &importer{config: config, logger: logger, client: client}
}

// Start ...
func (i *importer) Start(localData map[string]string) {

	// Create some local variables
	var ops structs.OperationMatrix
	var liveData map[string]string

	// Populate our Consul live data
	liveData = i.createLiveData()

	// Create our operations Matrix
	ops = i.createOperationMatrix(liveData, localData)

	// Check if it's a dry run
	if i.config.GetStrategy() == configuration.StrategyDry {
		// Print matrix and exit
		i.printOperations(ops, structs.OperationAll)

		return
	}

	// Process our operations matrix
	i.processOperations(ops)

	// Print result summary
	i.logger.PrintInfo(fmt.Sprintf("Finished: %d Inserts, %d Updates %d Deletes", ops.GetTotalInserts(), ops.GetTotalUpdates(), ops.GetTotalDeletes()))
}

func (i *importer) processOperations(matrix structs.OperationMatrix) {
	// Did we got any deletes and are we allowed to delete them?
	if !i.config.AllowDeletes() && matrix.HasDeletes() {
		// We're not supposed to trigger Consul deletes, output report and exit with error
		i.logger.PrintError("We're stopping as there are deletes and Gonsul is running without delete permission")
		i.logger.PrintError("Below is all the Consul KV paths that should be deleted")

		// Print matrix (or set in logger messages if in hook mode) and exit
		if i.config.GetStrategy() == configuration.StrategyHook {
			i.setDeletesToLogger(matrix)
		} else {
			i.printOperations(matrix, structs.OperationDelete)
		}
		util.ExitError(errors.New(""), util.ErrorDeleteNotAllowed, i.logger)
	}

	var transactions []structs.ConsulTxn

	// Fill our channel to indicate a non interruptible work (It stops here if interruption in progress)
	i.config.WorkingChan() <- true

	// Loop each operation
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

		if len(transactions) == consulTxnLimit {
			// Flush current transactions because we hit max operation per transaction
			// One day Consul will release an API endpoint from where we can get this limit
			// so we do can stop hardcoding this constant
			i.processConsulTransaction(transactions)
			// Reset our transaction array
			transactions = []structs.ConsulTxn{}
		}
	}

	// Do we have transactions to process
	if len(transactions) > 0 {
		i.processConsulTransaction(transactions)
	}

	// Consume our channel, to re-allow application interruption
	<-i.config.WorkingChan()
}
