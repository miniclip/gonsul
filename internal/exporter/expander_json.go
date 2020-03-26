package exporter

import (
	"github.com/miniclip/gonsul/internal/util"

	"encoding/json"
	"errors"
	"fmt"
)

// validateJSON ...
func (e *exporter) validateJSON(path string, jsonData string) map[string]interface{} {
	// Create "generic" json struct
	var arbitraryJSON map[string]interface{}

	// Decode data into "generic"
	err := json.Unmarshal([]byte(jsonData), &arbitraryJSON)

	// Decoded JSON ok?
	if err != nil {
		util.ExitError(
			errors.New(fmt.Sprintf("error parsing JSON file: %s with Message: %s", path, err.Error())),
			util.ErrorFailedJsonDecode,
			e.logger,
		)
	}

	return arbitraryJSON
}

// expandJSON ...
func (e *exporter) expandJSON(path string, jsonData string, localData map[string]string) {
	arbitraryJSON := e.validateJSON(path, jsonData)

	// Iterate over our "generic" JSON structure
	e.traverseMap(path, arbitraryJSON, localData)
}
