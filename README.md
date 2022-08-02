# the Admission Controller Pod Naming 

This Repository is for the Pod Naming Admission Controller. The Pod Naming Admission Controller is
A simple example on how to write the must basic admission controller for testing if the name of the pod
contains a string which is mentioned by the POD_NAMING environment variable.

## How to depoy 

The Deployment process it very simple. it focus on the fast that we deploy it on OpenShift 4 and we let
the Certificate operator to manage the TLS certificate communication between the Kubernetes API to the 
webhook we have created.

the Deployment Steps are :

1. create the namespace
2. create the service (that will create the secret)
3. create the deployment resource.
4. create the validating webhook resource and point it to the service

### creating the Name Space