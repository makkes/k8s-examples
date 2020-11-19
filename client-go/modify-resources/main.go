/*
Example for modifying existing resources in a cluster using client-go.

This code will retrieve a deployment called "whoami" from the "default"
namespace and change its first containers args' "--port" flag to a random
value between 1024 and 65535.

Prerequisites:

- Use the "deployment.yaml" file to create the deployment before running this
code.
- Set the "KUBECONFIG" environment variable.
*/
package main

import (
	"context"
	"math/rand"
	"os"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}

	clientset := kubernetes.NewForConfigOrDie(config)
	deploymentsClient := clientset.AppsV1().Deployments("default")

	// we're using the retry package's helper function RetryOnConflict here
	// so that parallel modifications of the deployment don't cause our code
	// to just fail but rather retry until it succeeded.
	err = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		deployment, err := deploymentsClient.Get(context.Background(), "whoami", metav1.GetOptions{})
		if err != nil {
			panic(err)
		}

		args := deployment.Spec.Template.Spec.Containers[0].Args
		for idx, arg := range args {
			if arg == "--port" {
				args[idx+1] = strconv.Itoa(rand.Intn(65535-1024) + 1024)
			}
		}

		_, err = deploymentsClient.Update(context.Background(), deployment, metav1.UpdateOptions{})
		return err
	})

	if err != nil {
		panic(err)
	}
}
