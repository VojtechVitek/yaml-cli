package yaml_test

import (
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
	}

	for i, tc := range tt {
		got, err := yaml.Delete(tc.in, tc.delete)
		if err != nil {
			t.Error(err)
		}

		if diff := cmp.Diff(tc.out, got); diff != "" {
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
