package instance

import (
	"context"
	"github.com/46bit/discovery/rainbow"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
	"sync"
	"syscall"
	"time"
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

type Instance struct {
	namespace string
	desc      rainbow.Instance
	state     State
	container *containerd.Container
	task      *containerd.Task
	sync.Mutex
}

func NewInstance(namespace string, desc rainbow.Instance) *Instance {
	return &Instance{
		namespace: namespace,
		desc:      desc,
		state:     Described,
	}
}

func (i *Instance) Create(client *containerd.Client) error {
	i.Lock()
	defer i.Unlock()
	ctx := i.context()
	image, err := client.Pull(ctx, i.desc.Remote, containerd.WithPullUnpack)
	if err != nil {
		return errors.Wrap(err, "Error pulling image")
	}
	imageConfig := containerd.WithImageConfig(image)
	withHostNamespace := containerd.WithHostNamespace(specs.NetworkNamespace)
	container, err := client.NewContainer(
		ctx,
		i.desc.ID(),
		containerd.WithImage(image),
		containerd.WithNewSnapshot("snapshot-"+i.desc.ID(), image),
		containerd.WithNewSpec(imageConfig, withHostNamespace),
	)
	if err != nil {
		return errors.Wrap(err, "Error creating new container")
	}
	i.state = Created
	i.container = &container
	return nil
}

func (i *Instance) Task() error {
	i.Lock()
	defer i.Unlock()
	task, err := (*i.container).NewTask(i.context(), containerd.Stdio)
	if err != nil {
		return errors.Wrap(err, "Error creating containerd task")
	}
	i.state = Tasked
	i.task = &task
	return nil
}

func (i *Instance) Start() error {
	i.Lock()
	defer i.Unlock()
	err := (*i.task).Start(i.context())
	if err != nil {
		return errors.Wrap(err, "Error starting containerd task")
	}
	i.state = Started
	return nil
}

func (i *Instance) Stop() error {
	i.Lock()
	defer i.Unlock()
	ctx := i.context()
	exitStatusC, err := (*i.task).Wait(ctx)
	if err != nil {
		return errors.Wrap(err, "Error waiting for containerd task")
	}
	err = (*i.task).Kill(ctx, syscall.SIGTERM, containerd.WithKillAll)
	if err != nil {
		return errors.Wrap(err, "Error SIGTERMing containerd task")
	}
	select {
	case <-exitStatusC:
	case <-time.After(5 * time.Second):
		err = (*i.task).Kill(ctx, syscall.SIGKILL, containerd.WithKillAll)
		if err != nil {
			return errors.Wrap(err, "Error SIGKILLing containerd task")
		}
		<-exitStatusC
	}
	i.state = Stopped
	return nil
}

func (i *Instance) Untask() error {
	i.Lock()
	defer i.Unlock()
	_, err := (*i.task).Delete(i.context())
	if err != nil {
		return errors.Wrap(err, "Error deleting containerd task")
	}
	i.state = Untasked
	i.task = nil
	return nil
}

func (i *Instance) Delete() error {
	i.Lock()
	defer i.Unlock()
	err := (*i.container).Delete(i.context(), containerd.WithSnapshotCleanup)
	if err != nil {
		return errors.Wrap(err, "Error deleting container")
	}
	i.state = Deleted
	i.container = nil
	return nil
}

func (i *Instance) Status() State {
	i.Lock()
	defer i.Unlock()
	if i.state == Started {
		status, err := (*i.task).Status(i.context())
		if err == nil && status.Status != containerd.Running {
			i.state = Stopped
		}
	}
	return i.state
}

func (i *Instance) context() context.Context {
	return namespaces.WithNamespace(context.Background(), i.namespace)
}
