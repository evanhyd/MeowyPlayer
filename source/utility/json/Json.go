package json

import (
	"encoding/json"
	"os"
)

func WriteFile(fileName string, object any) error {
	data, err := json.MarshalIndent(object, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, os.ModePerm)
}

/*
Expect object to be a pointer type
*/
func ReadFile(fileName string, object any) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, object)
}
