package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// type Transformation struct {
// 	match  map[string]string `yaml:"match"`
// 	delete string            `yaml:"delete"`
// }

func Delete(input []byte, path string) ([]byte, error) {
	var doc yaml.Node

	err := yaml.Unmarshal(input, &doc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal input")
	}

	if len(doc.Content) != 1 {
		// Looks like any YAML doc has always 1 node..
		return nil, errors.Errorf("len(doc.Content)=%v", len(doc.Content))
	}

	selectors := strings.Split(path, ".")

	err = deleteMatchingNode(doc.Content[0], selectors)
	if err != nil {
		return input, errors.Wrapf(err, "failed to match %q", path)
	}

	return yaml.Marshal(&doc)
}

func deleteMatchingNode(node *yaml.Node, selectors []string) error {
	currentSelector := selectors[0]
	lastSelector := len(selectors) == 1

	// Iterate over the keys (the slice is key/value pairs).
	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == currentSelector {
			// Key matches the selector.
			if !lastSelector {
				// Try to match the rest of the selector path in the value.
				return deleteMatchingNode(node.Content[i+1], selectors[1:])
			}

			node.Content[i] = nil   // Delete key.
			node.Content[i+1] = nil // Delete value.
			node.Content = append(node.Content[:i], node.Content[i+2:]...)
			return nil
		}
	}

	return errors.Errorf("can't find %q", strings.Join(selectors, "."))
}
