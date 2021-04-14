# Proposed Egeria Operator capability

Having reviewed a number of pre-existing operators on github and in catalogs 
this is a brief summary of the kind of capability that makes sense in an operator

We need to put ourselves in the shoes of those trying to deploy egeria. 
Let's all be garygeeke perhaps, or one of his colleagues. This isn't about the details
of creating assets, but it is about the deployment of the software components
needed to run egeria, and the logical deployment we build on top of that

## Kubernetes

Kubernetes is assumed here, since we're talking about an operator, so this provides
the platform upon which we can run code (think deployments, pods, jobs etc) and 
services (ie connectivity as well as manage storage needs)

# What are the Things?

So what are the entities that we deal with when deploying egeria?
Without thinking too much about client side:

## that we can run (now)...

* Egeria Server Chassis
* UI Server chassis
* Static content server (UI)
* React based UI
* Storage back-end for JanusGraph

So basically we have our general compute server, storage & UI

## but is modelling this alone helpful?

To some extent -- we can easily see how we might create multiple 
unique instances of the above, and/or replicas - very core k8s 
function.

But it may not really get into the true operational things involved
in deploying egeria -- The Egeria Server Chassis is very rich, and can itself
be thought of as a container for running many other things - metadata servers, view servers, discovery services, metadata connectors

To truly make operations easier we need to focus at this level too

## So what shall we do?

When developing an operator, a 'CRD' or 'Custom Resource Descriptor' provides a definition
of what an entity should look like - whilst the operator code itself acts
as Captain Picard and 'Makes it So'

A) We could go with two CRDs.

The first is the 'EgeriaPlatform' CRD. This is relatively simple (!) as it will
be deploying containers to run processes, and expose services

This could offer capabilities such as
 - deploying 1 or more egeria platforms
 - deploying 1 or more UI front ends