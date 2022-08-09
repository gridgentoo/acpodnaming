package main

import (
	"fmt"
	"net/http"
	"os"
	"io/ioutil"
	"encoding/json"
	"regexp"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)
	
func isKubeNamespace(ns_name string) bool {
	nsFlag := false
	nsNaming := regexp.MustCompile(`^kube`)
	if nsNaming.MatchString(ns_name) {
		nsFlag = true
	}
	nsNaming = regexp.MustCompile(`^openshift`)
	if nsNaming.MatchString(ns_name) {
		nsFlag = true
	}

	return nsFlag

}

func (gs *myValidServerhandler) serve(w http.ResponseWriter, r *http.Request) {
	
	var Body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			Body = data
		} else {
			fmt.Fprintf(os.Stderr, "Unable to Copy the Body\n")
		}
 	}

	debug := getEnv("DEBUG", "no")
	if debug == "yes" {
		fmt.Fprintf(os.Stdout,"Webhook revieved a request\n")
	}

	if r.URL.Path != "/validate" {
		fmt.Fprintf(os.Stdout,"Not a Valid URL Path\n")
		http.Error(w, "Not a valid URL Path", http.StatusBadRequest)
		return
	}

		if len(Body) == 0 {
		fmt.Fprintf(os.Stderr, "Unable to retrieve Body from the WebHook\n")
		http.Error(w, "Unable to retrieve Body from the API" , http.StatusBadRequest )
		return
	} else if debug == "yes" {
		fmt.Fprintf(os.Stdout, "Body retrieved\n")
	}

	arRequest := &admissionv1.AdmissionReview{}

	if err := json.Unmarshal(Body , arRequest); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to Marshal the Body Request\n")
		http.Error(w,"Unable to Marshal the Body Request",http.StatusBadRequest)
		return
	}


	raw := arRequest.Request.Object.Raw
	pod := corev1.Pod{}

	if err := json.Unmarshal(raw, &pod); err != nil {
		fmt.Fprintf(os.Stderr,"Error Deserializing Pod\n")
		http.Error(w,"Error Deserializing Pod",http.StatusBadRequest)
		return
	}

	arResponse := admissionv1.AdmissionReview {
		Response: &admissionv1.AdmissionResponse{
			Result: &metav1.Status{Status: "Failure", 
			Message: "The Pod name is NOT up to code with the pod naming requirements", 
			Code: 401},
			UID: arRequest.Request.UID,
			Allowed: false,
		},
	}
	
	podName := gs.podNaming

	if gs.podRegextype == "starts" {
		podName = "^" + podName
	}
	podNamingReg := regexp.MustCompile(podName)
	if podNamingReg.MatchString(string(pod.Name)) || isKubeNamespace(arRequest.Request.Namespace) {
		fmt.Fprintf(os.Stdout, "The Pod is up to the naming standard\n")
		arResponse.Response.Allowed = true
		arResponse.Response.Result = &metav1.Status{Status: "Success", 
		Message: "The Pod is up to Standard", 
		Code: 201}
	}

	arResponse.APIVersion = "admission.k8s.io/v1"
	arResponse.Kind = arRequest.Kind

	resp , resp_err := json.Marshal(arResponse)

	if resp_err != nil {
		fmt.Fprintf(os.Stderr, "Unable to Marshal the response\n")
		http.Error(w,"Unable to Marshal the response", http.StatusBadRequest)
		return
	}

	if _ , w_err := w.Write(resp); w_err != nil {
		fmt.Fprintf(os.Stderr,"Unable to write the response\n")
		http.Error(w,"Unable to write the response", http.StatusBadRequest)
	}
}