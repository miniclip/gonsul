package importer

import (
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/structs"

	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const ConsulTxnLimit = 64

func processConsulTransaction(transactions []structs.ConsulTxn, client *http.Client) {
	// Encode our transaction into a JSON payload
	jsonPayload, err := json.Marshal(transactions)
	if err != nil {
		errorutil.ExitError(errors.New("Marshal: "+err.Error()), errorutil.ErrorFailedJsonEncode, &logger)
	}

	// Create our URL
	consulUrl := config.GetConsulURL() + "/v1/txn"

	// build our request
	logger.PrintDebug("CONSUL: creating PUT request")
	req, err := http.NewRequest("PUT", consulUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		errorutil.ExitError(errors.New("NewRequestPUT: "+err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Set ACL token
	req.Header.Set("X-Consul-Token", config.GetConsulACL())

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	logger.PrintDebug("CONSUL: calling PUT request")
	resp, err := client.Do(req)
	if err != nil {
		errorutil.ExitError(errors.New("DoPUT: "+err.Error()), errorutil.ErrorFailedConsulConnection, &logger)
	}

	// Clean response after function ends
	defer resp.Body.Close()

	// Read the response body
	logger.PrintDebug("CONSUL: reading PUT response")
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errorutil.ExitError(errors.New("ReadPutResponse: "+err.Error()), errorutil.ErrorFailedReadingResponse, &logger)
	}

	// Cast response to string
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		errorutil.ExitError(errors.New("TransactionError: "+bodyString), errorutil.ErrorFailedConsulTxn, &logger)
	}

	// All good. Output some status for each transaction operation
	for _, txn := range transactions {
		logger.PrintInfo("Operation: " + *txn.KV.Verb + " Path: " + *txn.KV.Key)
	}
}
