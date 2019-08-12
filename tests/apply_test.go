package yaml_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/VojtechVitek/yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	yamlv3 "gopkg.in/yaml.v3"
)

func TestApply(t *testing.T) {
	tt := []struct {
		in  *os.File
		tf  *os.File
		out []byte
	}{
		{
			in:  openFile("apply/in.yml"),
			tf:  openFile("apply/transformations.yml"),
			out: readFile("apply/out.yml"),
		},
	}

	for i, tc := range tt {
		transformations, err := yaml.Transformations(tc.tf)
		if err != nil {
			t.Error(err)
		}

		dec := yamlv3.NewDecoder(tc.in)

		var b bytes.Buffer
		enc := yamlv3.NewEncoder(&b)
		enc.SetIndent(2)

		for { // For all YAML documents
			var doc yamlv3.Node
			if err := dec.Decode(&doc); err != nil {
				if err == io.EOF { // Last document.
					break
				}
				t.Fatal(errors.Wrap(err, "failed to decode YAML document(s)"))
			}

			for _, tf := range transformations {
				t.Logf("tf[%v]: %+v", tf.Matches, doc.Content)
				if err := tf.Apply(&doc); err != nil {
					t.Error(err)
				}
			}

			_ = enc.Encode(&doc)
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
