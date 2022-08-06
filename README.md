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
2. create the image on OpenShift
3. create the service (that will create the secret)
4. create the deployment resource.
5. create the validating webhook resource and point it to the service

### creating the Name Space

In a very simple way we can just create the yaml file 

```bash
# oc apply -f Yamls/namespace.yaml
```

### Running the build
Once the namespce is created we can start the build process which describes [here](build.md)

### Creating the Service

We need to create the service before we create the deployment because the service creates the secret 
the Pods needs to load and run internally 
```bash
# oc apply -f Yamls/service.yaml
```

### Run the Deployment
Now we can run the deployment with secret which was created automatically from the service :
I would recommand to change the POD_NAMING evironment variable to fit your needed naming request

```bash
# export MY_NAME=<your name>
# sed -i "s/kuku/$MY_NAME/g" Yamls/deployment.yaml
```
Now we can deploy the pods

```bash
# oc apply -f Yamls/deployment.yaml
```

For the final step we can deploy the validation webhook configuration :
```bash
# oc apply -f Yamls/ValidatingWebhookConfiguration.yaml
```

## Testing

if you want to test your deployment you can follow the steps describes [here](test.md)

That is it 

If you have any question feel free to responed/ leave a comment.  
You can find on linkedin at : https://www.linkedin.com/in/orenoichman  
Or twitter at : https://twitter.com/ooichman  