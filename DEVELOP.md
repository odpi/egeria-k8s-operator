<!-- SPDX-License-Identifier: CC-BY-4.0 -->
<!-- Copyright Contributors to the ODPi Egeria project. -->


# Developing the k8s operator

This document is to aid any developers working on building the k8s operator. It is not needed to just use the operator (the operator is not ready for use yet!!)
  
## Dependencies

These need to be installed and configured in order to build the k8s operator

* [operator sdk](https://github.com/operator-framework/operator-sdk) - version [1.10.1](https://github.com/operator-framework/operator-sdk/releases/tag/v1.10.1)
* [go](https://golang.org) 1.15.x - install from website, os distro, or homebrew (As of April 2021 1.16 is not supported - see https://github.com/operator-framework/operator-sdk/issues/4690 )
* Other dependencies as documented by [operator-sdk](https://sdk.operatorframework.io/docs/building-operators/golang/installation/) including docker, kubectl, kubernetes
* make - for the build process
* A *nix variant or macOS (shell script usage)
* Ensure you have 'GO111MODULE=on' set

## Known issues/gotchas

* Due to a bug in the toolkit, the kube-proxy-rbac image used is provided
by openshift. See https://github.com/operator-framework/operator-sdk/issues/4813#issuecomment-823669700 for more information. 
## Creation of the template project

These commands use the operator-sdk to create the initial project. Code is then edited manually. The commands used are included here in case we need to rebuild the template in future and remerge in customized code

```
operator-sdk init --domain egeria-project.org --license apache2 --owner 'Contributors to the Egeria project' --project-name 'egeria' --repo 'github.com/odpi/egeria-k8s-operator'                                                                                                                              [11:13:14]
```

```
operator-sdk create api --group egeria --version v1alpha1 --kind EgeriaPlatform   
 operator-sdk create api --group egeria --version v1alpha1 --kind EgeriaServer
```
## Dealing with platform-specifics

The operator SDK will install platform specific binaries when a project is created
(above)

The build scripts do what is needed in the CI environment. For local
use it's recommended to run:
```
make kustomize controller-gen
```
Further down the process, when looking at tests, `make test` will also download
required binaries for testing.

If you get any issues with binaries, clean out the 'bin' and 'testbin' directories
to remove any platform dependent files & then repeat these steps.

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
make install
```
## Building the project and image
```
make docker-build docker-push IMG=odpi/egeria-k8s-operator:0.1.0
```
## deploy the operator
```
make install
make deploy IMG=odpi/egeria-k8s-operator:0.1.0
```
## Check the operator controller is active
```
kubectl get deployment -n egeria-system 
```
## Checking logs of the controller
```
 kubectl get pods -n egeria-system 
```
Then use that pod id in the entry below:
```
kubectl logs egeria-k8s-operator-controller-manager-6bf887c74c-78mwc -n egeria-k8s-operator-system manager

```
## Create an instance of Egeria
```
 kubectl apply -f config/samples/egeria_v1alpha1_egeriaplatform.yaml         
```
## Lookiing for instances of the egeria CRD:
```
kubectl get EgeriaPlatform
```
## Changing properties of the sample instance (name returned from above)
```
kubectl edit EgeriaPlatform/egeriaplatform-sample
```
ie change the size to 10 to scale
## Delete the instance
```
kubectl delete EgeriaPlatform/egeriaplatform-sample
```
## Cleaning up the crd after
```
 kubectl delete -f config/crd/bases/egeria.egeria-project.org_egeriaplatforms.yaml        
 kustomize build config/default | kubectl delete -f -
```
# Design decisions

* We are using operator-sdk 1.0.0 for tooling
* golang is the implementation language for the oeprator
* operator is [cluster-scoped](https://sdk.operatorframework.io/docs/building-operators/golang/operator-scope/) - this is the default and can be revisited in future
* operator uses a single [group](https://book.kubebuilder.io/cronjob-tutorial/gvks.html) 'org.odpi.egeria' for it's APIs - also the default
* The initial implementation uses a single [kind](https://book.kubebuilder.io/cronjob-tutorial/gvks.html) called 'Egeria' - think of this as the k8s resource type we are dealing with
* We will start with version 1 ('v1')
 * 
# Bugs

# Discussion
 
See the #egeria-k8s channel on slack at https://slack.odpi.org

----
License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/),
Copyright Contributors to the ODPi Egeria project.
