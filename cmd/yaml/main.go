package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/VojtechVitek/transform"
	"github.com/pkg/errors"
)

func main() {
	if err := runCLI(); err != nil {
		log.Fatal(err)
	}
}

func runCLI() error {
	if len(os.Args) != 3 {
		return errors.New("expected one command and argument")
	}

	switch os.Args[1] {
	case "delete":
		delete := os.Args[2]

		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read stdin")
		}

		transformer := &transform.Transformer{}
		buf, err := transformer.Delete(input, delete)
		if err != nil {
			return errors.Wrapf(err, "failed to delete %q", delete)
		}

		_, err = os.Stdout.Write(buf)
		if err != nil {
			return errors.Wrap(err, "failed to write to stdout")
		}
		return nil

	default:
		return errors.Errorf("%v: unknown command")
	}
}
