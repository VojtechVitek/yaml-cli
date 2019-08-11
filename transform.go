package yaml

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Transformation struct {
	Matches map[string]string `yaml:"match"`
	Sets    yaml.Node         `yaml:"set"`
	Deletes []string          `yaml:"delete"`
}

func ParseTransformation(b []byte) (*Transformation, error) {
	var t Transformation
	if err := yaml.Unmarshal(b, &t); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal transformation")
	}

	return &t, nil
}

func (t *Transformation) Apply(doc *yaml.Node) error {
	if err := t.ApplySet(doc); err != nil {
		return err
	}

	return nil
}

func (t *Transformation) ApplySet(doc *yaml.Node) error {
	for i := 0; i < len(t.Sets.Content); i += 2 {
		path, value := t.Sets.Content[i].Value, t.Sets.Content[i+1]

		if err := Set(doc, path, value); err != nil {
			return errors.Wrapf(err, "failed to set %q", path)
		}
	}

	return nil
}
