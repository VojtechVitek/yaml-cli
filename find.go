package yaml

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func findNodes(node *yaml.Node, selectors []string, create bool) ([]*yaml.Node, error) {
	var nodes []*yaml.Node
	currentSelector := selectors[0]

	// array[N] or array[*] selectors.
	if i := strings.LastIndex(currentSelector, "["); i > 0 && strings.HasSuffix(currentSelector, "]") {
		arrayIndex := currentSelector[i+1 : len(currentSelector)-1]
		currentSelector = currentSelector[:i]

		index, err := strconv.Atoi(arrayIndex)
		if err != nil {
			if arrayIndex == "*" {
				index = -1
			} else {
				return nil, errors.Wrapf(err, "can't parse array index from %v[%v]", currentSelector, arrayIndex)
			}
		} else if index < 0 {
			return nil, errors.Wrapf(err, "array index can't be negative %v[%v]", currentSelector, arrayIndex)
		}

		// Go into array node(s).
		arrayNodes, err := findNodes(node, []string{currentSelector}, create)
		if err != nil {
			return nil, errors.Errorf("can't find %v", currentSelector)
		}
		for _, arrayNode := range arrayNodes {
			if arrayNode.Kind != yaml.SequenceNode {
				return nil, errors.Errorf("%v is not an array", currentSelector)
			}
			if index >= len(arrayNode.Content) {
				return nil, errors.Errorf("%v array doesn't have index %v", currentSelector, index)
			}

			var visitArrayNodes []*yaml.Node
			if index >= 0 { // array[N]
				visitArrayNodes = []*yaml.Node{arrayNode.Content[index]}
			} else { // array[*]
				visitArrayNodes = arrayNode.Content
			}

			for i, node := range visitArrayNodes {
				if len(selectors) == 1 {
					// Last selector, use this as final node.
					nodes = append(nodes, node)
				} else {
					// Go deeper into a specific array.
					deeperNodes, err := findNodes(node, selectors[1:], create)
					if err != nil {
						return nil, errors.Wrapf(err, "failed to go deeper into %v[%v]", currentSelector, i)
					}
					nodes = append(nodes, deeperNodes...)
				}
			}
		}
		return nodes, nil
	}

	// Iterate over the keys (the slice is key/value pairs).
	switch node.Kind {
	case yaml.MappingNode:
		lastSelector := len(selectors) == 1

		for i := 0; i < len(node.Content); i += 2 {
			// Does the current key match the selector?
			if node.Content[i].Value == currentSelector {
				if !lastSelector {
					// Match the rest of the selector path, ie. go deeper
					// in to the value node.
					return findNodes(node.Content[i+1], selectors[1:], create)
				}

				// Found last key, return its value.
				return []*yaml.Node{node.Content[i+1]}, nil
			}
		}

		if !create {
			return nil, errors.Errorf("can't find node %q", currentSelector)
		}

	case yaml.ScalarNode:
		if create {
			// Overwrite any existing nodes.
			node.Kind = yaml.MappingNode
			node.Tag = "!!map"
			node.Value = ""
		}

	case yaml.SequenceNode:
		return nil, errors.Errorf("parent node is array, use [*] or [0]..[%v] instead of .%v to access its item(s) first", len(node.Content)-1, currentSelector)

	default:
		return nil, errors.Errorf("parent node is of unknown kind %v", node.Kind)
	}

	if create {
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

	return []*yaml.Node{node}, nil
}
