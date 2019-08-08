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
	if len(os.Args) <= 2 {
		return errors.New("usage: yaml apply [files..]")
	}

	switch os.Args[1] {
	case "apply":
		filenames := os.Args[2:]
		transformations := make([]*yaml.Transformation, len(filenames))

		for i, filename := range filenames {
			b, err := ioutil.ReadFile(filename)
			if err != nil {
				return errors.Wrapf(err, "failed to read transformation %v", filename)
			}

			transformations[i], err = yaml.NewTransformation(b)
			if err != nil {
				return errors.Wrapf(err, "failed to parse transformation %v", filename)
			}
		}

		in, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "failed to read stdin")
		}

		var doc yaml.Node
		err = yaml.Unmarshal(in, &doc) // rename to yaml.Input()
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal input")
		}

		for _, tf := range transformations {
			ok, err := tf.MatchesAll(&doc, tf.Matches)
			if ok {
				log.Printf("WE HAVE A MATCH")
			}
			log.Println(err)
		}

		return nil

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
