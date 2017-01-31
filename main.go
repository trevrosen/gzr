package main

import (
	"github.com/bypasslane/gzr/cmd"
	"github.com/bypasslane/gzr/comms"
)

func main() {
	err := comms.EstablishK8sConnection("/Users/trevor/.kube/config")

	if err != nil {
		panic(err)
	}

	cmd.Execute()
}
