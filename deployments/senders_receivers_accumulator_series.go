package main

import (
	"fmt"
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

	for i := uint(1); i <= 3; i++ {
		for j := uint(1); j <= 7; j++ {
			log.Printf("------\nSENDERS-RECEIVER SET WITH %d, %d\n------\n", i, j)

			sendersReceiver, err := client.Create(rainbow.Deployment{
				Name: fmt.Sprintf("senders-receiver-i%d-j%d", i, j),
				Jobs: []rainbow.Job{
					{
						Name:          "aggregator",
						Remote:        "docker.io/46bit/aggregator:latest",
						InstanceCount: 1,
					},
					{
						Name:          "receiver",
						Remote:        "docker.io/46bit/receiver:latest",
						InstanceCount: i,
					},
					{
						Name:          "sender",
						Remote:        "docker.io/46bit/sender:latest",
						InstanceCount: i * j,
					},
				},
			})
			if err != nil {
				log.Println(err)
			}
			time.Sleep(time.Duration(int64(i)) * time.Minute)

			log.Printf("------\n")
			if err := client.Delete(sendersReceiver.Name); err != nil {
				log.Println(err)
			}
			time.Sleep(30 * time.Second)
		}
	}

	if err := client.Delete(serviceDiscovery.Name); err != nil {
		log.Println(err)
	}
	time.Sleep(5 * time.Second)
}
