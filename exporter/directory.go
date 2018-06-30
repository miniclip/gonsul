package exporter

import (
	"github.com/miniclip/gonsul/structs"

	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func processDir(directory string, localData map[string]string) {
	// Read the entire directory
	files, _ := ioutil.ReadDir(directory)
	// Loop each entry
	for _, file := range files {
		if file.IsDir() {
			// We found a directory, recurse it
			newDir := directory + "/" + file.Name()
			processDir(newDir, localData)
		} else {
			filePath := directory + "/" + file.Name()
			ext := filepath.Ext(filePath)
			if !isExtensionValid(ext) {
				continue
			}
			content, err := ioutil.ReadFile(filePath) // just pass the file name
			if err != nil {
				fmt.Print(err)
			}
			parseFile(filePath, string(content), localData)
		}
	}
}

// isExtensionValid checks if given file extensions is valid for processing
func isExtensionValid(extension string) bool {
	for _, validExtension := range config.GetValidExtensions() {
		if strings.Trim(extension, ".") == strings.Trim(validExtension, ".") {
			return true
		}
	}

	return false
}

func parseFile(filePath string, value string, localData map[string]string) {
	// Extract our file extension and cleanup file path
	ext := filepath.Ext(filePath)
	path := cleanFilePath(filePath)
	// Check if we should parse JSON files
	if config.ShouldExpandJSON() {
		// Check if the file is a JSON one
		if ext == ".json" {
			// Great, we should iterate our JSON (And that's the value)
			expandJSON(path, value, localData)

			// we must return here, to avoid importing the file as blob
			return
		}
	}

	// Not expanding JSON files, create new single "piece" with the
	// value given (the file content) and add to collection
	piece := createPiece(path, value)
	localData[piece.KVPath] = piece.Value
}

func cleanFilePath(filePath string) string {
	// Set part of the config that should be removed from the current
	// file system path in order to build our final Consul KV path
	replace := config.GetRepoRootDir() + "/" + config.GetRepoBasePath()
	// Remove the above from the file system path
	entryFilePath := strings.Replace(filePath, replace, "", 1)
	// Remove any left slash
	entryFilePath = strings.Replace(entryFilePath, "/", "", 1)
	// Remove the file extension from the file system path
	entryFilePath = strings.TrimSuffix(entryFilePath, filepath.Ext(entryFilePath))

	return entryFilePath
}

func createPiece(path string, value string) structs.Entry {
	// Create our Consul base path variable
	var kvPath string

	// Check if we have a Consul KV base path
	if config.GetConsulbasePath() != "" {
		kvPath = config.GetConsulbasePath()
	}

	// Finally append the Consul KV base path to the file path, if base is not an empty string
	if kvPath != "" {
		return structs.Entry{KVPath: kvPath + "/" + path, Value: value}
	}

	return structs.Entry{KVPath: path, Value: value}
}
