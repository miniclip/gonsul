package exporter

import (
	"github.com/miniclip/gonsul/errorutil"
	"github.com/miniclip/gonsul/data"
	"encoding/json"
	"strconv"
	"errors"
	"fmt"
)

func expandJSON(path string, jsonData string, importData *data.EntryCollection) {
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
	traverseJSON(path, arbitraryJSON, importData)
}

func traverseJSON(path string, arbitraryJSON map[string]interface{}, importData *data.EntryCollection) {
	for key, value := range arbitraryJSON {
		// Append key to path
		newPath := path + "/" + key

		switch value.(type) {
		case string:
			// We have a string value, create piece and add to collection
			piece := createPiece(newPath, value.(string))
			importData.AddEntry(piece)

		case bool:
			// We have a string value, create piece and add to collection
			piece := createPiece(newPath, strconv.FormatBool(value.(bool)))
			importData.AddEntry(piece)

		case float64:
			// We have a "Javascript number" -> always floating point. Create piece and add to collection
			piece := createPiece(newPath, fmt.Sprint(value.(float64)))
			importData.AddEntry(piece)

		case []interface{}:
			// We have an array - ohoh
			// Array inside consul are... well are not! Insert as string for now
			piece := createPiece(newPath, fmt.Sprint(value.([]interface{})))
			importData.AddEntry(piece)

		case map[string]interface{}:
			// we have an object, recurse casting the value
			traverseJSON(newPath, value.(map[string]interface{}), importData)
		}
	}
}