package yaml

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Get(doc *yaml.Node, path string) (*yaml.Node, error) {
	selectors := strings.Split(path, ".")

	node, err := findOrCreateNode(doc.Content[0], selectors)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to match %q", path)
	}

	return node, nil
}

func findNode(node *yaml.Node, selectors []string) (*yaml.Node, error) {
	currentSelector := selectors[0]

	// array[N] selectors.
	if i := strings.LastIndex(currentSelector, "["); i > 0 && strings.HasSuffix(currentSelector, "]") {
		arrayIndex := currentSelector[i+1 : len(currentSelector)-1]
		currentSelector = currentSelector[:i]

		// TODO: Can we do array[*], ie. set value in each array item?
		index, err := strconv.Atoi(arrayIndex)
		if err != nil {
			return nil, errors.Wrapf(err, "can't parse array index from %v[%v]", currentSelector, arrayIndex)
		}

		// Go into an array.
		node, err = findNode(node, []string{currentSelector})
		if err != nil {
			return nil, errors.Errorf("can't find %v", currentSelector)
		}

		if node.Kind != yaml.SequenceNode {
			return nil, errors.Errorf("%v is not an array", currentSelector)
		}

		if index > len(node.Content) {
			return nil, errors.Errorf("%v array doesn't have index %v", currentSelector, index)
		}

		return findNode(node.Content[index], selectors[1:])
	}

	// Iterate over the keys (the slice is key/value pairs).
	switch node.Kind {
	case yaml.MappingNode:
		lastSelector := len(selectors) == 1

		for i := 0; i < len(node.Content); i += 2 {
			// Does current key match the selector?
			if node.Content[i].Value == currentSelector {
				if !lastSelector {
					// Match the rest of the selector path, ie. go deeper
					// in to the value node.
					return findNode(node.Content[i+1], selectors[1:])
				}

				// Found last key, return its value.
				return node.Content[i+1], nil
			}
		}

	case yaml.ScalarNode:
		return node, nil

	default:
		return nil, errors.Errorf("unknown node.Kind %v", node.Kind)
	}

	panic("unreachable")
}
