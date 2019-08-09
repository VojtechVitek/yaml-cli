package yaml

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Transformation struct {
	Matches map[string]string    `yaml:"match"`
	Sets    map[string]yaml.Node `yaml:"set"`
	Deletes []string             `yaml:"delete"`
}

func NewTransformation(b []byte) (*Transformation, error) {
	var t Transformation

	err := yaml.Unmarshal(b, &t)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal input")
	}

	return &t, nil
}

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
