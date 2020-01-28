package main

import (
	"fmt"
	"os"
	"strings"

	versionutil "k8s.io/apimachinery/pkg/util/version"

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
	version, err := discovery.ServerVersion()
	if err != nil {
		panic(err)
	}
	serverVersion := versionutil.MustParseSemantic(version.String())
	fmt.Printf("The cluster is running Kubernets %s\n", serverVersion)
}
