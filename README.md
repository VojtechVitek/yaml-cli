# Streaming YAML CLI processor <!-- omit in toc -->
A CLI tool for querying and transforming YAML stream data:
- Grep matching documents (ie. K8s objects)
- Join multiple YAML files
- Get/add/edit/delete YAML nodes matching given selector
- Loop over documents and/or data arrays
- etc.

`[input.yml] => [query or transformations] => [output.yml]`

*Note: The input YAML documents in a [YAML stream](https://yaml.org/spec/1.2/spec.html#id2801681) are separated by `---`.*

- [One-liner commands](#one-liner-commands)
  - [yaml get selector](#yaml-get-selector)
    - [yaml get selector --print-key](#yaml-get-selector---print-key)
    - [yaml get array[*] --print-key](#yaml-get-array---print-key)
    - [yaml get spec.containers[*].image --no-separator](#yaml-get-speccontainersimage---no-separator)
  - [yaml set "selector: value"](#yaml-set-%22selector-value%22)
  - [yaml default "selector: value"](#yaml-default-%22selector-value%22)
  - [yaml delete selector](#yaml-delete-selector)
  - [yaml cat file1.yml file2.yml fileN.yml](#yaml-cat-file1yml-file2yml-filenyml)
  - [yaml count](#yaml-count)
  - [yaml doc $index](#yaml-doc-index)
  - [yaml len arraySelector](#yaml-len-arrayselector)
- [Grep YAML documents (objects)](#grep-yaml-documents-objects)
  - [yaml grep "selector: value" ...](#yaml-grep-%22selector-value%22)
    - [Grep k8s deployment object by name](#grep-k8s-deployment-object-by-name)
  - [yaml grep -v "selector: value" ...](#yaml-grep--v-%22selector-value%22)
    - [Grep all k8s objects that don't create any Pods](#grep-all-k8s-objects-that-dont-create-any-pods)
    - [Print first container's image of linkerd2 deployment objects](#print-first-containers-image-of-linkerd2-deployment-objects)
- [Useful Kubernetes examples](#useful-kubernetes-examples)
    - [Push all non-pod objects to k8s](#push-all-non-pod-objects-to-k8s)
    - [Rollout k8s deployments from desired-state files sequentially](#rollout-k8s-deployments-from-desired-state-files-sequentially)
- [Known issues](#known-issues)
- [Working with JSON](#working-with-json)
  - [JSON to YAML](#json-to-yaml)
  - [YAML to JSON](#yaml-to-json)
- [Transformation files](#transformation-files)
  - [yaml apply file1.yt file2.yt fileN.yt](#yaml-apply-file1yt-file2yt-filenyt)
    - [Examples of transformation YAML files](#examples-of-transformation-yaml-files)
- [Feedback](#feedback)
- [License](#license)

# One-liner commands

## yaml get selector
Print value of a YAML node matching the given selector.
```bash
$ kubectl get pod/nats-8576dfb67-vg6v7 -o yaml | yaml get spec.containers[0].image
nats-streaming:0.10.0
```

### yaml get selector --print-key

Since we're printing a value, the output might not necesarilly be a valid YAML.

If we're printing value of a primitive type (ie. string) and we need the output in a valid YAML format, so it can be processed further, we can explicitly print the node key in front of the value:

```bash
$ kubectl get pod/nats-8576dfb67-vg6v7 -o yaml | yaml get spec.containers[0].image --print-key
image: nats-streaming:0.10.0
```

We can print multiple values, and they will be printed as separate objects:

```bash
$ kubectl get pod/nats-8576dfb67-vg6v7 -o yaml | yaml get spec.containers[0].image spec.containers[1].image --print-key
image: nats-streaming:0.10.0
---
image: sidecar:1.0.1
```

### yaml get array[*] --print-key

We can print all array items at once with a wildcard (`array[*]`) too:

```bash
$ kubectl get pod/nats-8576dfb67-vg6v7 -o yaml | yaml get spec.containers[*].image --print-key
image: nats-streaming:0.10.0
---
image: sidecar:1.0.1
```

### yaml get spec.containers[*].image --no-separator

Need to print values only?

```bash
$ kubectl get pod/nats-8576dfb67-vg6v7 -o yaml | yaml get spec.containers[*].image --no-separator
nats-streaming:0.10.0
sidecar:1.0.1
```

## yaml set "selector: value"
Add/overwrite field's value.
```bash
$ cat input.yml | yaml set "metadata.labels.environment: staging" > output.yml
```

## yaml default "selector: value"
Set field's value, if no such value exists yet.
```bash
$ cat input.yml | yaml default "metadata.labels.environment: staging" > output.yml
```

## yaml delete selector
Delete specific field.
```bash
$ cat input.yml | yaml delete "metadata.labels.environment" > output.yml
```

## yaml cat file1.yml file2.yml fileN.yml
Join multiple YAML files into a single file with multiple documents separated by `---`.
```bash
$ yaml cat k8s-apps/*.yml > output.yml
```

## yaml count
Print number of YAML documents within the input YAML stream.

input.yml
```yml
document: this is doc 1
---
document: this is doc 2
```

```bash
$ yaml count input.yml
2
```

## yaml doc $index
Print Nth (index=0..N-1) YAML document from the input YAML stream.

input.yml
```yml
document: this is doc 1
---
document: this is doc 2
```

```bash
$ yaml doc 1
document: this is doc 2
```

## yaml len arraySelector
Print number of items within an array matching the given selector.

Useful for ranging over arrays.
```bash
pods=$(kubectl get pods -o yaml)
count=$(echo "$pods" | yaml len items)
for ((i=0; i < $count; i++)); do
    echo "$pods" | yaml get items[$i].status.phase
done
```

# Grep YAML documents (objects)
Grep documents/objects matching all of the given `selector: value` pairs.

If a provided value is an array (ie. `selector: [first, second]`), the matching value must match at least one of the provided values (logical "OR").

## yaml grep "selector: value" ...
### Grep k8s deployment object by name
```bash
$ cat desired-state.yml | yaml grep "kind: Deployment" "metadata.name: linkerd"
```

## yaml grep -v "selector: value" ...
Inverse grep.

### Grep all k8s objects that don't create any Pods
```bash
$ cat desired-state.yml | yaml grep -v "kind: [Deployment, Pod, Job, ReplicaSet, ReplicationController]"
```

### Print first container's image of linkerd2 deployment objects
```bash
$ cat linkerd.yml | yaml grep "kind: Deployment" | yaml get "spec.template.spec.containers[0].image"
gcr.io/linkerd-io/controller:stable-2.4.0
gcr.io/linkerd-io/controller:stable-2.4.0
gcr.io/linkerd-io/web:stable-2.4.0
prom/prometheus:v2.10.0
gcr.io/linkerd-io/grafana:stable-2.4.0
gcr.io/linkerd-io/controller:stable-2.4.0
gcr.io/linkerd-io/controller:stable-2.4.0
gcr.io/linkerd-io/controller:stable-2.4.0
```

# Useful Kubernetes examples
### Push all non-pod objects to k8s
```bash
$ cat desired-state.yml | yaml grep -v "kind: [Deployment, Pod, Job]"
```

### Rollout k8s deployments from desired-state files sequentially
```bash
$ for file in *.yml; do
    out=$(cat $file | yaml grep "kind: Deployment" | kubectl apply -f -)

    for deploy in $(echo "$out" | cut -d' ' -f1); do
        kubectl rollout status --timeout 180s $deploy || {
            kubectl rollout undo $deploy
            exit 1
        }
    done
  done
```

# Known issues

1. Merging complex nodes doesn't work well
```yaml
set:
    metadata:
        we:
            cant:
                merge:
                    complex: objects
                    such:
                        as: this
                        properly:
                            just: yet
```
We'll want to fix this later. For now, use explicit paths to the final nodes:
```yaml
set:
    metadata.we.cant.merge.complex: objects
    metadata.we.cant.merge.such.as: this
    metadata.we.cant.merge.such.properly.just: yet
```

2. Wildcard array[*] matching doesn't work yet
```yaml
match:
    spec.template.spec.containers[*].name: prometheus
```

3. Selectors with `.` dots in the selector path, ie.
```yml
metadata:
    annotations:
        linkerd.io/inject: enabled
```
Since the `match` selectors are separated with `.` dots, we'll have to figure out how to support these selector keys with inner `.` dots.

We might wanna support
```yml
delete:
    - metadata.annotations."linkerd.io/inject"
set:
    metadata.annotations."rbac.authorization.kubernetes.io/autoupdate": true
```

# Working with JSON

## JSON to YAML
```
$ cat file.json | yaml --from=json grep 'kind: Pod'
```

## YAML to JSON
```
$ cat file.yml | yaml grep 'kind: Pod' --to=json
```

# Transformation files

All of the above examples, and more, can be described in YAML transformation file syntax. Multiple such transformations can be applied at once.

## yaml apply file1.yt file2.yt fileN.yt
Apply multiple YAML "transformations", see the `.yt` file syntax below.
```bash
$ yaml cat k8s-apps/*.yml | yaml apply staging.yt enable-linkerd.yt > staging/desired-state.yml
```

staging.yt:
```yml
match:
    # all YAML objects
set:
    metadata.labels.environment: staging
---
match:
    kind: Deployment
    metadata.name: api
set:
    metadata.labels.first: updated-label
    spec.replicas: 3
```

enable-linkerd.yt:
```yml
match:
  kind: [Deployment, Pod]
default:
  metadata.annotations:
    linkerd.io/inject: enabled
```

Changes applied to the original object:
```diff
 apiVersion: apps/v1
 kind: Deployment
 metadata:
     name: api
     labels:
-        first: label
+        first: updated-label
         second: label
+        environment: staging
+    annotations:
+        linkerd.io/inject: enabled
 spec:
-    replicas: 1
+    replicas: 3
 ...
```

### Examples of transformation YAML files
```yml
match:
    kind: Deployment
set:
    spec.template.spec.containers[*]:
        imagePullPolicy: IfNotPresent
```
```yml
match:
    kind: Deployment
set:
    spec.template.spec.nodeSelector:
        worker: generic
```
```yml
match:
    kind: Deployment
    metadata.name: api
set:
    spec.replicas: 3
```
```yml
match:
    kind: Ingress
set:
    metadata:
        annotations:
            kubernetes.io/ingress.class: nginx
```
```yml
match:
    kind: Deployment
delete: spec.replicas
```

```yml
match:
    kind: Deployment
    spec.template.spec.containers[0].image: nats-streaming
set:
    spec.template.spec.containers[0].image: nats-streaming:0.15.1
```

# Feedback
Any feedback welcome! Please open issues and feature requests..

# License
Licensed under the [MIT License](./LICENSE).
