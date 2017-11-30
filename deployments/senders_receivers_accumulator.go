package main

import (
	"fmt"
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
		Name: "service-discovery",
		Jobs: []rainbow.Job{
			{
				Name:      "discoverer",
				Remote:    "docker.io/46bit/discoverer:latest",
				Instances: 1,
			},
		},
	}
	depl.Add <- serviceDiscovery
	time.Sleep(5 * time.Second)

	for i := uint(1); i <= 3; i++ {
		for j := uint(1); j <= 7; j++ {
			log.Printf("------\nSENDERS-RECEIVER SET WITH %d, %d\n------\n", i, j)

			sendersReceiver := rainbow.Deployment{
				Name: fmt.Sprintf("senders-receiver-i%d-j%d", i, j),
				Jobs: []rainbow.Job{
					{
						Name:      "aggregator",
						Remote:    "docker.io/46bit/aggregator:latest",
						Instances: 1,
					},
					{
						Name:      "receiver",
						Remote:    "docker.io/46bit/receiver:latest",
						Instances: i,
					},
					{
						Name:      "sender",
						Remote:    "docker.io/46bit/sender:latest",
						Instances: i * j,
					},
				},
			}
			depl.Add <- sendersReceiver
			time.Sleep(time.Duration(int64(i)) * time.Minute)

			log.Printf("------\n")
			depl.Remove <- sendersReceiver.Name
			time.Sleep(30 * time.Second)
		}
	}

	depl.Remove <- serviceDiscovery.Name
	time.Sleep(5 * time.Second)
}