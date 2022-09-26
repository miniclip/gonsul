package importer

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func (i *importer) exportToDirectory(directoryPath string, data map[string]string) error {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := data[key]
		filePath := filepath.Join(directoryPath, key+".txt")
		dirPath := filepath.Dir(filePath)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, 0700)
			if err != nil {
				return err
			}
		}
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		data, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return err
		}
		fmt.Println(key)
		file.WriteString(string(data))
	}
	return nil
}
