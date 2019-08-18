package cli_test

import (
	"testing"
)

func TestApply(t *testing.T) {
	tt := []*cliTestCase{
		{
			in:  readFile("_testfiles/apply/in.yml"),
			cmd: []string{"yaml", "apply", "_testfiles/apply/transformations.yml"},
			out: readFile("_testfiles/apply/out.yml"),
		},
	}

	for _, tc := range tt {
		tc.runTest(t)
	}
}
