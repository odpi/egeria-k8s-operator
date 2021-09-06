#!/bin/sh
# SPDX-License-Identifier: CC-BY-4.0
# Copyright Contributors to the ODPi Egeria project.

# For the coco labs we do
kubectl create configmap corecfg --from-file cocoMDS1 --from-file cocoMDS2 --from-file cocoMDS3 --from-file cocoMDS4 --from-file   cocoMDS5 --from-file cocoMDS6 --from-file cocoMDSx -o yaml --dry-run=client | kubectl replace --force -f -
kubectl create configmap lakecfg --from-file cocoMDS1 --from-file cocoMDS4 --from-file cocoView1 -o yaml --dry-run=client | kubectl replace --force -f -
kubectl create configmap devcfg --from-file cocoMDSx -o yaml --dry-run=client | kubectl replace --force -f -

