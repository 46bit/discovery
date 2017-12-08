package main

import (
	"github.com/46bit/discovery/rainbow/executor"
	"github.com/46bit/discovery/rainbow/operator"
	cd "github.com/containerd/containerd"
	"log"
	"time"
)

func main() {
	client, err := cd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	exec := executor.NewExecutor(client)
	go exec.Run()

	op := operator.NewOperator(exec.CmdChan, exec.EventChan)
	go op.Run()

	time.Sleep(2 * time.Second)
}
