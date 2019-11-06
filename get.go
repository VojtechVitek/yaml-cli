package yaml

import (
	"gopkg.in/yaml.v3"
)

func Get(doc *yaml.Node, selectors []string) ([]*yaml.Node, error) {
	return findNodes(doc.Content[0], selectors, false)
}
