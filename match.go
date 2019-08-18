package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (t *Transformation) MustMatchAll(doc *yaml.Node) (bool, error) {
	for path, want := range t.Matches {
		selectors := strings.Split(path, ".")
		got, err := findNode(doc.Content[0], selectors, false)
		if err != nil {
			return false, errors.Wrapf(err, "failed to match %q", path)
		}

		if err := match(&want, got); err != nil {
			return false, errors.Wrapf(err, "failed to match %q", path)
		}
	}
	return true, nil
}

func match(want *yaml.Node, got *yaml.Node) error {
	switch want.Kind {
	case yaml.ScalarNode:
		if got.Value == want.Value {
			return nil
		}
		return errors.Errorf("want %q, got %q", want.Value, got.Value)

	case yaml.SequenceNode:
		for _, wantOneOf := range want.Content {
			if got.Value == wantOneOf.Value {
				return nil
			}
		}
		return errors.Errorf("want one of %v values, got %q", want.Content, got.Value)

	default:
		return errors.Errorf("invalid type; match maches against scalar values or values defined in an array")
	}

	panic("unreachable")
}
