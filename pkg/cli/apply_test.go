package cli_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/VojtechVitek/yaml/pkg/cli"
	"github.com/google/go-cmp/cmp"
)

func TestApply(t *testing.T) {
	tt := []struct {
		in  *os.File
		cmd []string
		out []byte
	}{
		{
			in:  openFile("../../tests/apply/in.yml"),
			cmd: []string{"yaml", "apply", "../../tests/apply/transformations.yml"},
			out: readFile("../../tests/apply/out.yml"),
		},
	}

	for i, tc := range tt {
		var b bytes.Buffer

		if err := cli.Run(&b, tc.in, tc.cmd); err != nil {
			t.Errorf("tc[%v]: %v", i, err)
		}

		if diff := cmp.Diff(tc.out, b.Bytes()); diff != "" {
			t.Errorf("tc[%v] mismatch (-want +got):\n%s", i, diff)
		}
	}
}

func readFile(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return b
}

func openFile(filename string) *os.File {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	return f
}
