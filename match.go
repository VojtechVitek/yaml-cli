package yaml

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func (t *Transformation) MustMatchAll(doc *yaml.Node) (bool, error) {
	for path, want := range t.Matches {
		var regexString string

		switch want.Kind {
		case yaml.ScalarNode: // Single value. Ok.
			regexString = want.Value
		case yaml.SequenceNode: // Obsolete syntax, ie. "kind: [Deployment, Pod]". Convert to regex.
			var values []string
			for _, node := range want.Content {
				values = append(values, node.Value)
			}
			regexString = fmt.Sprintf("^%v$", strings.Join(values, "|"))
		default:
			panic(errors.Errorf("Unexpected match kind %v (expected: single value or array)", want.Kind))
		}

		re, err := regexp.Compile(regexString)
		if err != nil {
			return false, errors.Errorf("%q is not a valid regex, see https://github.com/google/re2/wiki/Syntax", regexString)
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
