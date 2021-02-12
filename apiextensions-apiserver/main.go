package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	v1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}

	crdClient := v1.NewForConfigOrDie(config)
	crd, err := crdClient.CustomResourceDefinitions().Get(context.Background(), "federatedconfigmaps.types.kubefed.io", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	crdBytes, err := json.Marshal(crd.Spec.Versions[0].Schema)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", crdBytes)
}
