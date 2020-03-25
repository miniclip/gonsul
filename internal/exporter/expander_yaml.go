package exporter

import (
	"github.com/miniclip/gonsul/internal/util"

	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
)

// validateJSON ...
func (e *exporter) validateYAML(path string, yamlData string) map[string]interface{} {
	// Create "generic" yaml struct
	var arbitraryYAML map[string]interface{}

	// Decode data into "generic"
	err := yaml.Unmarshal([]byte(yamlData), &arbitraryYAML)

	// Decoded YAML ok?
	if err != nil {
		util.ExitError(
			errors.New(fmt.Sprintf("error parsing YAML file: %s with Message: %s", path, err.Error())),
			util.ErrorFailedJsonDecode,
			e.logger,
		)
	}

	return arbitraryYAML
}

// expandJSON ...
func (e *exporter) expandYAML(path string, jsonData string, localData map[string]string) {
	arbitraryYAML := e.validateYAML(path, jsonData)

	// Iterate over our "generic" JSON structure
	e.traverseMap(path, arbitraryYAML, localData)
}
