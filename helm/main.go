package main

import (
	"fmt"
	"os"

	helmlib "k8s.io/helm/pkg/helm"
	hapirelease "k8s.io/helm/pkg/proto/hapi/release"
)

var (
	allCodes = []hapirelease.Status_Code{
		hapirelease.Status_UNKNOWN,
		hapirelease.Status_DEPLOYED,
		hapirelease.Status_DELETED,
		hapirelease.Status_SUPERSEDED,
		hapirelease.Status_FAILED,
		hapirelease.Status_DELETING,
		hapirelease.Status_PENDING_INSTALL,
		hapirelease.Status_PENDING_UPGRADE,
		hapirelease.Status_PENDING_ROLLBACK,
	}
)

func main() {

	release := os.Args[1]

	client := helmlib.NewClient(helmlib.Host("localhost:44134"))
	rels, err := client.ListReleases(helmlib.ReleaseListFilter(release), helmlib.ReleaseListStatuses(allCodes))
	if err != nil {
		panic(err)
	}

	// b, err := json.Marshal(rels.Releases)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", b)

	// rel := rels.Releases[rels.GetCount()-1]
	for _, rel := range rels.Releases {
		fmt.Printf("%s: %s\n", rel.Chart.Metadata.Name, rel.Chart.Metadata.Version)
	}
}
