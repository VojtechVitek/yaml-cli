package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (t *Transformation) MustMatchAll(doc *yaml.Node) (bool, error) {
nextMatch:
	for path, want := range t.Matches {
		selectors := strings.Split(path, ".")

		got, err := findMatchingNode(doc.Content[0], selectors)
		if err != nil {
			return false, errors.Wrapf(err, "failed to match %q", path)
		}

		switch want.Kind {
		case yaml.ScalarNode:
			if got.Value == want.Value {
				continue nextMatch
			}
			return false, errors.Errorf("failed to match %q (want %q, got %q)", path, want.Value, got.Value)

		case yaml.SequenceNode:
			for _, wantOneOf := range want.Content {
				if got.Value == wantOneOf.Value {
					continue nextMatch
				}
			}
			return false, errors.Errorf("failed to match %q (want one of %v values, got %q)", path, want.Content, got.Value)

		default:
			return false, errors.Errorf("match %q: invalid type; match can be only scalar value or array", path)
		}
	}
	return true, nil
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
