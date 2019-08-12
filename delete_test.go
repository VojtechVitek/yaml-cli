package yaml_test

import (
	"fmt"
	"testing"

	"github.com/VojtechVitek/yaml"
	"github.com/google/go-cmp/cmp"
)

func TestDelete(t *testing.T) {
	tt := []struct {
		delete string
		in     []byte
		out    []byte
	}{
		{
			in:     deployment,
			delete: "kind",
			out:    deploymentWithoutKind,
		},
		{
			in:     deployment,
			delete: "metadata.labels.app",
			out:    deploymentWithoutMetadataLabelsApp,
		},
		{
			in:     []byte(fmt.Sprintf("a:\n  b:\n  c:\n    d: value\nkey: value\n")),
			delete: "a",
			out:    []byte("key: value\n"),
		},
	}

	for i, tc := range tt {
		doc, err := yaml.Parse(tc.in)
		if err != nil {
			t.Error(err)
		}

		if err := yaml.Delete(doc, tc.delete); err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(tc.out, yaml.Bytes(doc)); diff != "" {
			t.Errorf("tc[%v] mismatch (-want +got):\n%s", i, diff)
		}
	}
}

var deployment = []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  labels:
    app: api
    other: label
spec:
  replicas: 3
`)

var deploymentWithoutKind = []byte(`apiVersion: apps/v1
metadata:
  name: api
  labels:
    app: api
    other: label
spec:
  replicas: 3
`)

var deploymentWithoutMetadataLabelsApp = []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  labels:
    other: label
spec:
  replicas: 3
`)
