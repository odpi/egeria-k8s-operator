<!-- SPDX-License-Identifier: CC-BY-4.0 -->
<!-- Copyright Contributors to the ODPi Egeria project. -->


# k8s operator design

This project aims to deliver a k8s operator that can manage the deployment of Egeria platforms and servers in a highly available manner. That configuration can subsequently be changed and updated.


## Connectors

The following connectors are needed for this approach to work

### Configuration Connector

* Must allow multiple platforms to be able to READ configurations
* Only one platform will TYPICALLY WRITE, but for handling redundancy this should allow for multiple to

* This should use a k8s-native approach to storing the documents, potentially as k8s resources (configmap or custom)

### Security Connector
* A nominated 'admin' platform should be able to write configuration
* Only the operator should be able to start/stop/add/remove servers/platforms

### Audit Log Connector
* Audit Log Events should be pushed to a k8s-friendly repository for exploration

### Crux Connector
* The existing crux backend supports a highly available & scalable topology & should be the initial supported backend
* Basic yaml, a chart, or an operator to manage a crux deployment
* Not integrated part of egeria operator (just a requirement to have at least one option that works)

## Deployment sample / charts
* An integrated demo environment should be created that deploys the operator, alongside a Kafka environment, relevant UIs, and a backend such as crux, ideally with samples

# Other required changes

* The React UI should respect controls set by the platform security connector, and ensure it is able to purely read/write server config documents, without performing operations
* Support for templating (replaceable values) in the config document may be required
----
License: [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/),
Copyright Contributors to the ODPi Egeria project.
