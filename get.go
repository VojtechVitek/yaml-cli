package yaml

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func Get(doc *yaml.Node, path string) (*yaml.Node, error) {
	selectors := strings.Split(path, ".")

	node, err := findNode(doc.Content[0], selectors, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to match %q", path)
	}

	return node, nil
}
