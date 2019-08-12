package yaml

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Node = yaml.Node

func Parse(in []byte) (*Node, error) {
	var doc Node
	if err := yaml.Unmarshal(in, &doc); err != nil {
		return nil, errors.Wrap(err, "failed to parse YAML")
	}
	return &doc, nil
}

func Write(w io.Writer, in *Node) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	return enc.Encode(in)
}

func Bytes(in *Node) []byte {
	var b bytes.Buffer
	_ = Write(&b, in)
	return b.Bytes()
}
