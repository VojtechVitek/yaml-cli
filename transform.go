package yaml

import (
	"io"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Transformation struct {
	Matches  map[string]yaml.Node `yaml:"match"`
	Sets     yaml.Node            `yaml:"set"`
	Defaults yaml.Node            `yaml:"default"`
	Deletes  []string             `yaml:"delete"`
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

	if err := t.ApplyDeletes(doc); err != nil {
		return err
	}

	if err := t.ApplySet(doc); err != nil {
		return err
	}

	if err := t.ApplyDefault(doc); err != nil {
		return err
	}

	return nil
}

func (t *Transformation) ApplyDeletes(doc *yaml.Node) error {
	for _, path := range t.Deletes {
		// Ignore errors. We don't care if we didn't find nodes etc.
		_ = Delete(doc, path)
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

func (t *Transformation) ApplyDefault(doc *yaml.Node) error {
	for i := 0; i < len(t.Defaults.Content); i += 2 {
		path, value := t.Defaults.Content[i].Value, t.Defaults.Content[i+1]

		if err := SetDefault(doc, path, value); err != nil {
			log.Printf("going to set default %v=%v", path, value)
			return errors.Wrapf(err, "failed to set default %q", path)
		}
	}

	return nil
}
