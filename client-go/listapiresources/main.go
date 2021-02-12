package main

import (
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/tools/clientcmd"
)

func versions(gvs []v1.GroupVersionForDiscovery, pref v1.GroupVersionForDiscovery) string {
	res := make([]string, len(gvs))
	for idx, gv := range gvs {
		res[idx] = "    " + gv.GroupVersion
		if gv.GroupVersion == pref.GroupVersion {
			res[idx] += "*"
		}
	}
	return strings.Join(res, "\n")
}

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}
	discovery, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err)
	}
	resourceList, err := discovery.ServerPreferredResources()
	if err != nil {
		panic(err)
	}

	for _, apiResourceList := range resourceList {
		for _, apiResource := range apiResourceList.APIResources {
			fmt.Printf("%s.%s\n", apiResourceList.GroupVersion, apiResource.Kind)
		}
	}
}
