package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Set(doc *yaml.Node, path string, value *yaml.Node) error {
	selectors := strings.Split(path, ".")

	node, err := findNode(doc.Content[0], selectors, true)
	if err != nil {
		return errors.Wrapf(err, "failed to match %q", path)
	}

	if node.Kind == yaml.ScalarNode {
		// Overwrite an existing scalar value with a new value (whatever kind).
		*node = *value
	} else if node.Kind == yaml.MappingNode && value.Kind == yaml.MappingNode {
		// Append new values onto an existing map node.
		node.Content = append(node.Content, value.Content...)
	} else if node.Kind == yaml.MappingNode && node.Content == nil {
		// Overwrite a new map node we created in findNode(), as confirmed
		// by the nil check (the node.Content wouldn't be nil otherwise).
		*node = *value
	} else if node.Kind == yaml.SequenceNode && value.Kind == yaml.SequenceNode {
		// Append new values onto an existing array node.
		node.Content = append(node.Content, value.Content...)
	} else {
		return errors.Errorf("can't overwrite %v value with %v value", node.Tag, value.Tag)
	}

	return nil
}

func SetDefault(doc *yaml.Node, path string, value *yaml.Node) error {
	selectors := strings.Split(path, ".")

	_, err := findNode(doc.Content[0], selectors, false)
	if err == nil {
		return nil
	}

	return Set(doc, path, value)
}
