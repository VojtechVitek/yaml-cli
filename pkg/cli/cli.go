package cli

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/VojtechVitek/yaml-cli"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
	"golang.org/x/sync/errgroup"
	yamlv3 "gopkg.in/yaml.v3"
)

var (
	// TODO: Move these vars to a struct, they shouldn't be global. A left-over from "main" pkg.
	flags         = flag.NewFlagSet("yaml", flag.ExitOnError)
	from          = flags.String("from", "yaml", "input data format [json]")
	to            = flags.String("to", "yaml", "output data format [json]")
	ignoreMissing = flags.Bool("ignore-missing", false, "yaml get: ignore missing nodes")
	printKey      = flags.Bool("print-key", false, "yaml get: print node key in front of the value, so the output is valid YAML")
	noSeparator   = flags.Bool("no-separator", false, "yaml get: don't print `---' separator between YAML documents")
	invert        = flags.BoolP("invert-match", "v", false, "yaml grep -v: select non-matching documents")
)

// TODO: Split into multiple files/functions. This function grew too much
//       over time while adding new commands and functionality.
func Run(out io.Writer, in io.Reader, args []string) error {
	var b bytes.Buffer
	flags.Parse(args)

	switch *from {
	case "yaml", "yml": // Nop.
	case "json":
		in = &b
		if err := jsonToYAML(&b, os.Stdin); err != nil {
			log.Fatal(errors.Wrap(err, "failed to convert output to JSON"))
		}
	default:
		return errors.Errorf("unknown --from format %q", *to)
	}

	switch *to {
	case "yaml", "yml": // Nop.
	case "json":
		out = &b
	default:
		return errors.Errorf("unknown --to format %q", *to)
	}

	if err := run(out, in, flags.Args()); err != nil {
		return errors.Wrap(err, "failed to run")
	}

	if *to == "json" {
		if err := yamlToJSON(os.Stdout, &b); err != nil {
			return errors.Wrap(err, "failed to convert output to JSON")
		}
	}

	return nil
}

func run(out io.Writer, in io.Reader, args []string) error {
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

		case "grep":
			selectors := args[1:]

			tfs, err := yaml.Transformations(strings.NewReader(fmt.Sprintf("match:\n  %v", strings.Join(selectors, "\n  "))))
			if err != nil {
				return errors.Wrapf(err, "failed to parse grep `key: value' pairs from %q", selectors)
			}
			tf := tfs[0]

			ok, _ := tf.MustMatchAll(&doc)
			//fmt.Fprintf(os.Stderr, "MustMatchAll: %v\n", err)
			if ok != *invert { // match
				if err := enc.Encode(&doc); err != nil {
					return errors.Wrap(err, "failed to encode YAML node")
				}
			}

		case "len":
			selector := args[1]

			nodes, err := yaml.Get(&doc, strings.Split(selector, "."))
			if err != nil {
				return errors.Wrapf(err, "failed to get %q", selector)
			}

			for _, node := range nodes {
				fmt.Println(len(node.Content))
			}
			return nil

		case "get":
			useTopLevelEnc := true

			var allMatchedNodes []*yamlv3.Node

			for _, selector := range args[1:] {
				selectors := strings.Split(selector, ".")
				lastSelector := selectors[len(selectors)-1]

				nodes, err := yaml.Get(&doc, selectors)
				if err != nil {
					if !*ignoreMissing {
						return errors.Wrapf(err, "failed to get %q", selector)
					}
				}

				for _, node := range nodes {
					if *printKey {
						node = &yamlv3.Node{
							Kind: yamlv3.MappingNode,
							Tag:  "!!map",
							Content: []*yamlv3.Node{
								&yamlv3.Node{
									Kind:  yamlv3.ScalarNode,
									Value: lastSelector,
								},
								node,
							},
						}
					}

					allMatchedNodes = append(allMatchedNodes, node)
				}
			}

			for _, node := range allMatchedNodes {
				if useTopLevelEnc {
					useTopLevelEnc = false
					// Top level enc will print separator between documents.
					if err := enc.Encode(node); err != nil {
						return errors.Wrap(err, "failed to encode YAML node")
					}
				} else {
					// Omit separator between one YAML document's nodes.
					enc := yamlv3.NewEncoder(out)
					enc.SetIndent(2)
					if err := enc.Encode(node); err != nil {
						return errors.Wrap(err, "failed to encode YAML node")
					}
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
