package main

import (
	"k8s.io/api/core/v1"
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// getUnstructuredList fetches an k8s.io/apimachinery/pkg/apis/meta/v1/unstructured.UnstructuredList of the given
// GVK and transforms it into the type denoted by res. This comes in handy when res doesn't implement the
// k8s.io/apimachinery/pkg/runtime.Object interface so you can't just c.List() it.
func getUnstructuredList(c client.Client, gvk schema.GroupVersionKind, res interface{}) error {
	var unstructuredList unstructured.UnstructuredList
	unstructuredList.SetGroupVersionKind(gvk)
	err := c.List(context.Background(), &unstructuredList)
	if err != nil {
		return err
	}
	unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&unstructuredList)
	if err != nil {
		return err
	}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredObject, res)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	v1.PodCondition
	
	controllerruntime.SetLogger(zap.Logger(true))
	c, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", c)
}
