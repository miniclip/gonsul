package exporter

import (
	"github.com/miniclip/gonsul/internal/util"

	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	e.traverseJSON(path, arbitraryJSON, localData)
}

// traverseJSON ...
func (e *exporter) traverseJSON(path string, arbitraryJSON map[string]interface{}, localData map[string]string) {
	for key, value := range arbitraryJSON {
		// Append key to path
		newPath := path + "/" + key

		switch value.(type) {
		case string:
			// We have a string value, create piece and add to collection
			piece := e.createPiece(newPath, value.(string))
			localData[piece.KVPath] = piece.Value

		case bool:
			// We have a string value, create piece and add to collection
			piece := e.createPiece(newPath, strconv.FormatBool(value.(bool)))
			localData[piece.KVPath] = piece.Value

		case float64:
			// We have a "Javascript number" -> always floating point. Create piece and add to collection
			piece := e.createPiece(newPath, fmt.Sprint(value.(float64)))
			localData[piece.KVPath] = piece.Value

		case []interface{}:
			// We have an array - ohoh
			// Array inside consul are... well are not! Insert as string for now
			piece := e.createPiece(newPath, fmt.Sprint(value.([]interface{})))
			localData[piece.KVPath] = piece.Value

		case map[string]interface{}:
			// we have an object, recurse casting the value
			e.traverseJSON(newPath, value.(map[string]interface{}), localData)
		}
	}
}
