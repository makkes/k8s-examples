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
	groups, err := discovery.ServerGroups()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s/%s\n", groups.Kind, groups.APIVersion)
	for idx, group := range groups.Groups {
		fmt.Printf("%s:\n%s\n", group.Name, versions(group.Versions, group.PreferredVersion))
		if idx < len(groups.Groups)-1 {
			fmt.Printf("\n")
		}
	}
}
