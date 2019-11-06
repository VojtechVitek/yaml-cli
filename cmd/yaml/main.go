package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/VojtechVitek/yaml/pkg/cli"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

var (
	flags = flag.NewFlagSet("yaml", flag.ExitOnError)
	from  = flags.String("from", "yaml", "input data format [json]")
	to    = flags.String("to", "yaml", "output data format [json]")
)

func main() {
	var (
		in  io.Reader
		out io.Writer
		b   bytes.Buffer
	)

	flags.Parse(os.Args[1:])

	switch *from {
	case "yaml", "yml": // Nop.
		in = os.Stdin
	case "json":
		in = &b
		if err := jsonToYAML(&b, os.Stdin); err != nil {
			log.Fatal(errors.Wrap(err, "failed to convert output to JSON"))
		}
	default:
		log.Fatalf("unknown --from format %q", *to)
	}

	switch *to {
	case "yaml", "yml": // Nop.
		out = os.Stdout
	case "json":
		out = &b
	default:
		log.Fatalf("unknown --to format %q", *to)
	}

	if err := cli.Run(out, in, flags.Args()); err != nil {
		log.Fatal(err)
	}

	if *to == "json" {
		if err := yamlToJSON(os.Stdout, &b); err != nil {
			log.Fatal(errors.Wrap(err, "failed to convert output to JSON"))
		}
	}
}
