package main

import (
	"github.com/46bit/discovery/rainbow"
	"github.com/46bit/discovery/rainbow/executor"
	cd "github.com/containerd/containerd"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	client, err := cd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	runtime := executor.NewRuntime(client)
	go runtime.Run()

	depl := rainbow.NewDeployer(runtime)
	go depl.Run()

	serviceDiscovery := rainbow.Deployment{
		Name: "hello-world",
		Jobs: []rainbow.Job{
			{
				Name:      "hello-world",
				Remote:    "docker.io/46bit/hello-world:latest",
				Instances: 1,
			},
		},
	}
	depl.Add <- serviceDiscovery
	time.Sleep(time.Minute)

	depl.Remove <- serviceDiscovery.Name
	time.Sleep(10 * time.Second)
}
