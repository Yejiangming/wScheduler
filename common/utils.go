package common

import "log"

type Result struct {
	Success bool
	Message string
}

func PanicIf(err error) {
	if err != nil {
		log.Panic(err)
	}
}
