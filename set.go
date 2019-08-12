package yaml

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Set(doc *yaml.Node, path string, value *yaml.Node) error {
	selectors := strings.Split(path, ".")

	node, err := findOrCreateNode(doc.Content[0], selectors)
	if err != nil {
		return errors.Wrapf(err, "failed to match %q", path)
	}

	if node.Kind == yaml.MappingNode && value.Kind == yaml.MappingNode {
		// Append new values onto an existing map node.
		node.Content = append(node.Content, value.Content...)
	} else if node.Kind == yaml.MappingNode && node.Content == nil {
		// Overwrite a new map node we created in findOrCreateNode(), as confirmed
		// by the nil check (the node.Content wouldn't be nil otherwise).
		*node = *value
	} else if node.Kind == yaml.ScalarNode && value.Kind == yaml.ScalarNode {
		// Overwrite an existing scalar value with a new value.
		*node = *value
	} else if node.Kind == yaml.SequenceNode && value.Kind == yaml.SequenceNode {
		// Append new values onto an existing array node.
		node.Content = append(node.Content, value.Content...)
	} else {
		return errors.Errorf("can't overwrite %v value with %v value", node.Tag, value.Tag)
	}

	return nil
}

func findOrCreateNode(node *yaml.Node, selectors []string) (*yaml.Node, error) {
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
		node, err = findOrCreateNode(node, []string{currentSelector})
		if err != nil {
			return nil, errors.Errorf("can't find %v", currentSelector)
		}

		if node.Kind != yaml.SequenceNode {
			return nil, errors.Errorf("%v is not an array", currentSelector)
		}

		if index > len(node.Content) {
			return nil, errors.Errorf("%v array doesn't have index %v", currentSelector, index)
		}

		return findOrCreateNode(node.Content[index], selectors[1:])
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
					return findOrCreateNode(node.Content[i+1], selectors[1:])
				}

				// Found last key, return its value.
				return node.Content[i+1], nil
			}
		}

	case yaml.ScalarNode:
		// Overwrite any existing nodes.
		node.Kind = yaml.MappingNode
		node.Tag = "!!map"
		node.Value = ""

	default:
		return nil, errors.Errorf("unknown node.Kind %v", node.Kind)
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
