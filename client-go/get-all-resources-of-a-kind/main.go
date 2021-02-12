package main

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.APIPath = "/api"
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}

	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}

	var res unstructured.UnstructuredList
	res.DeepCopyObject()
	err = client.Get().
		NamespaceIfScoped("", true).
		Resource("pods").
		Context(context.Background()).
		Do().
		Into(&res)
	if err != nil {
		panic(err)
	}

	for _, pod := range res.Items {
		fmt.Printf("%s/%s\n", pod.GetNamespace(), pod.GetName())
	}
}
