package runtime

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"log"
)

type Runtime struct {
	Client     *containerd.Client
	Namespaces map[string]context.Context
	Containers map[string]*Container
	Add        chan Container
	Remove     chan string
	exit       chan taskExit
}

func NewRuntime(client *containerd.Client) *Runtime {
	return &Runtime{
		Client:     client,
		Namespaces: map[string]context.Context{},
		Containers: map[string]*Container{},
		Add:        make(chan Container),
		Remove:     make(chan string),
		exit:       make(chan taskExit),
	}
}

func (r *Runtime) Run() {
	for {
		select {
		case container := <-r.Add:
			r.runAdd(container)
		case containerID := <-r.Remove:
			r.runRemove(containerID)
		case taskExit := <-r.exit:
			r.runExit(taskExit.containerID, taskExit.exitCode)
		}
	}
}

func (r *Runtime) runAdd(container Container) {
	if container.State != Described {
		log.Println(fmt.Sprintf("Added container was not in Described state: %s", container))
		return
	}
	r.Containers[container.ID] = &container

	namespace, ok := r.Namespaces[container.Namespace]
	if !ok {
		namespace := namespaces.WithNamespace(context.Background(), container.Namespace)
		r.Namespaces[container.Namespace] = namespace
	}

	_, err := container.Create(r.Client, namespace)
	if err != nil {
		log.Println(fmt.Sprintf("Error creating container %s: %s", container.ID, err))
		return
	}
	_, err = container.Prepare(r.Client, namespace)
	if err != nil {
		container.Erase(r.Client, namespace)
		log.Println(fmt.Sprintf("Error preparing container %s: %s", container.ID, err))
		return
	}
	exitStatusC, err := container.Start(r.Client, namespace)
	if err != nil {
		container.Delete(r.Client, namespace)
		container.Erase(r.Client, namespace)
		log.Println(fmt.Sprintf("Error preparing container %s: %s", container.ID, err))
		return
	}
	go func(containerID string, exit chan taskExit, exitStatusC <-chan containerd.ExitStatus) {
		exitStatus := <-exitStatusC
		exit <- taskExit{
			containerID: containerID,
			exitCode:    exitStatus,
		}
	}(container.ID, r.exit, exitStatusC)
}

func (r *Runtime) runRemove(containerID string) {
	container := r.Containers[containerID]
	namespace := r.Namespaces[container.Namespace]

	err := container.Sigkill(r.Client, namespace)
	if err != nil {
		log.Println(fmt.Sprintf("Error killing container %s: %s", containerID, err))
		return
	}
}

func (r *Runtime) runExit(containerID string, _ containerd.ExitStatus) {
	container := r.Containers[containerID]
	namespace := r.Namespaces[container.Namespace]

	container.State = Exited

	err := container.Delete(r.Client, namespace)
	if err != nil {
		log.Println(fmt.Sprintf("Error deleting container %s: %s", containerID, err))
		return
	}
	_, err = container.Prepare(r.Client, namespace)
	if err != nil {
		container.Erase(r.Client, namespace)
		log.Println(fmt.Sprintf("Error preparing container %s: %s", container.ID, err))
		return
	}
	exitStatusC, err := container.Start(r.Client, namespace)
	if err != nil {
		container.Delete(r.Client, namespace)
		container.Erase(r.Client, namespace)
		log.Println(fmt.Sprintf("Error preparing container %s: %s", container.ID, err))
		return
	}
	go func(containerID string, exit chan taskExit, exitStatusC <-chan containerd.ExitStatus) {
		exitStatus := <-exitStatusC
		exit <- taskExit{
			containerID: containerID,
			exitCode:    exitStatus,
		}
	}(container.ID, r.exit, exitStatusC)
}

type taskExit struct {
	containerID string
	exitCode    containerd.ExitStatus
}
