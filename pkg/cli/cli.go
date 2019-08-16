package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/VojtechVitek/yaml"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	yamlv3 "gopkg.in/yaml.v3"
)

func Run(out io.Writer, in io.Reader, args []string) error {
	if len(args) <= 2 {
		return errors.New("usage: yaml apply [files..]")
	}

	enc := yamlv3.NewEncoder(out)
	enc.SetIndent(2)

	if args[1] == "cat" {
		for _, filename := range args[2:] {
			f, err := os.Open(filename)
			if err != nil {
				return errors.Wrap(err, "failed to open file")
			}

			dec := yamlv3.NewDecoder(f)
			for { // For all YAML documents in the file (they might be separated by '---')
				var doc yamlv3.Node
				if err := dec.Decode(&doc); err != nil {
					if err == io.EOF { // Last document.
						break
					}
					return errors.Wrapf(err, "failed to decode file %v", filename)
				}

				if err := enc.Encode(&doc); err != nil {
					return errors.Wrap(err, "failed to write to stdout")
				}
			}

			if err := f.Close(); err != nil {
				return errors.Wrap(err, "failed to close file")
			}
		}
		return nil
	}

	dec := yamlv3.NewDecoder(in)
	for { // For all YAML documents in STDIN.
		var doc yamlv3.Node
		if err := dec.Decode(&doc); err != nil {
			if err == io.EOF { // Last document.
				return nil
			}
			return errors.Wrap(err, "failed to decode YAML document(s) from stdin")
		}

		switch args[1] {
		case "apply":
			filenames := args[2:]
			tfs := make([][]*yaml.Transformation, len(filenames))

			var g errgroup.Group
			for i, filename := range filenames {
				i, filename := i, filename
				g.Go(func() error {
					f, err := os.Open(filename)
					if err != nil {
						return errors.Wrapf(err, "failed to read transformation %v", filename)
					}

					tfs[i], err = yaml.Transformations(f)
					if err != nil {
						return errors.Wrapf(err, "failed to parse transformation %v", filename)
					}

					return nil
				})
			}
			if err := g.Wait(); err != nil {
				return err
			}

			for i, filename := range filenames {
				for _, tf := range tfs[i] {
					if err := tf.Apply(&doc); err != nil {
						return errors.Wrapf(err, "failed to apply %v. transformation from %v", i+1, filename)
					}
				}
			}

		case "grep":
			selectors := args[2:]

			// yaml grep -v "something"
			invert := false
			if selectors[0] == "-v" {
				invert = true
				selectors = selectors[1:]
			}

			tfs, err := yaml.Transformations(strings.NewReader(fmt.Sprintf("match:\n  %v", strings.Join(selectors, "\n  "))))
			if err != nil {
				return errors.Wrapf(err, "failed to parse grep `key: value' pairs from %q", selectors)
			}
			tf := tfs[0]

			ok, _ := tf.MustMatchAll(&doc)
			if ok == invert {
				continue // Do not print anything out.
			}

			// Print out original doc.

		case "get":
			selector := args[2]

			node, err := yaml.Get(&doc, selector)
			if err != nil {
				return errors.Wrapf(err, "failed to get %q", selector)
			}

			// Don't reuse top level encoder; we don't want to render
			// multiple YAML documents separated by `---`.
			enc := yamlv3.NewEncoder(out)
			enc.SetIndent(2)

			if err := enc.Encode(node); err != nil {
				return errors.Wrap(err, "failed to write to stdout")
			}

			continue

		case "delete":
			selector := args[2]

			if err := yaml.Delete(&doc, selector); err != nil {
				if false { // TODO: --strict mode, where we'd error out on non-existent selectors?
					return errors.Wrapf(err, "failed to delete %q", selector)
				}
			}

		default:
			return errors.Errorf("%v: unknown command", args[1])
		}

		if err := enc.Encode(&doc); err != nil {
			return errors.Wrap(err, "failed to write to stdout")
		}
	}

	return nil
}
