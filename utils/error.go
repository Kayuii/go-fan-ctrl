package utils

import "log"

func Fatal(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: error '%+v'", msg, err)
	}
}
