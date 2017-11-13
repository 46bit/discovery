package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/namespaces"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	// command := os.Args[1]

	// switch command {
	// case "run":
	// 	containerName := os.Args[2]
	// 	imageReference := "docker.io/46bit/" + containerName + ":latest"

	// }

	fmt.Println("running")
	containerName := os.Args[1]
	imageReference := os.Args[2]
	if err := redisExample(containerName, imageReference); err != nil {
		log.Fatal(err)
	}
}

func redisExample(containerName, imageReference string) error {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "default")

	// pull the redis image from DockerHub
	image, err := client.Pull(ctx, imageReference, containerd.WithPullUnpack)
	if err != nil {
		return fmt.Errorf("Error pulling image: %s", err.Error())
	}

	// create a container
	snapshotViewName := fmt.Sprintf("%s-snapshotView-%d", containerName, rand.Uint64())
	container, err := client.NewContainer(
		ctx,
		containerName,
		containerd.WithImage(image),
		containerd.WithNewSnapshotView(snapshotViewName, image),
		containerd.WithNewSpec(containerd.WithImageConfig(image), withHostNetworkNamespace),
	)
	if err != nil {
		return fmt.Errorf("Error creating new container: %s", err.Error())
	}

	// create a task from the container
	task, err := container.NewTask(ctx, containerd.Stdio)
	if err != nil {
		return fmt.Errorf("Error creating new task: %s", err.Error())
	}

	// make sure we wait before calling start
	_, err = task.Wait(ctx)
	if err != nil {
		fmt.Printf("Ignored error when waiting for task: %s\n", err)
	}

	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		return fmt.Errorf("Error starting task: %s", err.Error())
	}

	time.Sleep(10 * time.Second)

	return nil
}
