# Build the Pod on OpenShift

In case we want to take the code and build the Pod in our own environment we can use the following
steps to build and use the OpenShift tools to create the Image.

## steps

1. create an image with buildah and the source code 
2. generate the image 
3. use the image to build the acpodnaming pod
4. use the code within Openshift's internal registry

### create a custom build image

In order to create a custom build image we can use the Docker in the corrent directory with the build.sh
script and use the oc commands to genreate it (in our new Namespace)

#### Creating the Name Space

```bash
# oc apply -f ../Yamls/namespace.yaml
```
