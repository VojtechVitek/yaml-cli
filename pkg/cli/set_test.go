package cli_test

import (
	"testing"
)

func TestSet(t *testing.T) {
	tt := []*cliTestCase{
		{
			in:   "",
			args: []string{"set", ""},
			out:  "",
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
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
