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

	depl := deployer.NewDeployer(runtime)
	go depl.Run()

	serviceDiscovery := deployer.Deployment{
		Name: "service-discovery",
		Jobs: []deployer.Job{
			{
				Name:      "discoverer",
				Remote:    "docker.io/46bit/discoverer:latest",
				Instances: 1,
			},
		},
	}
	depl.Add <- serviceDiscovery
	time.Sleep(5 * time.Second)

	sendersReceiver := deployer.Deployment{
		Name: "senders-receiver",
		Jobs: []deployer.Job{
			deployer.Job{
				Name:      "aggregator",
				Remote:    "docker.io/46bit/aggregator:latest",
				Instances: 1,
			},
			deployer.Job{
				Name:      "receiver",
				Remote:    "docker.io/46bit/receiver:latest",
				Instances: 7,
			},
			deployer.Job{
				Name:      "sender",
				Remote:    "docker.io/46bit/sender:latest",
				Instances: 28,
			},
		},
	}
	depl.Add <- sendersReceiver
	time.Sleep(4 * time.Minute)

	depl.Remove <- sendersReceiver.Name
	time.Sleep(5 * time.Second)

	depl.Remove <- serviceDiscovery.Name
	time.Sleep(5 * time.Second)
}
