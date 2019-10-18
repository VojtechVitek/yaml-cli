# YAML CLI processor <!-- omit in toc -->
A CLI tool for querying and transforming YAML data: Grep matching objects, join YAML documents, get/add/edit/delete YAML nodes matching given selector, loop over objects and/or data arrays etc.

`[input.yml] => [query or transformations] => [output.yml]`

*Note: The input YAML data might contain multiple YAML documents separated by `---`.*

- [One-liner commands](#one-liner-commands)
- [yaml get <selector>](#yaml-get-selector)
- [yaml print <selector>](#yaml-print-selector)
  - [yaml set "<selector>: value"](#yaml-set-%22selector-value%22)
  - [yaml default "<selector>: value"](#yaml-default-%22selector-value%22)
  - [yaml delete <selector>](#yaml-delete-selector)
  - [yaml cat file1.yml file2.yml fileN.yml](#yaml-cat-file1yml-file2yml-filenyml)
  - [yaml count <selector>](#yaml-count-selector)
- [Transformation files](#transformation-files)
  - [yaml apply file1.yt file2.yt fileN.yt](#yaml-apply-file1yt-file2yt-filenyt)
    - [Examples of transformation YAML files](#examples-of-transformation-yaml-files)
- [Print selected YAML nodes](#print-selected-yaml-nodes)
  - [yaml get <selector>](#yaml-get-selector-1)
    - [kubectl - print pod's main container image](#kubectl---print-pods-main-container-image)
- [Grep documents/objects](#grep-documentsobjects)
  - [yaml grep "<selector>: value" ...](#yaml-grep-%22selector-value%22)
    - [Grep k8s deployment object by name](#grep-k8s-deployment-object-by-name)
  - [yaml grep -v "<selector>: value" ...](#yaml-grep--v-%22selector-value%22)
    - [Grep all k8s objects that don't create Pods](#grep-all-k8s-objects-that-dont-create-pods)
    - [Print first container's image of linkerd2 deployment objects](#print-first-containers-image-of-linkerd2-deployment-objects)
- [Useful Kubernetes examples](#useful-kubernetes-examples)
    - [Push all non-pod objects to k8s](#push-all-non-pod-objects-to-k8s)
    - [Rollout k8s deployments from desired-state files sequentially](#rollout-k8s-deployments-from-desired-state-files-sequentially)
- [Known issues](#known-issues)
- [Feedback](#feedback)
- [License](#license)

# One-liner commands

# yaml get <selector>
Get value of a node matching the given selector.
```bash
$ kubectl get pod/nats-8576dfb67-vg6v7 -o yaml | yaml get spec.containers[0].image
nats-streaming:0.10.0
```

# yaml print <selector>
Print full node matching the given selector.
```bash
$ kubectl get pod/nats-8576dfb67-vg6v7 -o yaml | yaml print kind apiVersion metadata.annotations
kind: Pod
apiVersion: v1
metadata.annotations:
  kubernetes.io/psp: eks.privileged
```

## yaml set "<selector>: value"
Add/overwrite field's value.
```bash
$ cat input.yml | yaml set "metadata.labels.environment: staging" > output.yml
```

## yaml default "<selector>: value"
Set field's value, if no such value exists yet.
```bash
$ cat input.yml | yaml default "metadata.labels.environment: staging" > output.yml
```

## yaml delete <selector>
Delete specific field.
```bash
$ cat input.yml | yaml delete "metadata.labels.environment" > output.yml
```

## yaml cat file1.yml file2.yml fileN.yml
Join multiple YAML files into a single file with multiple documents separated by `---`.
```bash
$ yaml cat k8s-apps/*.yml > output.yml
```

## yaml count <selector>
Print number of items within an array matching the given selector. Useful for ranging over arrays.
```bash
pods=$(kubectl get pods -o yaml)
count=$(echo "$pods" | yaml count items)
for ((i=0; i < $count; i++)); do
    echo "$pods" | yaml get items[$i].status.phase
done
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

# Print selected YAML nodes
## yaml get <selector>
### kubectl - print pod's main container image
```bash
$ kubectl get pods/nats-8576dfb67-vg6v7 -o yaml | yaml get spec.containers[0].image
nats-streaming:0.10.0
```

# Grep documents/objects
Grep documents/objects matching all of the given `selector: value` pairs.

If a provided value is an array (ie. `selector: [first, second]`), the matching value must match at least one of the provided values (logical "OR").

## yaml grep "<selector>: value" ...
### Grep k8s deployment object by name
```bash
$ cat desired-state.yml | yaml grep "kind: Deployment" "metadata.name: linkerd"
```

## yaml grep -v "<selector>: value" ...
### Grep all k8s objects that don't create Pods
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

# Feedback
Any feedback welcome! Please open issues and feature requests..

# License
Licensed under the [MIT License](./LICENSE).
