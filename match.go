package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (t *Transformation) MustMatchAll(doc *yaml.Node) (bool, error) {
	for path, want := range t.Matches {
		selectors := strings.Split(path, ".")
		got, err := findMatchingNode(doc.Content[0], selectors)
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

func findMatchingNode(node *yaml.Node, selectors []string) (*yaml.Node, error) {
	currentSelector := selectors[0]
	lastSelector := len(selectors) == 1

	// Iterate over the keys (the slice is key/value pairs).
	for i := 0; i < len(node.Content); i += 2 {
		// Does current key match the selector?
		if node.Content[i].Value == currentSelector {
			if !lastSelector {
				// Try to match the rest of the selector path in the value.
				return findMatchingNode(node.Content[i+1], selectors[1:])
			}

			// Found last key, return its value.
			return node.Content[i+1], nil
		}
	}

	return nil, errors.Errorf("can't find %q", strings.Join(selectors, "."))
}
