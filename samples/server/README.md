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

----
License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/),
Copyright Contributors to the ODPi Egeria project.
