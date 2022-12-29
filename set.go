package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func literalNodes(node *yaml.Node) {
	if node.Style == yaml.FlowStyle {
		node.Style = yaml.LiteralStyle
	}

	for _, subNode := range node.Content {
		if subNode.Style == yaml.FlowStyle {
			subNode.Style = yaml.LiteralStyle
		}

		literalNodes(subNode)
	}
}

func Set(doc *yaml.Node, path string, value *yaml.Node) error {
	selectors := strings.Split(path, ".")

	root := getRootNode(doc)
	nodes, err := findNodes(root, selectors, true)
	if err != nil {
		return errors.Wrapf(err, "failed to match %q", path)
	}

	for _, node := range nodes {
		if node.Kind == yaml.ScalarNode {
			// Overwrite an existing scalar value with a new value (whatever kind).
			*node = *value
		} else if node.Kind == yaml.MappingNode && value.Kind == yaml.MappingNode {
			// Append new values onto an existing map node.
			node.Content = append(value.Content, node.Content...)
		} else if node.Kind == yaml.MappingNode && node.Content == nil {
			// Overwrite a new map node we created in findNode(), as confirmed
			// by the nil check (the node.Content wouldn't be nil otherwise).
			*node = *value
		} else if node.Kind == yaml.SequenceNode && value.Kind == yaml.SequenceNode {
			// Append new values onto an existing array node.
			node.Content = append(value.Content, node.Content...)
		} else {
			return errors.Errorf("can't overwrite %v value (line: %v, column: %v) with %v value", node.Tag, node.Line, node.Column, value.Tag)
		}

		// always literal style
		literalNodes(node)
	}

	return nil
}

func SetDefault(doc *yaml.Node, path string, value *yaml.Node) error {
	selectors := strings.Split(path, ".")

	root := getRootNode(doc)
	_, err := findNodes(root, selectors, false)
	if err == nil {
		return nil
	}

	return Set(doc, path, value)
}
