package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/containerd/containerd"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	newGroups := make(chan Group)
	deleteGroups := make(chan string)
	e := NewExecutor("default", client, newGroups, deleteGroups)
	go func(e *Executor) {
		e.run()
	}(e)

	groupA := NewGroup("a", []string{"docker.io/46bit/hello-world:latest", "docker.io/46bit/long-running:latest"})
	newGroups <- groupA
	time.Sleep(1 * time.Second)

	groupB := NewGroup("b", []string{"docker.io/46bit/sender:latest", "docker.io/46bit/receiver:latest"})
	newGroups <- groupB
	time.Sleep(10 * time.Second)

	deleteGroups <- groupA.Name
	time.Sleep(10 * time.Second)

	deleteGroups <- groupB.Name
	time.Sleep(10 * time.Second)
}
