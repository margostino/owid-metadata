package common

import "log"

func Check(err error) {
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
}
