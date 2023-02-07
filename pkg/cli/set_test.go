package cli_test

import (
	"testing"
)

func TestSet(t *testing.T) {
	tt := []*cliTestCase{
		{
			in:   "",
			args: []string{"set", ""},
			out:  "{}\n",
		},
		{
			in:   "",
			args: []string{"set", "first: value"},
			out:  "first: value\n",
		},
		{
			in:   "first: value",
			args: []string{"set", "second: value", "third: value"},
			out: `first: value
second: value
third: value
`,
		},
		{
			in:   "first: value",
			args: []string{"set", "first: [several, values]"},
			out:  "first: [several, values]\n",
		},
		{
			in:   "first: [a, b]",
			args: []string{"set", "first: [several, values]"},
			out:  "first: [a, b, several, values]\n",
		},
		{
			in:   "metadata:\n  overwrite: me\n",
			args: []string{"set", "metadata.overwrite: with", "metadata.overwrite: value"},
			out:  "metadata:\n  overwrite: value\n",
		},
		{
			in:   "first: object\n---\nsecond: object\n",
			args: []string{"set", "foo: bar"},
			out:  "first: object\nfoo: bar\n---\nsecond: object\nfoo: bar\n",
		},
		{
			in:   "first: object\n---\nsecond.foo: object\n",
			args: []string{"set", "second\\.foo: bar"},
			out:  "first: object\nsecond.foo: bar\n---\nsecond.foo: bar\n",
		},
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
