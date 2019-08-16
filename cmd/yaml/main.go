package main

import (
	"log"
	"os"

	"github.com/VojtechVitek/yaml/pkg/cli"
)

func main() {
	if err := cli.Run(os.Stdout, os.Stdin, os.Args); err != nil {
		log.Fatal(err)
	}
}
