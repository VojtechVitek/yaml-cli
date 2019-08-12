package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/VojtechVitek/yaml"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	yamlv3 "gopkg.in/yaml.v3"
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

	dec := yamlv3.NewDecoder(os.Stdin)
	enc := yamlv3.NewEncoder(os.Stdout)
	enc.SetIndent(2)

	for { // For all YAML documents (they might be separated by '---')
		var doc yamlv3.Node
		if err := dec.Decode(&doc); err != nil {
			if err == io.EOF { // Last document.
				return nil
			}
			return errors.Wrap(err, "failed to decode YAML document(s) from stdin")
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
				ok, _ := tf.MustMatchAll(&doc)
				if !ok {
					continue
				}

				if err := tf.Apply(&doc); err != nil {
					return errors.Wrapf(err, "failed to apply transformation %v", filenames[i])
				}
			}

		case "grep":
			selectors := os.Args[2:]

			// yaml grep -v "something"
			invert := false
			if selectors[0] == "-v" {
				invert = true
				selectors = selectors[1:]
			}

			tf, err := yaml.ParseTransformation([]byte(fmt.Sprintf("match:\n  %v", strings.Join(selectors, "\n  "))))
			if err != nil {
				return errors.Wrapf(err, "failed to parse grep `key: value' pairs from %q", selectors)
			}

			ok, _ := tf.MustMatchAll(&doc)
			if ok == invert {
				continue // Do not print anything out.
			}

			// Print out original doc.

		case "delete":
			selector := os.Args[2]

			if err := yaml.Delete(&doc, selector); err != nil {
				if false { // TODO: --strict mode, where we'd error out on non-existent selectors?
					return errors.Wrapf(err, "failed to delete %q", selector)
				}
			}

		default:
			return errors.Errorf("%v: unknown command", os.Args[1])
		}

		if err := enc.Encode(&doc); err != nil {
			return errors.Wrap(err, "failed to write to stdout")
		}
	}

	return nil
}
