package importer

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
		file.WriteString(string(data))
	}
	return nil
}

func (i *importer) exportToFile(filePath string, data map[string]string, base64Encode bool) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("opening dict file for write filename=%s err=%v", filePath, err)
	}
	defer f.Close()

	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := data[key]
		var saveStr string
		if base64Encode {
			saveStr = key + "=" + value
		} else {
			decodedValue, err := base64.StdEncoding.DecodeString(value)
			if err != nil {
				return fmt.Errorf("error decoding key to %s err=%v", key, err)
			}
			decodedStr := strings.Replace(string(decodedValue), "\n", "\n+", -1)
			decodedStr = strings.Replace(decodedStr, "\r\n+", "\nr+", -1)
			saveStr = key + "=" + decodedStr
		}
		_, err := f.WriteString(saveStr)
		if err != nil {
			return fmt.Errorf("error writting to %s err=%v", filePath, err)
		}
		_, err = f.WriteString("\n")
		if err != nil {
			return fmt.Errorf("error writting to %s err=%v", filePath, err)
		}
	}
	f.Sync()
	return nil
}
