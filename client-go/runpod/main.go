package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	podClient := clientset.CoreV1().Pods("default")
	pod, err := podClient.Create(context.Background(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "debug",
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{{
				Name:  "debug",
				Image: "digitalocean/doks-debug:latest",
				Args: []string{
					"curl",
					"-sk",
					"https://kubernetes/version",
				},
			}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	defer func() {
		podClient.Delete(context.Background(), "debug", metav1.DeleteOptions{})
	}()

	wait.Poll(1*time.Second, 60*time.Second, func() (bool, error) {
		var err error
		pod, err = podClient.Get(context.Background(), "debug", metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return pod.Status.Phase == corev1.PodSucceeded, nil
	})

	res, err := podClient.GetLogs("debug", &corev1.PodLogOptions{}).Do(context.Background()).Raw()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", res)
}
