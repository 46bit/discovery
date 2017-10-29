package main

import (
  "os"
  "context"
  "fmt"
  "log"
  "syscall"
  "time"

  "github.com/containerd/containerd"
  "github.com/containerd/containerd/namespaces"
  "github.com/containerd/containerd/containers"
  specs "github.com/opencontainers/runtime-spec/specs-go"
)

func main() {
  fmt.Println("running")
  containerName := os.Args[1]
  imageReference := os.Args[2]
  if err := redisExample(containerName, imageReference); err != nil {
    log.Fatal(err)
  }
}

// adapted from WithHtop https://github.com/containerd/containerd/blob/a6ce1ef2a140d79856a8647e1d1ae5ac9ab581eb/docs/client-opts.md
func withHostNetworkNamespace(context context.Context, client *containerd.Client, container *containers.Container, s *specs.Spec) error {
  // make sure we are in the host network namespace
  if err := containerd.WithHostNamespace(specs.NetworkNamespace)(context, client, container, s); err != nil {
    return err
  }
  return nil
}

func redisExample(containerName, imageReference string) error {
  // create a new client connected to the default socket path for containerd
  client, err := containerd.New("/run/containerd/containerd.sock")
  if err != nil {
    return err
  }
  defer client.Close()

  // create a new context with an "example" namespace
  ctx := namespaces.WithNamespace(context.Background(), "example")

  // pull the redis image from DockerHub
  image, err := client.Pull(ctx, imageReference, containerd.WithPullUnpack)
  if err != nil {
    return err
  }

  // create a container
  container, err := client.NewContainer(
    ctx,
    containerName,
    containerd.WithImage(image),
    containerd.WithNewSnapshot(containerName+"-snapshot", image),
    containerd.WithNewSpec(containerd.WithImageConfig(image), withHostNetworkNamespace),
  )
  if err != nil {
    return err
  }
  defer container.Delete(ctx, containerd.WithSnapshotCleanup)

  // create a task from the container
  task, err := container.NewTask(ctx, containerd.Stdio)
  if err != nil {
    return err
  }
  defer task.Delete(ctx)

  // make sure we wait before calling start
  exitStatusC, err := task.Wait(ctx)
  if err != nil {
    fmt.Println(err)
  }

  // call start on the task to execute the redis server
  if err := task.Start(ctx); err != nil {
    return err
  }

  // sleep for a lil bit to see the logs
  time.Sleep(300 * time.Second)

  // kill the process and get the exit status
  if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
    return err
  }

  // wait for the process to fully exit and print out the exit status

  status := <-exitStatusC
  code, _, err := status.Result()
  if err != nil {
    return err
  }
  fmt.Printf(containerName+" exited with status: %d\n", code)

  return nil
}
