package container

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"sync"
	"syscall"
	"time"
)

const (
	namespace = "default"
)

type State uint

const (
	Described State = iota
	Created
	Tasked
	Started
	Stopped
	Untasked
	Deleted
)

type Container struct {
	ID        string
	Remote    string
	State     State
	container *containerd.Container
	task      *containerd.Task
	sync.Mutex
}

func NewContainer(id, remote string) *Container {
	return &Container{
		ID:     id,
		Remote: remote,
		State:  Described,
	}
}

func (c *Container) Create(client *containerd.Client) error {
	c.Lock()
	defer c.Unlock()
	ctx := c.context()
	image, err := client.Pull(ctx, c.Remote, containerd.WithPullUnpack)
	if err != nil {
		return errors.Wrap(err, "Error pulling image")
	}
	imageConfig := containerd.WithImageConfig(image)
	withHostNamespace := containerd.WithHostNamespace(specs.NetworkNamespace)
	container, err := client.NewContainer(
		ctx,
		c.ID,
		containerd.WithImage(image),
		containerd.WithNewSnapshot("snapshot-"+c.ID, image),
		containerd.WithNewSpec(imageConfig, withHostNamespace),
	)
	if err != nil {
		return errors.Wrap(err, "Error creating new container")
	}
	c.State = Created
	c.container = &container
	return nil
}

func (c *Container) Task() error {
	c.Lock()
	defer c.Unlock()
	task, err := (*c.container).NewTask(c.context(), containerd.Stdio)
	if err != nil {
		return errors.Wrap(err, "Error creating containerd task")
	}
	c.State = Tasked
	c.task = &task
	return nil
}

func (c *Container) Start() (<-chan containerd.ExitStatus, error) {
	c.Lock()
	defer c.Unlock()
	ctx := c.context()
	exitStatusC, err := (*c.task).Wait(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error waiting for containerd task")
	}
	err = (*c.task).Start(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Error starting containerd task")
	}
	c.State = Started
	return exitStatusC, nil
}

func (c *Container) Stop() error {
	c.Lock()
	defer c.Unlock()
	ctx := c.context()
	exitStatusC, err := (*c.task).Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "Error waiting for containerd task")
	}
	err = (*c.task).Kill(ctx, syscall.SIGTERM, containerd.WithKillAll)
	if err != nil {
		return errors.Wrap(err, "Error SIGTERMing containerd task")
	}
	select {
	case <-exitStatusC:
	case <-time.After(5 * time.Second):
		err = (*c.task).Kill(ctx, syscall.SIGKILL, containerd.WithKillAll)
		if err != nil {
			return errors.Wrap(err, "Error SIGKILLing containerd task")
		}
		<-exitStatusC
	}
	c.State = Stopped
	return nil
}

func (c *Container) Untask() error {
	c.Lock()
	defer c.Unlock()
	_, err := (*c.task).Delete(c.context())
	if err != nil {
		return errors.Wrap(err, "Error deleting containerd task")
	}
	c.State = Untasked
	c.task = nil
	return nil
}

func (c *Container) Delete() error {
	c.Lock()
	defer c.Unlock()
	err := (*c.container).Delete(c.context(), containerd.WithSnapshotCleanup)
	if err != nil {
		return errors.Wrap(err, "Error deleting container")
	}
	c.State = Deleted
	c.container = nil
	return nil
}

func (c *Container) Status() State {
	c.Lock()
	defer c.Unlock()
	if c.State == Started {
		status, err := (*c.task).Status(c.context())
		if err == nil && status.Status != containerd.Running {
			c.State = Stopped
		}
	}
	return c.State
}

func (c *Container) context() context.Context {
	return namespaces.WithNamespace(context.Background(), namespace)
}
