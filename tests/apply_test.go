package yaml_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/VojtechVitek/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestApply(t *testing.T) {
	tt := []struct {
		in    []byte
		apply *os.File
		out   []byte
	}{
		{
			in:    readFile("apply/in.yml"),
			apply: openFile("apply/apply.yml"),
			out:   readFile("apply/out.yml"),
		},
	}

	for i, tc := range tt {
		doc, err := yaml.Parse(tc.in)
		if err != nil {
			t.Error(err)
		}

		transformations, err := yaml.Transformations(tc.apply)
		if err != nil {
			t.Error(err)
		}

		for _, tf := range transformations {
			if err := tf.Apply(doc); err != nil {
				t.Error(err)
			}
		}

		if diff := cmp.Diff(tc.out, yaml.Bytes(doc)); diff != "" {
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
