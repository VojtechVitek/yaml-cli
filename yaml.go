package yaml

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Transformation struct {
	Matches map[string]string `yaml:"match"`
	Sets    map[string]string `yaml:"set"`
	Deletes []string          `yaml:"delete"`
}

func NewTransformation(b []byte) (*Transformation, error) {
	var t Transformation

	err := yaml.Unmarshal(b, &t)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal input")
	}

	return &t, nil
}

// TODO: Provide helper func to parse input doc instead of
// duplicating upstream types/methods.
type Node = yaml.Node

func Marshal(in interface{}) ([]byte, error)     { return yaml.Marshal(in) }
func Unmarshal(in []byte, out interface{}) error { return yaml.Unmarshal(in, out) }
