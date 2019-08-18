package cli_test

import (
	"testing"
)

func TestSet(t *testing.T) {
	tt := []*cliTestCase{
		{
			in:  ``,
			cmd: []string{"yaml", "set", ""},
			out: ``,
		},
		{
			in:  `first: value`,
			cmd: []string{"yaml", "set", "second: value", "third: value"},
			out: `first: value
second: value
third: value
`,
		},
		{
			in:  `first: value`,
			cmd: []string{"yaml", "set", "first: [several, values]"},
			out: "first: [several, values]\n",
		},
		{
			in:  `first: [a, b]`,
			cmd: []string{"yaml", "set", "first: [several, values]"},
			out: "first: [a, b, several, values]\n",
		},
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
