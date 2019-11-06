package cli

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/VojtechVitek/yaml"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	yamlv3 "gopkg.in/yaml.v3"
)

// TODO: Split into multiple files/functions. This function grew too much
//       over time while adding new commands and functionality.
func Run(out io.Writer, in io.Reader, args []string) error {
	if len(args) == 1 && args[0] == "count" {
		count := 0
		dec := yamlv3.NewDecoder(in)
		for { // For all YAML documents in STDIN.
			var doc yamlv3.Node
			if err := dec.Decode(&doc); err != nil {
				if err == io.EOF { // Last document.
					break
				}
				return errors.Wrap(err, "failed to decode YAML document(s) from stdin")
			}
			count++
		}
		fmt.Println(count)
		return nil
	}

	if len(args) <= 1 {
		return errors.New(`usage: see https://github.com/VojtechVitek/yaml-cli/blob/master/README.md`)
	}

	enc := yamlv3.NewEncoder(out)
	enc.SetIndent(2)

	if args[0] == "cat" {
		for _, filename := range args[1:] {
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
					return errors.Wrap(err, "failed to encode YAML node")
				}
			}

			if err := f.Close(); err != nil {
				return errors.Wrap(err, "failed to close file")
			}
		}
		return nil
	}

	var tfs []*yaml.Transformation

	switch args[0] {
	case "apply":
		filenames := args[1:]

		fileTfs := make([][]*yaml.Transformation, len(filenames))

		var g errgroup.Group
		for i, filename := range filenames {
			i, filename := i, filename // local copy for goroutine
			g.Go(func() error {
				f, err := os.Open(filename)
				if err != nil {
					return errors.Wrapf(err, "failed to read transformation %v", filename)
				}

				fileTfs[i], err = yaml.Transformations(f)
				if err != nil {
					return errors.Wrapf(err, "failed to parse transformation %v", filename)
				}

				return f.Close()
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
		for _, tf := range fileTfs {
			tfs = append(tfs, tf...)
		}
		// do not return

	case "doc":
		if len(args) != 2 {
			return errors.New("usage: yaml doc $index")
		}

		index, err := strconv.Atoi(args[1])
		if err != nil {
			return errors.Wrapf(err, "yaml doc $index: failed to parse $index %q", args[1])
		}

		count := 0
		dec := yamlv3.NewDecoder(in)
		for { // For all YAML documents in STDIN.
			var doc yamlv3.Node
			if err := dec.Decode(&doc); err != nil {
				if err == io.EOF { // Last document.
					break
				}
				return errors.Wrap(err, "failed to decode YAML document(s) from stdin")
			}
			if count == index {
				if err := enc.Encode(&doc); err != nil {
					return errors.Wrap(err, "failed to encode YAML node")
				}
				io.Copy(ioutil.Discard, in)
				return nil
			}
			count++
		}
		return errors.Errorf("doc %v not found (count=%v)", index, count)
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

		switch args[0] {
		case "apply":
			for _, tf := range tfs {
				if err := tf.Apply(&doc); err != nil {
					return errors.Wrapf(err, "failed to apply transformation")
				}
			}

			if err := enc.Encode(&doc); err != nil {
				return errors.Wrap(err, "failed to encode YAML node")
			}

		case "match":
			for _, selector := range args[1:] {
				_, err := yaml.Get(&doc, selector)
				if err != nil {
					return errors.Wrapf(err, "failed to get %q", selector)
				}
			}

			return nil

		case "grep":
			selectors := args[1:]

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
			if ok != invert { // match
				if err := enc.Encode(&doc); err != nil {
					return errors.Wrap(err, "failed to encode YAML node")
				}
			}

		case "print":
			obj := &yamlv3.Node{
				Kind: yamlv3.MappingNode,
				Tag:  "!!map",
			}

			for _, selector := range args[1:] {
				node, err := yaml.Get(&doc, selector)
				if err != nil {
					return errors.Wrapf(err, "failed to get %q", selector)
				}

				obj.Content = append(obj.Content,
					&yamlv3.Node{
						Kind:  yamlv3.ScalarNode,
						Value: selector,
					},
					node,
				)
			}

			if err := enc.Encode(obj); err != nil {
				return errors.Wrap(err, "failed to encode YAML node")
			}

		case "len":
			selector := args[1]

			node, err := yaml.Get(&doc, selector)
			if err != nil {
				return errors.Wrapf(err, "failed to get %q", selector)
			}

			fmt.Println(len(node.Content))
			return nil

		case "get":
			for _, selector := range args[1:] {
				// TODO: Get() []node .... get array[*]
				// TODO: encode each node in a loop
				node, err := yaml.Get(&doc, selector)
				if err != nil {
					return errors.Wrapf(err, "failed to get %q", selector)
				}

				// Don't reuse top level encoder; we don't want to render
				// multiple YAML documents separated by `---`.
				enc = yamlv3.NewEncoder(out)
				enc.SetIndent(2)
				if err := enc.Encode(node); err != nil {
					return errors.Wrap(err, "failed to encode YAML node")
				}
			}

		case "set":
			keyValues := args[1:]

			tfs, err := yaml.Transformations(strings.NewReader(fmt.Sprintf("set:\n  %v", strings.Join(keyValues, "\n  "))))
			if err != nil {
				return errors.Wrapf(err, "failed to parse `key: value' pairs from %v", keyValues)
			}
			tf := tfs[0]

			if err := tf.Apply(&doc); err != nil {
				return errors.Wrapf(err, "failed to set %v", tf.Sets)
			}

			if err := enc.Encode(&doc); err != nil {
				return errors.Wrap(err, "failed to encode YAML node")
			}

		case "default":
			keyValues := args[1:]

			tfs, err := yaml.Transformations(strings.NewReader(fmt.Sprintf("default:\n  %v", strings.Join(keyValues, "\n  "))))
			if err != nil {
				return errors.Wrapf(err, "failed to parse `key: value' pairs from %v", keyValues)
			}
			tf := tfs[0]

			if err := tf.Apply(&doc); err != nil {
				return errors.Wrapf(err, "failed to set default %v", tf.Defaults)
			}

			if err := enc.Encode(&doc); err != nil {
				return errors.Wrap(err, "failed to encode YAML node")
			}

		case "delete":
			selector := args[1]

			if err := yaml.Delete(&doc, selector); err != nil {
				if false { // TODO: --strict mode, where we'd error out on non-existent selectors?
					return errors.Wrapf(err, "failed to delete %q", selector)
				}
			}

		default:
			return errors.Errorf("%v: unknown command", args[0])
		}
	}

	return nil
}
