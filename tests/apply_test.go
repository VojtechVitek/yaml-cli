package yaml_test

import (
	"io/ioutil"
	"testing"

	"github.com/VojtechVitek/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestApply(t *testing.T) {
	tt := []struct {
		in    []byte
		apply []byte
		out   []byte
	}{
		{
			in:    readFile("apply/in.yml"),
			apply: readFile("apply/apply.yml"),
			out:   readFile("apply/out.yml"),
		},
	}

	for i, tc := range tt {
		doc, err := yaml.Parse(tc.in)
		if err != nil {
			t.Error(err)
		}

		transformation, err := yaml.ParseTransformation(tc.apply)
		if err != nil {
			t.Error(err)
		}

		if err := transformation.Apply(doc); err != nil {
			t.Error(err)
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
