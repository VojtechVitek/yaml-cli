package cli_test

import (
	"fmt"
	"testing"
)

func TestGrep(t *testing.T) {
	var fooYAML = `---
foo: bar
---
foo: baz
---
# empty
---
no-foo: nope
---
foo: nope
---
foo:
  - bar
---
foo:
  - bar
  - baz
---
foo:
  bar: baz
`

	tt := []*cliTestCase{
		{
			in:   fooYAML,
			args: []string{"grep", "foo: bar"},
			out:  "foo: bar\n",
		},
		{
			in:   fooYAML,
			args: []string{"grep", "foo: [bar, baz]"},
			out:  "foo: bar\n---\nfoo: baz\n",
		},
		{
			in:   fooYAML,
			args: []string{"grep", "doesnt: exist"},
			out:  "",
		},
		{
			in:   fmt.Sprintf("foo: bar\n---\nfoo: baz\n---\nfoo: x"),
			args: []string{"grep", "foo: [bar, baz]"},
			out:  fmt.Sprintf("foo: bar\n---\nfoo: baz\n"),
		},
		{
			in:   fmt.Sprintf("foo: bar\n---\nfoo: baz\n---\nfoo: x"),
			args: []string{"grep", "-v", "foo: [bar, baz]"},
			out:  fmt.Sprintf("foo: x\n"),
		},
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
