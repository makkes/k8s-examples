package main

import (
	"context"
	"fmt"
	"os"

	rbacv1 "k8s.io/api/rbac/v1"
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

	rbacv1.AddToScheme(scheme.Scheme)
	config.ContentConfig.GroupVersion = &rbacv1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}
	client, err := rest.RESTClientFor(config)
	if err != nil {
		panic(err)
	}
	newRole := rbacv1.Role{}
	newRole.SetName("a-new-role")
	res := rbacv1.Role{}
	if err := client.
		Post().
		Resource("roles").
		Namespace("default").
		Body(&newRole).
		Do(context.Background()).
		Into(&res); err != nil {
		panic(err)
	}
	fmt.Printf("New Role:\n%#v\n", res)
}
