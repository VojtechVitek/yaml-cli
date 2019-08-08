package transform

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTransform(t *testing.T) {
	input := []byte(`
match:
    kind: Deployment
delete: spec.replicas
`)

	{
		var transformer *Transformer
		err := yaml.Unmarshal(input, &transformer)
		if err != nil {
			t.Fatal(err)
		}

		result, err := transformer.Delete(deployment, "kind")
		if err != nil {
			t.Error(err)
		}

		t.Logf("\n%s", result)
	}

	{
		var transformer *Transformer
		err := yaml.Unmarshal(input, &transformer)
		if err != nil {
			t.Fatal(err)
		}

		result, err := transformer.Delete(deployment, "metadata.labels.app")
		if err != nil {
			t.Error(err)
		}

		t.Logf("\n%s", result)
	}
}

/*
match:
    kind: Deployment
    metadata.name: api

add:
    spec.replicas: 3

match:
    kind: Deployment
    metadata.name: service
set:
    spec.replicas: 3

	match:
    kind: Deployment
delete: spec.replicas
*/

var deployment = []byte(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  labels:
    app: api
    some: other-label
spec:
  replicas: 3
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
  labels:
    app: api
spec:
  replicas: 3
`)
