# yaml transformations
A CLI tool for transforming YAML files (add, edit, delete YAML nodes based on selector)

## Work in progress.. not stable yet

## Goals
I was frustrated by the K8s Kustomize tool and its useless error messages.

I'm seeking a simple tool to transform YAML files, ie. match some Kubernetes objects and apply some common transformations onto the fields per specific environment.

```bash
$ cat deployment.yml | yaml apply staging/transform.yml [...] > _desired/staging-deployment.yml`
```

## Ideas for syntax

### Syntax for advanced transformations
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

## Feedback welcome! Please open issues and feature requests..
