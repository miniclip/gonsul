package exporter

import (
	"github.com/miniclip/gonsul/internal/entities"

	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
)

// traverse is our entry point function to start traversing a given directory.
// this is a recursive function, as it will call itself whenever we hit a sub folder
func (e *exporter) parseDir(directory string, localData map[string]string) {
	// Read the entire directory
	files, _ := ioutil.ReadDir(directory)
	// Loop each entry
	for _, file := range files {
		if file.IsDir() {
			// We found a directory, recurse it
			newDir := directory + "/" + file.Name()
			e.parseDir(newDir, localData)
		} else {
			filePath := directory + "/" + file.Name()
			ext := filepath.Ext(filePath)
			if !e.isExtensionValid(ext) {
				continue
			}
			content, err := ioutil.ReadFile(filePath) // just pass the file name
			if err != nil {
				fmt.Print(err)
			}
			e.parseFile(filePath, string(content), localData)
		}
	}
}

// isExtensionValid checks if given file extensions is valid for processing
func (e *exporter) isExtensionValid(extension string) bool {
	for _, validExtension := range e.config.GetValidExtensions() {
		if strings.Trim(extension, ".") == strings.Trim(validExtension, ".") {
			return true
		}
	}

	return false
}

// parseFile ...
func (e *exporter) parseFile(filePath string, value string, localData map[string]string) {
	// Extract our file extension and cleanup file path
	ext := filepath.Ext(filePath)
	cleanedPath := e.cleanFilePath(filePath)

	// Check if the file is a JSON one
	if ext == ".json" {
		// Check if we should parse JSON files
		if e.config.ShouldExpandJSON() {
			// Great, we should iterate our JSON (And that's the value)
			e.expandJSON(cleanedPath, value, localData)

			// we must return here, to avoid importing the file as blob
			return
		}

		// Not expanding json file, but we should validate anyways
		// HEADS UP: Below function will exit program if any error found
		_ = e.validateJSON(cleanedPath, value)
	}

	// Check if the file is a YAML one
	if ext == ".yaml" {
		// Check if we should parse JSON files
		if e.config.ShouldExpandYAML() {
			// Great, we should iterate our JSON (And that's the value)
			e.expandYAML(cleanedPath, value, localData)

			// we must return here, to avoid importing the file as blob
			return
		}

		// Not expanding json file, but we should validate anyways
		// HEADS UP: Below function will exit program if any error found
		_ = e.validateYAML(cleanedPath, value)
	}

	// Not expanding JSON files, create new single "piece" with the
	// value given (the file content) and add to collection
	piece := e.createPiece(cleanedPath, value)
	localData[piece.KVPath] = piece.Value
}

// cleanFilePath ...
func (e *exporter) cleanFilePath(filePath string) string {
	// Set part of the config that should be removed from the current
	// file system path in order to build our final Consul KV path
	replace := path.Join(e.config.GetRepoRootDir(), e.config.GetRepoBasePath())
	// Remove the above from the file system path
	entryFilePath := strings.Replace(filePath, replace, "", 1)
	// Remove any left slash
	entryFilePath = strings.Replace(entryFilePath, "/", "", 1)
	// Set or not the file extension when importing to consul k/v the file
	if !e.config.KeepFileExt() {
		entryFilePath = strings.TrimSuffix(entryFilePath, filepath.Ext(entryFilePath))
	}

	return entryFilePath
}

// createPiece ...
func (e *exporter) createPiece(piecePath string, value string) entities.Entry {
	// Create our Consul base path variable
	var kvPath string

	// Check if we have a Consul KV base path
	if e.config.GetConsulBasePath() != "" {
		kvPath = e.config.GetConsulBasePath()
	}

	// Finally append the Consul KV base path to the file path, if base is not an empty string
	if kvPath != "" {
		fullPath := path.Join(kvPath, piecePath)
		return entities.Entry{KVPath: fullPath, Value: value}
	}

	return entities.Entry{KVPath: piecePath, Value: value}
}
