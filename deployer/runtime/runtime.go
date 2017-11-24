package runtime

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"log"
	"syscall"
)

type Runtime struct {
	Client     *containerd.Client
	Namespaces map[string]context.Context
	Containers map[string]Container
	Tasks      map[string]taskRuntime
	Add        chan Container
	Remove     chan string
	Exit       chan taskExit
}

func NewRuntime(client *containerd.Client) *Runtime {
	return &Runtime{
		Client:     client,
		Namespaces: map[string]context.Context{},
		Containers: map[string]Container{},
		Tasks:      map[string]taskRuntime{},
		Add:        make(chan Container),
		Remove:     make(chan string),
		Exit:       make(chan taskExit),
	}
}

func (r *Runtime) Run() {
	for {
		select {
		case container := <-r.Add:
			if _, ok := r.Namespaces[container.Namespace]; !ok {
				r.Namespaces[container.Namespace] = namespaces.WithNamespace(context.Background(), container.Namespace)
			}
			r.Containers[container.ID] = container
			err := r.run(container.ID)
			if err != nil {
				log.Println(fmt.Sprintf("Error running container %s: %s", container.ID, err))
			}
		case containerID := <-r.Remove:
			err := r.kill(containerID, syscall.SIGTERM)
			if err != nil {
				log.Println(fmt.Sprintf("Error killing container %s: %s", containerID, err))
			}
			delete(r.Containers, containerID)
		case taskExit := <-r.Exit:
			err := r.delete(taskExit.ContainerID)
			if err != nil {
				log.Println(fmt.Sprintf("Error deleting container %s: %s", taskExit.ContainerID, err))
			}
			_, ok := r.Containers[taskExit.ContainerID]
			if ok {
				err = r.run(taskExit.ContainerID)
				if err != nil {
					log.Println(fmt.Sprintf("Error re-running container %s: %s", taskExit.ContainerID, err))
				}
			}
		}
	}
}

func (r *Runtime) run(id string) error {
	container := r.Containers[id]
	ctx := r.Namespaces[container.Namespace]
	containerdContainer, err := createContainer(r.Client, ctx, container.ID, container.Remote)
	if err != nil {
		return err
	}
	containerdTask, exitStatusC, err := createTask(ctx, containerdContainer)
	if err != nil {
		containerdContainer.Delete(ctx, containerd.WithSnapshotCleanup)
		return err
	}
	r.Tasks[container.ID] = taskRuntime{
		Task:      containerdTask,
		Namespace: container.Namespace,
	}
	go func(containerID string, exit chan taskExit, exitStatusC <-chan containerd.ExitStatus) {
		exitStatus := <-exitStatusC
		exit <- taskExit{
			ContainerID: containerID,
			ExitCode:    exitStatus,
		}
	}(container.ID, r.Exit, exitStatusC)
	return nil
}

func (r *Runtime) kill(id string, signal syscall.Signal) error {
	taskRuntime, ok := r.Tasks[id]
	if !ok {
		return nil
	}
	ctx := r.Namespaces[taskRuntime.Namespace]
	status, err := taskRuntime.Task.Status(ctx)
	if err != nil {
		return err
	}
	switch status.Status {
	case containerd.Running:
		fallthrough
	case containerd.Paused:
		fallthrough
	case containerd.Pausing:
		err = taskRuntime.Task.Kill(ctx, signal, containerd.WithKillAll)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Runtime) delete(id string) error {
	taskRuntime := r.Tasks[id]
	ctx := r.Namespaces[taskRuntime.Namespace]
	_, err := taskRuntime.Task.Delete(ctx)
	if err != nil {
		return fmt.Errorf("Error deleting task %s: %s", id, err)
	}
	delete(r.Tasks, id)

	containerdContainer, err := r.Client.LoadContainer(ctx, id)
	if err != nil {
		return fmt.Errorf("Error loading container %s: %s", id, err)
	}

	err = containerdContainer.Delete(ctx, containerd.WithSnapshotCleanup)
	if err != nil {
		return fmt.Errorf("Error deleting container %s: %s", id, err)
	}

	return nil
}

type taskRuntime struct {
	Task      containerd.Task
	Namespace string
}

type taskExit struct {
	ContainerID string
	ExitCode    containerd.ExitStatus
}

type Container struct {
	ID        string
	Remote    string
	Namespace string
}
