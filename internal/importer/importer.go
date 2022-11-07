package importer

import (
	"github.com/miniclip/gonsul/internal/config"
	"github.com/miniclip/gonsul/internal/entities"
	"github.com/miniclip/gonsul/internal/util"

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
	config config.IConfig
	logger util.ILogger
	client *http.Client
}

// NewImporter
func NewImporter(config config.IConfig, logger util.ILogger, client *http.Client) IImporter {
	return &importer{config: config, logger: logger, client: client}
}

// Start ...
func (i *importer) Start(localData map[string]string) {

	// Create some local variables
	var ops entities.OperationMatrix

	// Populate our Consul live data
	liveData := i.createLiveData()

	if i.config.GetStrategy() == config.StrategyRead && i.config.GetOutputDir() == "" && i.config.GetOutputFile() == "" {
		util.ExitError(errors.New("undefined output-dir or output-file"), util.ErrorWrite, i.logger)
	}

	if i.config.GetOutputDir() != "" {
		if err := i.exportToDirectory(i.config.GetOutputDir(), liveData); err != nil {
			util.ExitError(errors.New(err.Error()), util.ErrorWrite, i.logger)
		}
	}
	if i.config.GetOutputFile() != "" {
		if err := i.exportToFile(i.config.GetOutputFile(), liveData, false); err != nil {
			util.ExitError(errors.New(err.Error()), util.ErrorWrite, i.logger)
		}
	}
	// Check if it's read-only
	if i.config.GetStrategy() == config.StrategyRead {
		return
	}

	// Create our operations Matrix
	ops = i.createOperationMatrix(liveData, localData)

	// Print operation table
	i.printOperations(ops, entities.OperationAll, i.config.GetPrintValues())
	// Check if it's a dry run
	if i.config.GetStrategy() == config.StrategyDry {
		// Exit after having printed the operations table
		return
	}

	// Process our operations matrix
	i.processOperations(ops)

	// Print result summary
	i.logger.PrintInfo(fmt.Sprintf("Finished: %d Inserts, %d Updates %d Deletes", ops.GetTotalInserts(), ops.GetTotalUpdates(), ops.GetTotalDeletes()))
}

func (i *importer) processOperations(matrix entities.OperationMatrix) {
	// Did we got any deletes and are we allowed to delete them?
	if i.config.AllowDeletes() == "false" && matrix.HasDeletes() {
		// We're not supposed to trigger Consul deletes, output report and exit with error
		i.logger.PrintError("We're stopping as there are deletes and Gonsul is running without delete permission")
		i.logger.PrintError("Below is all the Consul KV paths that would be deleted")

		// Print matrix (or set in logger messages if in hook mode) and exit
		if i.config.GetStrategy() == config.StrategyHook {
			i.setDeletesToLogger(matrix)
		} else {
			i.printOperations(matrix, entities.OperationDelete, i.config.GetPrintValues())
		}
		util.ExitError(errors.New(""), util.ErrorDeleteNotAllowed, i.logger)
	}

	// Initialize the batch counter
	batch := 1

	var transactions []entities.ConsulTxn
	var newTransactions []entities.ConsulTxn

	// Fill our channel to indicate a non interruptible work (It stops here if interruption in progress)
	i.config.WorkingChan() <- true

	// Loop each operation
	for _, op := range matrix.GetOperations() {
		// We need to get the values to use pointers for our structure
		// so we can clearly identify nil values, as in https://willnorris.com/2014/05/go-rest-apis-and-pointers
		verb := op.GetVerb()
		path := op.GetPath()

		var TxnKV entities.ConsulTxnKV

		if op.GetType() == entities.OperationDelete {
			TxnKV = entities.ConsulTxnKV{Verb: &verb, Key: &path}
		} else {
			val := op.GetValue()
			TxnKV = entities.ConsulTxnKV{Verb: &verb, Key: &path, Value: &val}
		}

		// add the next transaction and check payload lenght
		newTransactions = transactions
		newTransactions = append(transactions, entities.ConsulTxn{KV: TxnKV})
		newPayloadSize := i.getTransactionsPayloadSize(&newTransactions)

		if newPayloadSize > maximumPayloadSize || len(transactions) == consulTxnLimit {
			i.processConsulTransaction(transactions, batch)
			transactions = []entities.ConsulTxn{}

			batch++
		}

		transactions = append(transactions, entities.ConsulTxn{KV: TxnKV})
	}

	// Do we have transactions to process
	if len(transactions) > 0 {
		i.processConsulTransaction(transactions, batch)
	}

	// Consume our channel, to re-allow application interruption
	<-i.config.WorkingChan()
}
