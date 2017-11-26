package runtime

import (
	"context"
	"github.com/containerd/containerd"
	"syscall"
)

const (
	// An unused description of a Container
	Described ContainerState = iota
	// A Container whose containerd.Container has been created
	Created
	// A Container whose containerd.Task has been created
	Prepared
	// A Container whose containerd.Task has been started
	Started
	// A Container that has been SIGTERMed
	Sigtermed
	// A Container that has been SIGKILLed
	Sigkilled
	// A Container that has exited
	Exited
	// A Container whose containerd.Task has been deleted
	Deleted
	// A Container whose containerd.Container has been deleted
	Erased
)

type ContainerState uint

// RESTARTING BUT WITH A BRAND-NEW CONTAINER AND TASK:
//   Created -> Prepared -> Started -> Stopped -> Deleted -> Erased -> Created -> ...
// STOPPING:
//   Created -> Prepared -> Started -> Sigtermed -> Stopped -> Deleted -> Erased
// FORCE-STOPPING AFTER TIMEOUT:
//   Created -> Prepared -> Started -> Sigtermed -> Sigkilled -> Stopped -> Deleted -> Erased

type Container struct {
	ID        string
	Remote    string
	Namespace string
	State     ContainerState
}

func (c *Container) Create(client *containerd.Client, ctx context.Context) (*containerd.Container, error) {
	containerdContainer, err := createContainer(client, ctx, c.ID, c.Remote)
	if err != nil {
		return nil, err
	}
	c.State = Created
	return &containerdContainer, nil
}

func (c *Container) Prepare(client *containerd.Client, ctx context.Context) (*containerd.Task, error) {
	containerdContainer, err := c.containerdContainer(client, ctx)
	containerdTask, err := (*containerdContainer).NewTask(ctx, containerd.Stdio)
	if err != nil {
		return nil, err
	}
	c.State = Prepared
	return &containerdTask, nil
}

func (c *Container) Start(client *containerd.Client, ctx context.Context) (<-chan containerd.ExitStatus, error) {
	containerdTask, err := c.containerdTask(client, ctx)
	if err != nil {
		return nil, err
	}
	exitStatusC, err := (*containerdTask).Wait(ctx)
	if err != nil {
		return nil, err
	}
	err = (*containerdTask).Start(ctx)
	if err != nil {
		return nil, err
	}
	c.State = Prepared
	return exitStatusC, nil
}

func (c *Container) Sigterm(client *containerd.Client, ctx context.Context) error {
	err := c.signal(client, ctx, syscall.SIGTERM)
	if err != nil {
		return err
	}
	c.State = Sigtermed
	return nil
}

func (c *Container) Sigkill(client *containerd.Client, ctx context.Context) error {
	err := c.signal(client, ctx, syscall.SIGKILL)
	if err != nil {
		return err
	}
	c.State = Sigkilled
	return nil
}

func (c *Container) signal(client *containerd.Client, ctx context.Context, signal syscall.Signal) error {
	containerdTask, err := c.containerdTask(client, ctx)
	if err != nil {
		return err
	}
	status, err := (*containerdTask).Status(ctx)
	if err != nil {
		return err
	}
	switch status.Status {
	case containerd.Running:
		fallthrough
	case containerd.Paused:
		fallthrough
	case containerd.Pausing:
		err = (*containerdTask).Kill(ctx, signal, containerd.WithKillAll)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Container) Delete(client *containerd.Client, ctx context.Context) error {
	containerdTask, err := c.containerdTask(client, ctx)
	if err != nil {
		return err
	}
	_, err = (*containerdTask).Delete(ctx)
	if err != nil {
		return err
	}
	c.State = Deleted
	return nil
}

func (c *Container) Erase(client *containerd.Client, ctx context.Context) error {
	containerdContainer, err := c.containerdContainer(client, ctx)
	if err != nil {
		return err
	}
	err = (*containerdContainer).Delete(ctx, containerd.WithSnapshotCleanup)
	if err != nil {
		return err
	}
	c.State = Erased
	return nil
}

func (c *Container) containerdContainer(client *containerd.Client, ctx context.Context) (*containerd.Container, error) {
	containerdContainer, err := client.LoadContainer(ctx, c.ID)
	if err != nil {
		return nil, err
	}
	return &containerdContainer, err
}

func (c *Container) containerdTask(client *containerd.Client, ctx context.Context) (*containerd.Task, error) {
	containerdContainer, err := c.containerdContainer(client, ctx)
	if err != nil {
		return nil, err
	}

	containerdTask, err := (*containerdContainer).Task(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &containerdTask, nil
}
