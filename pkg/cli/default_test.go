package cli_test

import (
	"testing"
)

func TestDefault(t *testing.T) {
	tt := []*cliTestCase{
		{
			in:   "",
			args: []string{"default", ""},
			out:  ``,
		},
		{
			in:   "first: value",
			args: []string{"default", "first: updated"},
			out:  "first: value\n",
		},
		{
			in:   "first: value",
			args: []string{"default", "second: value", "third: value"},
			out:  "first: value\nsecond: value\nthird: value\n",
		},
		{
			in:   "first: value",
			args: []string{"default", "first: [several, values]"},
			out:  "first: value\n",
		},
		{
			in:   "first: [a, b]",
			args: []string{"default", "first: updated"},
			out:  "first: [a, b]\n",
		},
		{
			in:   "first: [a, b]",
			args: []string{"default", "second: [a, b]"},
			out:  "first: [a, b]\nsecond: [a, b]\n",
		},
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
