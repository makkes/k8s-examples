package main

import (
	"path"
	"context"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func curl(c typedcorev1.PodInterface, url string) (string, error) {
	pod, err := c.Create(context.Background(), &corev1.Pod{
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
					url,
				},
			}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("could not create pod: %w", err)
	}

	defer func() {
		c.Delete(context.Background(), "debug", metav1.DeleteOptions{})
	}()

	if err := wait.Poll(1*time.Second, 60*time.Second, func() (bool, error) {
		var err error
		pod, err = c.Get(context.Background(), "debug", metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return pod.Status.Phase == corev1.PodSucceeded, nil
	})
	err != nil{
		return "", fmt.Errorf("could not verify pod is running: %w", err)
	}

	res, err := c.GetLogs("debug", &corev1.PodLogOptions{}).Do(context.Background()).Raw()
	if err != nil {
		return "", fmt.Errorf("could not get pod logs: %w", err)
	}
	return string(res), nil
}

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		kubeconfig = path.Join(homeDir, ".kube", "config")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset := kubernetes.NewForConfigOrDie(config)
	podClient := clientset.CoreV1().Pods("default")
	res, err := curl(podClient, "https://kubernetes/version")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", res)
}
