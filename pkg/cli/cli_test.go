package cli_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/VojtechVitek/yaml/pkg/cli"
	"github.com/google/go-cmp/cmp"
)

type cliTestCase struct {
	in   string
	args []string
	out  string
	err  bool
}

func (tc *cliTestCase) runTest(t *testing.T) {
	t.Helper()

	var b bytes.Buffer

	err := cli.Run(&b, strings.NewReader(tc.in), tc.args)
	if err != nil && !tc.err {
		t.Errorf("%#v: %v", tc.args, err)
	} else if err == nil && tc.err {
		t.Errorf("%#v: expected error", tc.args)
	}

	if diff := cmp.Diff(tc.out, b.String()); diff != "" {
		t.Errorf("%#v mismatch (-want +got):\n%s", tc.args, diff)
	}
}

func readFile(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(b)
}
