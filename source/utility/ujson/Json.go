package ujson

import (
	"encoding/json"
	"os"
)

func WriteFile(fileName string, object any) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(object)
}

/*
Object must be a pointer type
*/
func ReadFile(fileName string, object any) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(object)
}
