<!-- SPDX-License-Identifier: CC-BY-4.0 -->
<!-- Copyright Contributors to the ODPi Egeria project. -->


# Developing the k8s operator

This document is to aid any developers working on building the k8s operator. It is not needed to just use the operator (the operator is not ready for use yet!!)
  
## Dependencies

These need to be installed and configured in order to build the k8s operator

* [operator sdk](https://github.com/operator-framework/operator-sdk) - version [1.0.0](https://github.com/operator-framework/operator-sdk/releases/tag/v1.0.0)
* [go](https://golang.org) 1.15 - install from website, os distro, or homebrew
* Other dependencies as documented by [operator-sdk](https://sdk.operatorframework.io/docs/building-operators/golang/installation/) including docker, kubectl, kubernetes
* make - for the build process
* A *nix variant or macOS (shell script usage)

## Creation of the template project

These commands use the oeprator-sdk to create the initial project. Code is then edited manually. The commands used are included here in case we need to rebuild the template in future and remerge in customized code

```
operator-sdk init --plugins "go.kubebuilder.io/v2" --project-name 'egeria-k8s-operator' --repo 'github.com/odpi/egeria-k8s-operator'
```

```
operator-sdk create api --group egeria --version v1 --kind Egeria --resource=true --controller=true
```
## Changing the API model

This is needed if the egeria type is modified -- it keeps the go type definitions in sync
```
make generate
```
Then we need to build the new CRD with
```
make manifests
```

## Changing the controller

```
make manifests
make install
make docker-build docker-push IMG=odpi/egeria-k8s-operator:0.0.1
```
## Building the project and image
```
make docker-build docker-push IMG=odpi/egeria-k8s-operator:0.0.1
```
## deploy the operator
```
make install
make deploy IMG=odpi/egeria-k8s-operator:0.0.1
```

# Design decisions

* We are using operator-sdk 1.0.0 for tooling
* golang is the implementation language for the oeprator
* operator is [cluster-scoped](https://sdk.operatorframework.io/docs/building-operators/golang/operator-scope/) - this is the default and can be revisited in future
* operator uses a single [group](https://book.kubebuilder.io/cronjob-tutorial/gvks.html) 'egeria' for it's APIs - also the default
* The initial implementation uses a single [kind](https://book.kubebuilder.io/cronjob-tutorial/gvks.html) called 'Egeria' - think of this as the k8s resource type we are dealing with
* We will start with version 1 ('v1')
----
License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/),
Copyright Contributors to the ODPi Egeria project.
