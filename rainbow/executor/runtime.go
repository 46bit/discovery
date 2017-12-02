package executor

import (
	"context"
	cd "github.com/containerd/containerd"
	ns "github.com/containerd/containerd/namespaces"
	"github.com/pkg/errors"
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
		case <-time.After(100 * time.Millisecond):
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
		return errors.Wrap(err, "Error creating new container")
	}
	r.containers[containerDesc.ID] = container
	return nil
}

func (r *Runtime) runRemove(containerID string) error {
	container, ok := r.containers[containerID]
	if !ok {
		return errors.New("Container not found to remove")
	}
	delete(r.containers, containerID)

	namespace := r.namespace(container.desc.Namespace)
	api := cdApi{client: r.client, context: namespace}
	if container.task.state == started {
		container.task.stop(api)
	}
	if container.task.state == created || container.task.state == stopped {
		container.task.delete(api)
	}
	if container.state == created {
		container.delete(api)
	}

	return nil
}

func (r *Runtime) runHealthchecks() error {
	for _, container := range r.containers {
		namespace := r.namespace(container.desc.Namespace)
		if container.task.state == started {
			status, err := (*container.task.task).Status(namespace)
			if err != nil {
				return err
			}
			if status.Status != cd.Created && status.Status != cd.Running {
				container.task.state = stopped
			}
		}

		api := cdApi{client: r.client, context: namespace}
		err := transition(api, container)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Runtime) namespace(namespaceName string) context.Context {
	if _, ok := r.namespaces[namespaceName]; !ok {
		r.namespaces[namespaceName] = ns.WithNamespace(context.Background(), namespaceName)
	}
	return r.namespaces[namespaceName]
}
