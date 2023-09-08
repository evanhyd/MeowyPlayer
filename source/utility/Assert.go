package utility

import (
	"log"
)

// Error must not occur, used in type asserting or asset reading
func MustOk(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

// Error may occur, used in checking if IO fails due to unknown reasons (network error)
func ShouldOk(err error) {
	if err != nil {
		log.Println(err)
	}
}

func MustNotNil(object any) {
	if object == nil {
		log.Panicf("%v is nil\n", object)
	}
}
