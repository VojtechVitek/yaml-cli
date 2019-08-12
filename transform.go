package yaml

import (
	"io"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Transformation struct {
	Matches map[string]string `yaml:"match"`
	Sets    yaml.Node         `yaml:"set"`
	Deletes []string          `yaml:"delete"`
}

func Transformations(r io.Reader) ([]*Transformation, error) {
	var transformations []*Transformation

	dec := yaml.NewDecoder(r)
	for { // For all YAML documents (they might be separated by '---')
		var t Transformation
		if err := dec.Decode(&t); err != nil {
			if err == io.EOF { // Last document.
				break
			}
			return nil, errors.Wrapf(err, "failed to decode YAML")
		}

		transformations = append(transformations, &t)
	}

	return transformations, nil
}

func (t *Transformation) Apply(doc *yaml.Node) error {
	ok, _ := t.MustMatchAll(doc)
	if !ok {
		return nil
	}

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
