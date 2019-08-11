package yaml

import (
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

func Bytes(in *Node) []byte {
	b, _ := yaml.Marshal(in)
	return b
}
