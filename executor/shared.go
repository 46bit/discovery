package main

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/satori/go.uuid"
)

type Machine struct {
	Remote string `json:"remote"`
	GUID   string `json:"guid"`
}

func NewMachine(remote string) Machine {
	return Machine{
		Remote: remote,
		GUID:   uuid.NewV4().String(),
	}
}

func (m *Machine) SnapshotGUID() string {
	return m.GUID + "-snapshot"
}

type Group struct {
	Name     string             `json:"name"`
	Machines map[string]Machine `json:"machines"`
}

func NewGroup(name string, remotes []string) Group {
	group := Group{
		Name:     name,
		Machines: map[string]Machine{},
	}
	for _, remote := range remotes {
		group.Machines[remote] = NewMachine(remote)
	}
	return group
}

func runTask(machine Machine, namespace string, client *containerd.Client) (containerd.Task, error) {
	ctx := namespaces.WithNamespace(context.Background(), namespace)

	image, err := client.Pull(ctx, machine.Remote, containerd.WithPullUnpack)
	if err != nil {
		return nil, fmt.Errorf("Error pulling image: %s", err)
	}

	spec, err := containerSpec(ctx, machine.GUID)
	if err != nil {
		return nil, fmt.Errorf("Error building container spec: %s", err)
	}

	container, err := client.NewContainer(
		ctx,
		machine.GUID,
		containerd.WithSpec(spec, containerd.WithImageConfig(image), containerd.WithHostNamespace(specs.NetworkNamespace)),
		containerd.WithImage(image),
		containerd.WithNewSnapshot(machine.SnapshotGUID(), image),
	)
	if err != nil {
		return nil, fmt.Errorf("Error creating container: %s", err)
	}

	task, err := container.NewTask(ctx, containerd.Stdio)
	if err != nil {
		container.Delete(ctx, containerd.WithSnapshotCleanup)
		return nil, fmt.Errorf("Error creating task: %s", err)
	}

	return task, nil
}
