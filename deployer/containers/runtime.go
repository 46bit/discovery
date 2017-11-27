package containers

import (
	"context"
	"errors"
	cd "github.com/containerd/containerd"
	ns "github.com/containerd/containerd/namespaces"
	"log"
	"time"
)

type ContainerDesc struct {
	ID        string
	Remote    string
	Namespace string
}

type Runtime struct {
	Add        chan ContainerDesc
	Remove     chan string
	client     *cd.Client
	namespaces map[string]context.Context
	containers map[string]*container
}

func NewRuntime(client *cd.Client) *Runtime {
	return &Runtime{
		Add:        make(chan ContainerDesc),
		Remove:     make(chan string),
		client:     client,
		namespaces: map[string]context.Context{},
		containers: map[string]*container{},
	}
}

func (r *Runtime) Run() {
	for {
		select {
		case containerDesc := <-r.Add:
			err := r.runAdd(containerDesc)
			if err != nil {
				log.Printf("Add error for container %s: %s", containerDesc.ID, err)
			}
		case containerID := <-r.Remove:
			err := r.runRemove(containerID)
			if err != nil {
				log.Printf("Remove error for container %s: %s", containerID, err)
			}
		case <-time.After(time.Second):
			err := r.runHealthchecks()
			if err != nil {
				log.Printf("Healthcheck error: %s", err)
			}
		}
	}
}

func (r *Runtime) runAdd(containerDesc ContainerDesc) error {
	namespace := r.namespace(containerDesc.Namespace)
	api := cdApi{client: r.client, context: namespace}

	container, err := newContainer(api, containerDesc)
	if err != nil {
		return err
	}
	r.containers[containerDesc.ID] = container

	container.task, err = newTask(api, container.container)
	if err != nil {
		return err
	}
	err = container.task.start(api)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runtime) runRemove(containerID string) error {
	container, ok := r.containers[containerID]
	if !ok {
		return errors.New("Container not found to delete")
	}
	namespace := r.namespace(container.desc.Namespace)
	api := cdApi{client: r.client, context: namespace}

	err := container.task.stop(api)
	if err != nil {
		return err
	}
	err = container.task.delete(api)
	if err != nil {
		return err
	}
	container.task = nil

	return nil
}

func (r *Runtime) runHealthchecks() error {
	// Poll all containers for status
	return nil
}

func (r *Runtime) namespace(namespaceName string) context.Context {
	if _, ok := r.namespaces[namespaceName]; !ok {
		r.namespaces[namespaceName] = ns.WithNamespace(context.Background(), namespaceName)
	}
	return r.namespaces[namespaceName]
}
