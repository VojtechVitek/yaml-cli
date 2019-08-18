# yaml CLI <!-- omit in toc -->
A CLI tool for transforming YAML files: Grep objects, join files, get/add/edit/delete YAML nodes based on selectors etc.

`[input.yml] => [apply transformations] => [output.yml]`

*Note: The `input.yml` file might contain multiple YAML documents/objects separated by `---`.*

- [Simple YAML transformations](#simple-yaml-transformations)
  - [yaml set $key $value](#yaml-set-key-value)
  - [yaml default $key $value](#yaml-default-key-value)
  - [yaml delete $key](#yaml-delete-key)
  - [yaml join file1.yml file2.yml ...](#yaml-join-file1yml-file2yml)
- [Transformation files](#transformation-files)
  - [yaml apply file1.yml file2.yml ...](#yaml-apply-file1yml-file2yml)
    - [Examples of transformation YAML files](#examples-of-transformation-yaml-files)
- [Print YAML nodes](#print-yaml-nodes)
  - [yaml get "key"](#yaml-get-%22key%22)
    - [kubectl - print pod's main container image](#kubectl---print-pods-main-container-image)
- [Grep objects](#grep-objects)
  - [yaml grep "key: value" ...](#yaml-grep-%22key-value%22)
    - [Grep k8s deployment object by name](#grep-k8s-deployment-object-by-name)
    - [Print first container's image of linkerd2 deployment objects](#print-first-containers-image-of-linkerd2-deployment-objects)
- [Useful Kubernetes examples](#useful-kubernetes-examples)
    - [Push all non-pod objects to k8s](#push-all-non-pod-objects-to-k8s)
    - [Rollout k8s deployments from desired-state files sequentially](#rollout-k8s-deployments-from-desired-state-files-sequentially)
- [Known issues](#known-issues)
- [Feedback](#feedback)
- [License](#license)

# Simple YAML transformations

## yaml set $key $value
```bash
# Add/overwrite field's value
$ cat input.yml | yaml set "metadata.labels.environment: staging" > output.yml
```

## yaml default $key $value
```bash
# Add default value (if no such value exists yet)
$ cat input.yml | yaml default "metadata.labels.environment: staging" > output.yml
```

## yaml delete $key
```bash
# Delete specific field
$ cat input.yml | yaml delete "metadata.labels.environment" > output.yml
```

## yaml join file1.yml file2.yml ...
```bash
# Join multiple YAML files into one, where all documents/objects are separated by `---`
$ yaml cat k8s-apps/*.yml > output.yml
```

# Transformation files

## yaml apply file1.yml file2.yml ...

```bash
$ yaml cat k8s-apps/*.yml | yaml apply staging.yml enable-linkerd.yml > staging/desired-state.yml
```

staging.yml:
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

enable-linkerd.yml:
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

# Print YAML nodes
## yaml get "key"
### kubectl - print pod's main container image
```bash
$ kubectl get pods/nats-8576dfb67-vg6v7 -o yaml | yaml get spec.containers[0].image
nats-streaming:0.10.0
```

# Grep objects
## yaml grep "key: value" ...
### Grep k8s deployment object by name
```bash
$ cat desired-state.yml | yaml grep "kind: Deployment" "metadata.name: linkerd"
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

2. Wildcard array matching doesn't work
```yaml
match:
    spec.template.spec.containers[*].name: prometheus
```

# Feedback
Any feedback welcome! Please open issues and feature requests..

# License
Licensed under the [MIT License](./LICENSE).
