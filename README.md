# barn

Barn watches remote Kubernetes manifests and applies them to your cluster.

# Purpose

Think of the `kubelet` but for your cluster. The `kubelet` runs on a single node and watches a source of Kubernetes
manifests, e.g. on some path on the filesystem or a remote HTTP endpoint, and tells the container runtime on that node
to run the respective containers.

Barn runs in a single cluster and watches a source of Kubernetes manifests and tells the Kubernetes API server of that
cluster  to run the respective manifests. Currently, any OAuth2 protected HTTP endpoint is supported as a source.

# Usage

```
$ barn \
    --cluster=https://kubernetes.default.svc.cluster.local \
    --source=https://raw.githubusercontent.com/kubernetes/kubernetes/master/examples/explorer/pod.yaml \
    [--source=another-source ... \]
    [--token=authentication-not-needed-in-this-case \]
```

This watches the file `./examples/explorer/pod.yaml` of the `github.com/kubernetes` repostiory on the `master` branch.
If there's any changes pushed to that file on the master branch, `barn` will apply them.

Note, that any HTTP server that serves `yaml` or `json` files and uses the `Authorization: Bearer <token>` method to
authenticate clients will work. So instead of an HTTP server pointing to your Github repository content, as seen above,
you could easily run your own more tailored solution, e.g. https://manifests.me/?cluster=foo&env=bar could return
dynamicly generated `yaml` specifically for that cluster.

# Caveats

The controller will only `kubectl apply` the files it watches, so it cannot handle deletions currently.
Furthermore, controllers that don't handle `apply` correctly will also not work as expected, e.g. daemon sets aren't
updated by `apply` at the moment.
