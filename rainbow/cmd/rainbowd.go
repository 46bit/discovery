package main

import (
	"github.com/46bit/discovery/rainbow/server"
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

	s := server.NewServer(client)
	s.Run()
}
