# yaml CLI transformer
A CLI tool for transforming YAML files (add, edit, delete YAML nodes based on selectors)

`[input.yml] => [apply transformations] => [output.yml]`

*Note: The input file might contain multiple YAML documents separated by `---`.*

## == Work in progress.. not stable yet! ==

## Simple YAML transformations

```bash
# Add/overwrite field's value
$ cat input.yml | yaml set "metadata.labels.environment" "staging" > output.yml
```

```bash
# Delete specific field
$ cat input.yml | yaml delete "metadata.labels.environment" > output.yml
```

```bash
# Add default value (if no such value exists yet)
$ cat input.yml | yaml default "metadata.labels.environment" "staging" > output.yml
```

## Applying transformations from multiple files

```bash
$ yaml cat *.yml | yaml apply staging.yml enable-linkerd.yml > staging/desired-state.yml
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
```

## Helpful CLI tools:

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

### Push all non-pod objects to k8s:
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

## Examples of transformation YAML files:
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

### Known issues

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

### Feedback
Any feedback welcome! Please open issues and feature requests..
