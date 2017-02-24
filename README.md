# manifest-controller
[![Build Status](https://travis-ci.org/linki/manifest-controller.svg?branch=master)](https://travis-ci.org/linki/manifest-controller)
[![Build Status](https://drone.factorio.linki.space/api/badges/linki/manifest-controller/status.svg)](https://drone.factorio.linki.space/linki/manifest-controller)
[![Coverage Status](https://coveralls.io/repos/github/linki/manifest-controller/badge.svg?branch=master)](https://coveralls.io/github/linki/manifest-controller?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/linki/manifest-controller)](https://goreportcard.com/report/github.com/linki/manifest-controller)

manifest-controller watches remote Kubernetes manifests and applies them to your cluster.

# Purpose

Think of the `kubelet` but for your cluster. The `kubelet` runs on a single node and watches a source of Kubernetes
manifests, e.g. on some path on the filesystem or a remote HTTP endpoint, and tells the container runtime on that node
to run the respective containers.

manifest-controller runs in a single cluster and watches a source of Kubernetes manifests and tells the Kubernetes API server of that
cluster  to run the respective manifests. Currently, any OAuth2 protected HTTP endpoint is supported as a source.

# Usage

```
$ manifest-controller \
    --cluster=http://127.0.0.1:8001 \
    --source=https://raw.githubusercontent.com/kubernetes/kubernetes/master/examples/elasticsearch/es-svc.yaml \
    --source=https://raw.githubusercontent.com/kubernetes/kubernetes/master/examples/elasticsearch/es-rc.yaml \
    --source=https://raw.githubusercontent.com/kubernetes/kubernetes/master/examples/elasticsearch/service-account.yaml \
    [--source=another-source ... \]
    --token=a-personal-github-access-token \
```

This watches the three example manifest files of the `github.com/kubernetes` repostiory on the `master` branch.
If there's any changes pushed to those files on the master branch, `manifest-controller` will apply them.

Note, that any HTTP server that serves `yaml` files and uses the `Authorization: Bearer <token>` method to
authenticate clients will work. So instead of an HTTP server pointing to your Github repository content, as seen above,
you could easily run your own more tailored solution, e.g. https://manifests.me/?cluster=foo&env=bar could return
dynamicly generated `yaml` specifically for that cluster.

# Caveats

The controller will only `kubectl apply` the files it watches, so it cannot handle deletions currently.
Furthermore, API objects that don't handle `apply` correctly will also not work as expected, e.g. daemon sets aren't
updated by `apply` at the moment.
