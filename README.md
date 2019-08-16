# yaml transformer
A CLI tool for transforming YAML files (add, edit, delete YAML nodes based on selectors)

`[input.yml] => [apply transformations] => [output.yml]`

*Note: The input file might contain multiple YAML documents separated by `---`.*

## == Work in progress.. not stable yet! ==

## Simple operations:
```bash
# Add/overwrite field's value.
$ cat input.yml | yaml set "metadata.labels.environment" "staging" > output.yml

# Delete specific field.
$ cat input.yml | yaml delete "metadata.labels.environment" > output.yml

# Add default value (if no such value exists yet)
$ cat input.yml | yaml default "metadata.labels.environment" "staging" > output.yml

# Grep YAML object(s):
$ cat desired-state.yml | yaml grep "kind: Deployment"
$ cat desired-state.yml | yaml grep -v "kind: [Deployment, Pod]"
```

## Kubernetes examples

### Rollout deployments sequentially
```bash
# Apply all K8s objects except for Deployments first:
$ cat desired-state.yml | yaml grep -v "kind: Deployment" | kubectl apply -f -

# Apply deployments sequentially, ie. define order in which specific microservices are rolled out.
$ for app in embedder api frontend; do
    cat desired-state.yml | yaml grep "kind: Deployment" "metadata.name: $app" | kubectl apply -f -
    kubectl rollout status --timeout 120 deploy/$app
    if [ $? -ne 0 ]; then
        kubectl rollout redo deploy/$app
        # Notify about failed deployment.
        exit 1
    fi
  done
```

## Match the YAML objects before applying transformation:
```bash
$ cat k8s-desired-state.yml | yaml match "kind: Deployment" "metadata.name: redis" set "spec.template.spec.containers[0].image" "redis:5.0.5" > output.yml
```

## Apply transformations files:
```bash
$ cat deployment.yml | yaml apply staging.yml > desired-state.yml
```

Original object (deployment.yml):
```yml
apiVersion: apps/v1
kind: Deployment
metadata:
    name: api
    labels:
        first: label
        second: label
spec:
    replicas: 1
```

Transformation (staging.yml):
```yml
match:
    kind: Deployment
    metadata.name: api
set:
    metadata.labels.first: updated-label
    metadata.labels.environment: staging
    spec.replicas: 3
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
 spec:
-    replicas: 1
+    replicas: 3
```

## Goals
I was frustrated by the K8s Kustomize tool and its useless error messages.

I'm seeking a simple tool to transform YAML files, ie. match some Kubernetes objects and apply some common transformations onto the fields per specific environment.

```bash
$ cat deployment.yml | yaml apply staging/transform.yml [...] > _desired/staging-deployment.yml`
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

1. Merging complex nodes
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

2. Match with arrays
```yaml
match:
    spec.template.spec.containers[0].image: pressly/api
```
Doesn't look at the containers[] array yet.

### Feedback
Any feedback welcome! Please open issues and feature requests..
