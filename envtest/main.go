package main

import (
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func main() {
	env := envtest.Environment{}
	cfg, err := env.Start()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", cfg)
}
