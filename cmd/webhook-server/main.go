package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	admission "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	tlsDir      = `/run/secrets/tls`
	tlsCertFile = `tls.crt`
	tlsKeyFile  = `tls.key`
)

var (
	podResource = metav1.GroupVersionResource{Version: "v1", Resource: "pods"}
	conf, _     = readConf("image_conf.yaml")
)

func applyRegistryDefaults(req *admission.AdmissionRequest) ([]patchOperation, error) {
	// This handler should only get called on Pod objects as per the MutatingWebhookConfiguration in the YAML file.
	// However, if (for whatever reason) this gets invoked on an object of a different kind, issue a log message but
	// let the object request pass through otherwise.
	if req.Resource != podResource {
		log.Printf("expect resource to be %s", podResource)
		return nil, nil
	}

	type containerDetail struct {
		name  string
		image string
	}

	// Parse the Pod object.
	raw := req.Object.Raw
	pod := corev1.Pod{}
	var allContainers = make([]containerDetail, 0)
	if _, _, err := universalDeserializer.Decode(raw, nil, &pod); err != nil {
		return nil, fmt.Errorf("could not deserialize pod object: %v", err)
	} else {
		for _, s := range pod.Spec.Containers {
			allContainers = append(allContainers, containerDetail{name: s.Name, image: s.Image})
		}
	}

	var errorContainers = make([]containerDetail, 0)

	//TODO: remove when all options for conf.Allow.Empty is implemented.
	if conf.Allow.Empty == "mutate" || conf.Allow.Empty == "true" {
		return nil, fmt.Errorf("configuration is not implemented yet: %v", conf.Allow.Empty)
	}

	for _, s := range allContainers {
		if !containerImageISOK(s.image, conf) {
			errorContainers = append(errorContainers, s)
		}
	}

	var errorMessage string = "Error in image, you need to fix"
	for _, s := range errorContainers {
		errorMessage = fmt.Sprintf("%s %s:%s", errorMessage, s.name, s.image)
		log.Println(errorMessage)
	}

	//TODO: Implement conf.Allow.Empty='mutate' here.
	var patches []patchOperation
	/*
		patches = append(patches, patchOperation{
			Op:   "add",
			Path: "/spec/containers/<a container name>",
			Value: <a-value>
		})
	*/

	if len(errorContainers) > 0 {
		return nil, errors.New(errorMessage)
	}

	return patches, nil
}

func main() {
	c, err := readConf("image_conf.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", c)

	certPath := filepath.Join(tlsDir, tlsCertFile)
	keyPath := filepath.Join(tlsDir, tlsKeyFile)

	mux := http.NewServeMux()
	mux.Handle("/mutate", admitFuncHandler(applyRegistryDefaults))
	server := &http.Server{
		// We listen on port 8443 such that we do not need root privileges or extra capabilities for this server.
		// The Service object will take care of mapping this port to the HTTPS port 443.
		Addr:    ":8443",
		Handler: mux,
	}
	log.Fatal(server.ListenAndServeTLS(certPath, keyPath))
}
