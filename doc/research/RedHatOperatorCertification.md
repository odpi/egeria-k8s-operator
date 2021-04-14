# RedHat Operator Certification

Red Hat is one organization promoting the use of operators on it's Kubernetes 
platform, and offers a certification process

This document will summarise the requirements for certification insofar as they
inform the design and work activities around the Egeria Operator - with the expectation 
that some of these are 'good ideas' that we should incorporate in our design.

Source: https://connect.redhat.com/partner-with-us/red-hat-openshift-operator-certification

## Benefits to operator developer

* Appears in Red Hat Marketplace
* Appears in OperatorHub (within openshift)
* Appears in the Red Hat Ecosystem Catalog
* Findability

## Benefits to User of operator

* Easier to find solutions
* Collaborative support Red Hat + development project

## General Requirements for operators

* Container images are trusted/secure - based on UBI 
  (Core Egeria already is, but we use other containers in our demos that aren't)
* automated operations (need more clarity on this)
* SPecific additional requirements for CNI/CSI integration for networking/storage solutions (not so relevant for Egeria)

## Operator Distribution

According to https://docs.openshift.com/container-platform/4.7/operators/understanding/olm-understanding-operatorhub.html#olm-understanding-operatorhub there is a 
category known as 'Community Operators' which do not have official support, but are
linked to their dev team through github. This is probably the most suitable
category for Egeria - but EXCLUDES availability in the marketplace.

