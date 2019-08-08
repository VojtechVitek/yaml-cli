package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/VojtechVitek/yaml"
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
		selector := os.Args[2]

		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read stdin")
		}

		buf, err := yaml.Delete(in, selector)
		if err != nil {
			if false { // TODO: --strict mode, where we'd error out on non-existent selectors?
				return errors.Wrapf(err, "failed to delete %q", selector)
			}
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
