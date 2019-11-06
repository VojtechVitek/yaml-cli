package yaml

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (t *Transformation) MustMatchAll(doc *yaml.Node) (bool, error) {
	for path, want := range t.Matches {
		if want.Kind != yaml.ScalarNode {
			return false, errors.Errorf("TODO: Support non-scalar match values?")
		}
		re, err := regexp.Compile(want.Value)
		if err != nil {
			return false, errors.Errorf("%q is not a valid regex, see https://github.com/google/re2/wiki/Syntax", want.Content)
		}

		selectors := strings.Split(path, ".")
		gotNodes, err := findNodes(doc.Content[0], selectors, false)
		if err != nil {
			return false, errors.Wrapf(err, "failed to match %q", path)
		}

		for _, gotNode := range gotNodes {
			if !re.MatchString(gotNode.Value) {
				return false, errors.Wrapf(err, "failed to match %q (line: %v, column: %v)", path, gotNode.Line, gotNode.Column)
			}
		}
	}

	return true, nil
}
