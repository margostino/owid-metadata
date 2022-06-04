package main

import (
	"fmt"
	"github.com/margostino/owid-metadata/tooling"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		message := fmt.Sprintf("Command Not Found!\n" +
			"Commands available: \n" +
			"- server: to start a local server\n" +
			"- schema-gen: to generate a new Graphql Schema")
		log.Panicln(message)
	}
	action := os.Args[1]
	if action == "metadata-gen" {
		tooling.GenerateMetadata()
	} else {
		message := fmt.Sprintf("Command Not Valid!\n")
		log.Panicln(message)
	}
}
