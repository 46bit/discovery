package main

import (
	"github.com/46bit/discovery/rainbow"
	"log"
	"time"
)

func main() {
	client := rainbow.NewClient("http://localhost:8080")

	serviceDiscovery, err := client.Create(rainbow.Deployment{
		Name: "service-discovery",
		Jobs: []rainbow.Job{
			{
				Name:          "discoverer",
				Remote:        "docker.io/46bit/discoverer:latest",
				InstanceCount: 1,
			},
		},
	})
	if err != nil {
		log.Println(err)
	}
	time.Sleep(5 * time.Second)

	sendersReceiver, err := client.Create(rainbow.Deployment{
		Name: "senders-receivers-aggregator",
		Jobs: []rainbow.Job{
			{
				Name:          "aggregator",
				Remote:        "docker.io/46bit/aggregator:latest",
				InstanceCount: 1,
			},
			{
				Name:          "receiver",
				Remote:        "docker.io/46bit/receiver:latest",
				InstanceCount: 2,
			},
			{
				Name:          "sender",
				Remote:        "docker.io/46bit/sender:latest",
				InstanceCount: 4,
			},
		},
	})
	if err != nil {
		log.Println(err)
	}
	time.Sleep(time.Minute)

	if err := client.Delete(sendersReceiver.Name); err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)
	if err := client.Delete(serviceDiscovery.Name); err != nil {
		log.Println(err)
	}
	time.Sleep(5 * time.Second)
}
