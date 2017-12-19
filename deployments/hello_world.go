package main

import (
	"github.com/46bit/discovery/rainbow"
	"log"
	"time"
)

func main() {
	client := rainbow.NewClient("http://localhost:4601")

	helloWorld := rainbow.Deployment{
		Name: "hello-world",
		Jobs: []rainbow.Job{
			{
				Name:          "hello-world",
				Remote:        "docker.io/46bit/hello-world:latest",
				InstanceCount: 1,
			},
		},
	}
	if err := client.CreateDeployment(helloWorld); err != nil {
		log.Println(err)
	}
	time.Sleep(time.Minute)

	if err := client.DeleteDeployment(helloWorld.Name); err != nil {
		log.Println(err)
	}
	time.Sleep(10 * time.Second)
}
