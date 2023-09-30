package json

import (
	"encoding/json"
	"os"
)

func Write(fileName string, object any) error {
	data, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, os.ModePerm)
}

func Read(fileName string, object any) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, object)
}
