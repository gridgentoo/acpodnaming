# Testing the Admission controller 

## Steps
The following steps will describe how to test the validation admission controller 

1. create a namespace
2. add the lable for the namespace
3. try to create a Pod with the string
4. without the string

### Creating the namespace 
First let's create the namespace (as in any case)
```bash
# oc new-project acpodnaming-testing
```

### Label the namespace
In order for the admission controller to work we need to add a label to the namespace 
```bash
# oc label namespace acpodnaming-testing "admission.kubernetes.io/podnaming=True"
```

### create the Pass Pod
Now let's create a pod with the name string as we defined 
```bash
# cat << EOF | oc apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: kuku-ubi
  namespace: acpodnaming-testing
spec:
  containers:
  - name: kuku-ubi
    image: ubi8/ubi-minimal
    command: ["/bin/tail"]
    args:
    - "-f"
    - "/dev/null"
EOF
```

### Create the Failed Pod
```bash
# cat << EOF | oc apply -f -
apiVersion: v1
kind: Pod
metadata:
  name: simple-ubi
  namespace: acpodnaming-testing
spec:
  containers:
  - name: simple-ubi
    image: ubi8/ubi-minimal
    command: ["/bin/tail"]
    args:
    - "-f"
    - "/dev/null"
EOF
```
Now you should recieved the following error message :
```bash
Error from server: error when creating "STDIN": admission webhook "acpodnaming.kubernetes.io" denied the request: The Pod name is NOT up to code with the pod naming requirements
```

That is it
Good Luck