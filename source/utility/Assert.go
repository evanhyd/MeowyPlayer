package utility

import (
	"log"
)

func MustOk(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func MustNotNil(object any) {
	if object == nil {
		log.Panicf("%v is nil", object)
	}
}
