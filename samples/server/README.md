<!-- SPDX-License-Identifier: CC-BY-4.0 -->
<!-- Copyright Contributors to the ODPi Egeria project. -->

This directory contains sample configuration files from the Egeria Coco Pharmaceuticals lab.

These have been extracted from a cloud deployed server, so hostnames/ports for message bus, local addresses etc
will need to be changed.

File change history is also included in each config file.

The servers with a metadata store make use of an in-memory repository, so these configurations will
not work correctly when replicated as each replica will have a unique repository.

One future likely approach will be for the 'server author' GUI (part of our ReactUI) to create
these configurations.

# Loading config files into kubernetes objects

Example:
```shell
$ kubectl create configmap servercfg --from-file cocoMDS1 --from-file cocoMDS2 --from-file cocoMDS3 --from-file cocoMDS4 --from-file cocoMDS5 --from-file cocoMDS6 --from-file cocoMDSx -o yaml --dry-run=client | kubectl replace --force -f -
configmap "servercfg" deleted
configmap/servercfg replaced
```

# Coco Pharma configuration

The script 'updatecfg' loads up three configmaps to match the coco pharma environment, as setup for the notebooks
* corecfg (cocoMDS 2/5/6/3)
* lakecfg (cocoMDS 1/4 & cocoView1)
* devcfg (cocoMDSx)

These can then be referred to in the associated CR for the platform

Any changes to configmaps are visible to platforms which already have them mounted.
Note that some caching of config maps does occur - for now wait 'a while' before updating any CR to reference them

----
License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/),
Copyright Contributors to the ODPi Egeria project.
