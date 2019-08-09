package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Set(doc *yaml.Node, path string, value *yaml.Node) error {
	selectors := strings.Split(path, ".")

	node, err := findOrCreateMatchingNode(doc.Content[0], selectors)
	if err != nil {
		return errors.Wrapf(err, "failed to match %q", path)
	}

	*node = *value

	return nil
}

func findOrCreateMatchingNode(node *yaml.Node, selectors []string) (*yaml.Node, error) {
	currentSelector := selectors[0]
	lastSelector := len(selectors) == 1

	// Iterate over the keys (the slice is key/value pairs).
	switch node.Kind {
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			// Does current key match the selector?
			if node.Content[i].Value == currentSelector {
				if !lastSelector {
					// Match the rest of the selector path, ie. go deeper
					// in to the value node.
					return findOrCreateMatchingNode(node.Content[i+1], selectors[1:])
				}

				// Found last key, return its value.
				return node.Content[i+1], nil
			}
		}
	case yaml.ScalarNode:
		// Override existing nodes.
		node.Kind = yaml.MappingNode
		node.Tag = "!!map"
		node.Value = ""
	}

	// Create the rest of the selector path.
	for _, selector := range selectors {
		node.Content = append(node.Content,
			&yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: selector,
			},
			&yaml.Node{
				Kind: yaml.MappingNode,
				Tag:  "!!map",
			},
		)

		node = node.Content[len(node.Content)-1]
	}

	return node, nil
}
