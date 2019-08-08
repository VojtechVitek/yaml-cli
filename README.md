# yaml transformer
A CLI tool for transforming YAML files (add, edit, delete YAML nodes based on selector)

```
[input YAML] => [apply given transformations] => [output YAML]
```

The input might contain multiple YAML documents separated by `---`.

## Work in progress.. not stable yet

## Example:

deployment.yml:
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

staging.yml:
```yml
match:
    kind: Deployment
    metadata.name: api
set:
    metadata.labels.environment: staging
    metadata.labels.first: updated-label
    spec.replicas: 3
```

```bash
cat deployment.yml | yaml apply staging.yml
```
```yml
apiVersion: apps/v1
kind: Deployment
metadata:
    name: api
    labels:
        environment: staging
        first: updated-label
        second: label
spec:
    replicas: 3
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

Any feedback welcome! Please open issues and feature requests..
