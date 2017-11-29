package main

import (
	"github.com/46bit/discovery/rainbow/containers"
	"github.com/46bit/discovery/rainbow/deployments"
	cd "github.com/containerd/containerd"
	"log"
	"time"
)

func main() {
	client, err := cd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	runtime := containers.NewRuntime(client)
	go runtime.Run()

	depl := deployments.NewDeployer(runtime)
	go depl.Run()

	time.Sleep(2 * time.Second)
}
