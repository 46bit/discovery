package main

import (
	"github.com/46bit/discovery/rainbow"
	"github.com/46bit/discovery/rainbow/client"
	"log"
	"time"
)

func main() {
	client := client.NewClient("http://localhost:8080")

	helloWorld, err := client.Create(rainbow.Deployment{
		Name: "hello-world",
		Jobs: []rainbow.Job{
			{
				Name:          "hello-world",
				Remote:        "docker.io/46bit/hello-world:latest",
				InstanceCount: 1,
			},
		},
	})
	if err != nil {
		log.Println(err)
	}
	time.Sleep(time.Minute)

	if err := client.Delete(helloWorld.Name); err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)
}
