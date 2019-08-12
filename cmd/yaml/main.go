package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/VojtechVitek/yaml"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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

	in, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "failed to read stdin")
	}

	doc, err := yaml.Parse(in)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal input")
	}

	switch os.Args[1] {
	case "apply":
		filenames := os.Args[2:]
		transformations := make([]*yaml.Transformation, len(filenames))

		var g errgroup.Group
		for i, filename := range filenames {
			i, filename := i, filename
			g.Go(func() error {
				b, err := ioutil.ReadFile(filename)
				if err != nil {
					return errors.Wrapf(err, "failed to read transformation %v", filename)
				}

				transformations[i], err = yaml.ParseTransformation(b)
				if err != nil {
					return errors.Wrapf(err, "failed to parse transformation %v", filename)
				}

				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}

		for i, tf := range transformations {
			ok, _ := tf.MustMatchAll(doc)
			if !ok {
				continue
			}

			if err := tf.Apply(doc); err != nil {
				return errors.Wrapf(err, "failed to apply transformation %v", filenames[i])
			}
		}

	case "grep":
		arg := 2
		invert := false

		// yaml grep -v "something"
		if os.Args[arg] == "-v" {
			invert = true
			arg++
		}

		// yaml grep -v "something"
		selector := os.Args[arg]

		tf, err := yaml.ParseTransformation([]byte(fmt.Sprintf("match:\n  %v", selector)))
		if err != nil {
			return errors.Wrapf(err, "failed to parse grep `key: value' pair from %q", selector)
		}

		ok, _ := tf.MustMatchAll(doc)
		if ok == invert {
			return nil // Return; do not print anything out.
		}

		// Print out original doc.

	case "delete":
		selector := os.Args[2]

		if err := yaml.Delete(doc, selector); err != nil {
			if false { // TODO: --strict mode, where we'd error out on non-existent selectors?
				return errors.Wrapf(err, "failed to delete %q", selector)
			}
		}

	default:
		return errors.Errorf("%v: unknown command", os.Args[1])
	}

	err = yaml.Write(os.Stdout, doc)
	if err != nil {
		return errors.Wrap(err, "failed to write to stdout")
	}

	return nil
}
