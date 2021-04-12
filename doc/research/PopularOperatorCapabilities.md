# Capabilities of popular Kubernetes Operators (April 2021)

## Purpose

Review popular operators & look at the kind of operators they offer, to help
inform decisions we make about what the Egeria Operator should be able to do

I will aim to look at operators relating to similar technologies as Egeria - 
so not focussing on those that support kubernetes, helm etc, but rather typical
middleware components or applications.

## Sources
https://github.com/operator-framework/awesome-operators

## Summary

## Detailed Analysis

### Tomcat

https://github.com/kube-incubator/tomcat-operator

This is a very simple operator, not updated for 2 years. It allows configuration
of a tomcat instance using a single CRD (of which of course there many be many instances of)
 - replicas
 - docker image & pull policy
 - a web archive to serve
 - the deployment directory to use

There's not a lot to see here beyond an early example of using operator-sdk

### Wordpress

https://github.com/presslabs/wordpress-operator

This is a more sophisticated example with many releases up to July 2020, but at
that point development appears to have stopped with just a few minor commits since that
point.

Configuration includes
 - replicas 
 - source code location for wordpress code -- such as from git (provide credentials, path)
 - media location. Typically in google cloud, could also be a volume
 - basic wordpress parms including  user, password, email, title
 - deployment info for the wp db such as host, user, password, 
 - TLS certificate
 - ingress annotations

The focus here is on initial deployment rather than day 2 operations

### PostGres

https://github.com/zalando/postgres-operator

This is one of several postgres operators. A UI is provided.

New clusters can be created with configuration include
 - name
 - namespace
 - team ownership
 - postgres version
 - replicas
 - DNS names
 - load balancing configuration
 - data volume size
 - resource limits
 - users (list)
 - databases (list)


### Open Liberty

https://github.com/rabbitmq/cluster-operator

This operator allows for applications to be deployed and managed by open liberty.

The key attributes include:
 - application name, version
 - container image name, and pull policy/secret
 - any additional containers needed for initialization / sidecars
 - service bindings, ports, types, routes, sso configuration
 - certificates, passwords
 - dependencies on other open liberty services
 - replicas & scaling - min, max, cpu, request management
 - storage volumes, sizes, classes
 - additional CRDs to configure logging & traces

### RabbitMQ

https://github.com/rabbitmq/cluster-operator

This allows for simple cluster creation

 - rabbitmq configuration parameters (intervals, mem settings etc)
 - storage volume size
 - resource limits
 - replicas

### MongoDB

https://docs.mongodb.com/kubernetes-operator/master/tutorial/install-k8s-operator/

An ops manager can be configured
 - replicas
 - version
 - credentials
 - connectivity type (ie load balancer)
 - application db 
 - backup target, credentials for s3
 - providing an aggregated status overview

You can then
 - create a mongoDB instance
 - more extenstive settings are in a configmap

### mySQL 

https://github.com/presslabs/mysql-operator

* Backups (s3)
* Cluster deployment (inc credentials), various config parms
* link to predefined storage
* restore cluster from backup
* orchestrator - topology mapping, master/slave management, replication rules

### Yugabyte DB

https://github.com/yugabyte/yugabyte-operator

* Deploy a new cluste
* replicas for master & TServers
* storage
* resource scheduling/management
* additional config flags

### Data Explorer operator

https://github.com/rhm-samples/data-explorer-operator

Deploys an example environment including notebooks
Minimal configuration