package yaml

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Get(doc *yaml.Node, selectors []string) (*yaml.Node, error) {
	node, err := findNode(doc.Content[0], selectors, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to match %v", selectors)
	}

	return node, nil
}
