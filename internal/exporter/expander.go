package exporter

import (
	"fmt"
	"strconv"
)

// traverseJSON ...
func (e *exporter) traverseMap(path string, arbitraryJSON map[string]interface{}, localData map[string]string) {
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
			e.traverseMap(newPath, value.(map[string]interface{}), localData)
		}
	}
}
