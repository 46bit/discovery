package main

import (
	"github.com/46bit/discovery/rainbow"
	"github.com/46bit/discovery/rainbow/executor"
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

	runtime := executor.NewRuntime(client)
	go runtime.Run()

	depl := rainbow.NewDeployer(runtime)
	go depl.Run()

	time.Sleep(2 * time.Second)
}
