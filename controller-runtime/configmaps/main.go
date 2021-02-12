package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"k8s.io/helm/pkg/strvals"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/helm/pkg/chartutil"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

const (
	priorityKey = "kubeaddons.d2iq.com/priority"
)

func main() {

	c, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		panic(err)
	}

	labelSelector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app.kubernetes.io/name":      "kubeaddons",
			"app.kubernetes.io/component": "controller",
		},
	})
	if err != nil {
		panic(err)
	}

	// fetch all ConfigMaps that match the Addon's desired label selectors
	cmList := corev1.ConfigMapList{}
	if err := c.List(context.Background(), &cmList, client.MatchingLabelsSelector{
		Selector: labelSelector,
	}); err != nil {
		panic(err)
	}

	// sort matching ConfigMaps by priority in ascending order
	sort.Slice(cmList.Items, func(i, j int) bool {
		prioi, err := strconv.Atoi(cmList.Items[i].Annotations[priorityKey])
		if err != nil {
			fmt.Printf("Could not convert: %s\n", err)
			return false
		}
		prioj, err := strconv.Atoi(cmList.Items[j].Annotations[priorityKey])
		if err != nil {
			fmt.Printf("Could not convert: %s\n", err)
			return false
		}
		return prioi < prioj
	})

	// merge values from all matching sorted ConfigMaps
	values := chartutil.Values{}
	for _, cm := range cmList.Items {
		cmValues, err := chartutil.ReadValues([]byte(cm.Data["values"]))
		if err != nil {
			panic(err)
		}
		values.MergeInto(cmValues)
	}
	fmt.Printf("merged values: %s\n", values)

	// this will later come from the Addon spec
	remapValues := map[string]string{
		"image.tag": "heartbeat.image.version",
		"image.repo": "heartbeat.image.repo",
	}
	remappedValues := make([]string, 0)
	for k, v := range remapValues {
		storeVal, err := values.PathValue(v)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not remap value %s from %s: %s\n", k, v, err)
			continue
		}
		remappedValues = append(remappedValues, fmt.Sprintf("%s=%v", k, storeVal))
	}

	remappedParameters := make(map[string]string, 0)
	for _, v := range remappedValues {
		kv := strings.Split(v, "=")
		remappedParameters[kv[0]] = kv[1]
	}
	fmt.Printf("=== Remapped parameters ===\n\n%s\n", remappedParameters)

	overrideValues, err := strvals.ToYAML(strings.Join(remappedValues, ","))
	if err != nil {
		panic(err)
	}
	fmt.Printf("=== Remapped values ===\n\n%s\n", overrideValues)
}
