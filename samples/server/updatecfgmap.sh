#!/bin/zsh
# SPDX-License-Identifier: CC-BY-4.0
# Copyright Contributors to the ODPi Egeria project.

# These are aggregate configmaps - no longer used as we currently need one config per server. We could access by key, but for now go with 1:1 (below)
# kubectl create configmap corecfg --from-file cocoMDS2 --from-file  cocoMDS5 --from-file cocoMDS6 --from-file cocoMDS3 -o yaml --dry-run=client | kubectl replace --force -f -
# kubectl create configmap lakecfg --from-file cocoMDS1 --from-file cocoMDS4 --from-file cocoView1 -o yaml --dry-run=client | kubectl replace --force -f -
# kubectl create configmap devcfg --from-file cocoMDSx -o yaml --dry-run=client | kubectl replace --force -f -

for file in coco*
do
  echo "Processing $file"
  # filename forced to lower case (configmaps MUST be lowercase RFC 1123 subdomain)
  kubectl create configmap ${file:l} --from-file $file -o yaml --dry-run=client | kubectl replace --force -f -
done


