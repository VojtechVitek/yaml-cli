package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (t *Transformation) MatchesAll(doc *yaml.Node, conds map[string]string) (bool, error) {
	for path, value := range conds {
		selectors := strings.Split(path, ".")

		match, err := findMatchingNode(doc.Content[0], selectors)
		if err != nil {
			return false, errors.Wrapf(err, "failed to match %q", path)
		}

		if match.Value != value {
			return false, errors.Errorf("failed to match %q (want %q, got %q)", path, value, match.Value)
		}
	}

	return true, nil
}

func findMatchingNode(node *yaml.Node, selectors []string) (*yaml.Node, error) {
	currentSelector := selectors[0]
	lastSelector := len(selectors) == 1

	// Iterate over the keys (the slice is key/value pairs).
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == currentSelector {
			// Key matches the selector.
			if !lastSelector {
				// Try to match the rest of the selector path in the value.
				return findMatchingNode(node.Content[i+1], selectors[1:])
			}

			// Return value
			return node.Content[i+1], nil
		}
	}

	return nil, errors.Errorf("can't find %q", strings.Join(selectors, "."))
}
