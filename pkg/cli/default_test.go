package cli_test

import (
	"testing"
)

func TestDefault(t *testing.T) {
	tt := []*cliTestCase{
		{
			in:  ``,
			cmd: []string{"yaml", "default", ""},
			out: ``,
		},
		{
			in:  `first: value`,
			cmd: []string{"yaml", "default", "first: updated"},
			out: "first: value\n",
		},
		{
			in:  `first: value`,
			cmd: []string{"yaml", "default", "second: value", "third: value"},
			out: "first: value\nsecond: value\nthird: value\n",
		},
		{
			in:  `first: value`,
			cmd: []string{"yaml", "default", "first: [several, values]"},
			out: "first: value\n",
		},
		{
			in:  `first: [a, b]`,
			cmd: []string{"yaml", "default", "first: updated"},
			out: "first: [a, b]\n",
		},
		{
			in:  `first: [a, b]`,
			cmd: []string{"yaml", "default", "second: [a, b]"},
			out: "first: [a, b]\nsecond: [a, b]\n",
		},
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
