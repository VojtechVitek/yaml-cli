package yaml

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (t *Transformation) MustMatchAll(doc *yaml.Node) (bool, error) {
	fmt.Fprintf(os.Stderr, "MustMatchAll: %v\n", t.Matches)

	for path, want := range t.Matches {
		selectors := strings.Split(path, ".")
		gotNodes, err := findNodes(doc.Content[0], selectors, false)
		if err != nil {
			return false, errors.Wrapf(err, "failed to match %q", path)
		}

		for _, gotNode := range gotNodes {
			if err := match(&want, gotNode); err != nil {
				return false, errors.Wrapf(err, "failed to match %q (line: %v, column: %v)", path, gotNode.Line, gotNode.Column)
			}
		}
	}

	return true, nil
}

func match(want *yaml.Node, got *yaml.Node) error {
	switch want.Kind {
	case yaml.ScalarNode:
		if got.Value == want.Value {
			return nil
		}
		return errors.Errorf("want %q, got %q", want.Value, got.Value)

	case yaml.SequenceNode:
		for _, wantOneOf := range want.Content {
			if got.Value == wantOneOf.Value {
				return nil
			}
		}
		return errors.Errorf("want one of %v values, got %q", want.Content, got.Value)

	default:
		return errors.Errorf("invalid type; match maches against scalar values or values defined in an array")
	}

	panic("unreachable")
}
