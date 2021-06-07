package importer

import (
	"strconv"

	"github.com/miniclip/gonsul/internal/entities"
	"github.com/miniclip/gonsul/internal/util"

	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const consulTxnLimit = 64
const maximumPayloadSize = 500000 // max size is actually 512kb

// processConsulTransaction ...
func (i *importer) processConsulTransaction(transactions []entities.ConsulTxn, batchNumber int) {
	batch := strconv.Itoa(batchNumber)

	// Encode our transaction into a JSON payload
	jsonPayload, err := json.Marshal(transactions)
	if err != nil {
		util.ExitError(errors.New("Marshal: "+err.Error()+" in Batch "+batch), util.ErrorFailedJsonEncode, i.logger)
	}

	// Create our URL
	consulUrl := i.config.GetConsulURL() + "/v1/txn"

	// build our request
	i.logger.PrintDebug("CONSUL: creating PUT request for Batch " + batch)
	req, err := http.NewRequest("PUT", consulUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		util.ExitError(errors.New("NewRequestPUT"+err.Error()+" in Batch "+batch), util.ErrorFailedConsulConnection, i.logger)
	}

	// Set ACL token
	if i.config.GetConsulACL() != "" {
		req.Header.Set("X-Consul-Token", i.config.GetConsulACL())
	}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	i.logger.PrintDebug("CONSUL: calling PUT request for Batch " + batch)
	resp, err := i.client.Do(req)
	if err != nil {
		util.ExitError(errors.New("DoPUT: "+err.Error()+" for Batch "+batch), util.ErrorFailedConsulConnection, i.logger)
	}

	// Clean response after function ends
	defer func() {
		if err := resp.Body.Close(); err != nil {
			i.logger.PrintError("Could not close Consul http body")
		}
	}()

	// Read the response body
	i.logger.PrintDebug("CONSUL: reading PUT response from Batch " + batch)
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.ExitError(errors.New("ReadPutResponse: "+err.Error()+" in Batch "+batch), util.ErrorFailedReadingResponse, i.logger)
	}

	// Cast response to string
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		util.ExitError(errors.New("TransactionError: "+bodyString+" in Batch "+batch), util.ErrorFailedConsulTxn, i.logger)
	}

	// All good. Output some status for each transaction operation
	for _, txn := range transactions {
		i.logger.PrintInfo("Operation: " + *txn.KV.Verb + " Path: " + *txn.KV.Key + " Batch: " + batch)
	}
}
