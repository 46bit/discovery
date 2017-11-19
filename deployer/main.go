package main

import (
	"github.com/46bit/discovery/deployer/deployer"
	"github.com/46bit/discovery/deployer/runtime"
	"github.com/containerd/containerd"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	runtime := runtime.NewRuntime(client)
	go runtime.Run()

	deployment := deployer.Deployment{
		Name: "senders-receiver",
		Jobs: []deployer.Job{
			deployer.Job{
				Name:      "sender",
				Remote:    "docker.io/46bit/sender:latest",
				Instances: 10,
			},
			deployer.Job{
				Name:      "receiver",
				Remote:    "docker.io/46bit/receiver:latest",
				Instances: 1,
			},
		},
	}

	deployer := deployer.NewDeployer(runtime)
	go deployer.Run()

	deployer.Add <- deployment
	time.Sleep(time.Minute)
	deployer.Remove <- deployment.Name

	for {
		time.Sleep(10 * time.Second)
	}
}
