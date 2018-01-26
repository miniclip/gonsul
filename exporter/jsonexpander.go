package exporter

import (
	"github.com/miniclip/gonsul/errorutil"

	"encoding/json"
	"strconv"
	"errors"
	"fmt"
)

func expandJSON(path string, jsonData string, localData map[string]string) {
	// Create "generic" json struct
	var arbitraryJSON map[string]interface{}

	// Decode data into "generic"
	err := json.Unmarshal([]byte(jsonData), &arbitraryJSON)

	// Decoded JSON ok?
	if err != nil {
		errorutil.ExitError(
			errors.New(fmt.Sprintf("error parsing JSON file: %s", err.Error())),
			errorutil.ErrorFailedJsonDecode,
			&logger,
			)
	}

	// Iterate over our "generic" JSON structure
	traverseJSON(path, arbitraryJSON, localData)
}

func traverseJSON(path string, arbitraryJSON map[string]interface{}, localData map[string]string) {
	for key, value := range arbitraryJSON {
		// Append key to path
		newPath := path + "/" + key

		switch value.(type) {
		case string:
			// We have a string value, create piece and add to collection
			piece := createPiece(newPath, value.(string))
			localData[piece.KVPath] = piece.Value

		case bool:
			// We have a string value, create piece and add to collection
			piece := createPiece(newPath, strconv.FormatBool(value.(bool)))
			localData[piece.KVPath] = piece.Value

		case float64:
			// We have a "Javascript number" -> always floating point. Create piece and add to collection
			piece := createPiece(newPath, fmt.Sprint(value.(float64)))
			localData[piece.KVPath] = piece.Value

		case []interface{}:
			// We have an array - ohoh
			// Array inside consul are... well are not! Insert as string for now
			piece := createPiece(newPath, fmt.Sprint(value.([]interface{})))
			localData[piece.KVPath] = piece.Value

		case map[string]interface{}:
			// we have an object, recurse casting the value
			traverseJSON(newPath, value.(map[string]interface{}), localData)
		}
	}
}