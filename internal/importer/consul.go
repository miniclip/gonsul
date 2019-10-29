package importer

import (
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
func (i *importer) processConsulTransaction(transactions []entities.ConsulTxn) {
	// Encode our transaction into a JSON payload
	jsonPayload, err := json.Marshal(transactions)
	if err != nil {
		util.ExitError(errors.New("Marshal: "+err.Error()), util.ErrorFailedJsonEncode, i.logger)
	}

	// Create our URL
	consulUrl := i.config.GetConsulURL() + "/v1/txn"

	// build our request
	i.logger.PrintDebug("CONSUL: creating PUT request")
	req, err := http.NewRequest("PUT", consulUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		util.ExitError(errors.New("NewRequestPUT: "+err.Error()), util.ErrorFailedConsulConnection, i.logger)
	}

	// Set ACL token
	if i.config.GetConsulACL() != "" {
		req.Header.Set("X-Consul-Token", i.config.GetConsulACL())
	}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	i.logger.PrintDebug("CONSUL: calling PUT request")
	resp, err := i.client.Do(req)
	if err != nil {
		util.ExitError(errors.New("DoPUT: "+err.Error()), util.ErrorFailedConsulConnection, i.logger)
	}

	// Clean response after function ends
	defer func() {
		if err := resp.Body.Close(); err != nil {
			i.logger.PrintError("Could not close Consul http body")
		}
	}()

	// Read the response body
	i.logger.PrintDebug("CONSUL: reading PUT response")
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.ExitError(errors.New("ReadPutResponse: "+err.Error()), util.ErrorFailedReadingResponse, i.logger)
	}

	// Cast response to string
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		util.ExitError(errors.New("TransactionError: "+bodyString), util.ErrorFailedConsulTxn, i.logger)
	}

	// All good. Output some status for each transaction operation
	for _, txn := range transactions {
		i.logger.PrintInfo("Operation: " + *txn.KV.Verb + " Path: " + *txn.KV.Key)
	}
}
