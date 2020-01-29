# Kubernetes programming examples

This repo is a place for collecting small snippets of Go code targeting a
specific task each, i.e. simple controllers, API clients etc.

## What you'll find here

The examples are split by the library they mainly employ:

* [controller-runtime](controller-runtime/): Examples making use of
  [sigs.k8s.io/controller-runtime](https://github.com/kubernetes-sigs/controller-runtime).
* [client-go](client-go/): Examples making use of
  [k8s.io/client-go](https://github.com/kubernetes/client-go).

## How to use the examples

All the examples in this repo depend on the `KUBECONFIG` variable pointing to a
kubeconfig file. So running them is as simple as

```sh
$ KUBECONFIG="/home/dan/.kube/config" go run main.go
```

## How to contribute

PRs are very welcome! Just open a PR adding a directory with your `main.go` and
any other accompanying files. ❤️
