package ujson

import (
	"encoding/json"
	"os"
)

/*
Write the object into file.
*/
func Write(fileName string, object any) error {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(object)
}

/*
Read the file into the object.
Object must be a pointer type.
*/
func Read(fileName string, object any) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(object)
}
