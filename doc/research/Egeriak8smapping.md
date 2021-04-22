# Mapping Egeria Server/Platform to Kubernetes

## Egeria Platform

The egeria platform (aka 'server chassis') is effectively a java process that is launched and acts as a container for
running Egeria servers (below). In itself it has relatively little configuration but does include
 - TLS configuration ie passwords, keys, certs
 - Spring configuration (security, plugins)

It offers an admin interface for defining servers (below) and starting/stopping them

The platform's connectors are mostly limited to configuration

## Egeria Servers

The primary resource dealt with when deploying egeria is the 'server', which as diagrammed, comes in a number of 
different forms including a repository proxy, a metadata repository, a view server etc

A 'server' is a logical concept that actually executes code within an Egeria Platform but is basically what is defined, started, stopped,
and also hosts the myriad of REST APIs that egeria supports such as for creating the definition of a data asset.

## Egeria Connectors

Much of the capability in Egeria is pluggable. For example we have a connector to the message bus, which could be kafka, but equally
may be RabbitMQ (if a connector were written). We have connectors that understand data sources such as for Postgres. These are 
implemented as java jars, and are added into the classpath of the egeria platform. Thus the platform can provide a certain capability
with these additional connectors, but it is the server which defines which ones to use (and refers to the class name to locate)

## Pods & Containers

The egeria platform is what is runnable, so the platform is what needs to be running in a container (most likely a RedHat UBI image, as we use today)


## Services

Egeria services take the form of REST APi calls to, for example,  GET https://platformhost:port/servers/myserver/open-metadata ....
Since we need servers to be highly available ie replicated, we need a 'service' which can point to any pod that is running a platform with this service
healthchecks need to be targeted at the specific server, not just the platform
The routing decision is probably best aided by additional metadata tags we can use as a selector - this is superior to simple DNS routing which might be done with a headless service
Each platform should also support a service, but this may need to be restricted to read-only for most users - or a subset - via configuration to avoid interacting with k8s

## Mapping Servers to Platforms

Three possible options are

### Anonymous, homogenous platforms

In this model we have a pool of 'n' platforms which are all configured identically
A server can therefore be configured/activated on any platform that is running

If we had 4 servers, we could run these across 2 platforms, or 4, or 6, 

### One platform, one server

Here we don't explicitly define the platform. Instead we always create one platform per server (replica). All configuration
is associated with the server in a single CRD including TLS, spring etc.

This avoids the need for the deployer to be concerned with the platform.

One downside is that 4 servers would effectively need 8 platforms, 8 servers would need 16 etc - which could end up
being inefficient if these servers are relatively small. 

### Explicit platforms and server

With this option we adopt a model whereby the deployer can explicitly create various platform configurations - say type1, type2, type3. 
Each of these platform configurations can have attributes like a number of replicas
The pods running a platform would expose a metadata tag including the name of the type

The deployer will also define a server, and will refer to a platform type
The operator will configure all platform instances of a type with the new server, but will only activate the server to the point where the desired number of
server replicas is reached
The platform will additionally get assigned metadata with the server name
A service will use this metadata to route requests



