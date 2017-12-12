package main

import (
	"github.com/46bit/discovery/rainbow"
	"github.com/46bit/discovery/rainbow/executor"
	"github.com/46bit/discovery/rainbow/operator"
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

	exec := executor.NewExecutor(client)
	go exec.Run()

	op := operator.NewOperator(exec.CmdChan, exec.EventChan)
	go op.Run()

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
	op.Add(helloWorld)
	time.Sleep(time.Minute)

	op.Remove(helloWorld.Name)
	time.Sleep(10 * time.Second)
}
