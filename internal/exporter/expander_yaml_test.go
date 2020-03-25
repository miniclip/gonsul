package exporter

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

func Test_exporter_validateYAML(t *testing.T) {
	var arbitraryJSON map[string]interface{}

	err := yaml.Unmarshal([]byte(data), &arbitraryJSON)

	fmt.Printf("%+v", err)
	fmt.Printf("%+v", arbitraryJSON)
}
