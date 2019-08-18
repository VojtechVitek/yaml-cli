package yaml

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func findNode(node *yaml.Node, selectors []string, create bool) (*yaml.Node, error) {
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
		node, err = findNode(node, []string{currentSelector}, create)
		if err != nil {
			return nil, errors.Errorf("can't find %v", currentSelector)
		}

		if node.Kind != yaml.SequenceNode {
			return nil, errors.Errorf("%v is not an array", currentSelector)
		}

		if index >= len(node.Content) {
			return nil, errors.Errorf("%v array doesn't have index %v", currentSelector, index)
		}

		if len(selectors) == 1 { // Last selector.
			return node.Content[index], nil
		}
		return findNode(node.Content[index], selectors[1:], create)
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
					return findNode(node.Content[i+1], selectors[1:], create)
				}

				// Found last key, return its value.
				return node.Content[i+1], nil
			}
		}

	case yaml.ScalarNode:
		if create {
			// Overwrite any existing nodes.
			node.Kind = yaml.MappingNode
			node.Tag = "!!map"
			node.Value = ""
		}

	default:
		return nil, errors.Errorf("unknown node.Kind %v", node.Kind)
	}

	if create {
		if node.Kind == yaml.ScalarNode {

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
	}

	return node, nil
}
