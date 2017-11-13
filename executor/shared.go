package main

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
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

// func runExecutionSet(execution_set ExecutionSet, namespace string, client containerd.Client) {
// 	for _, execution := range execution_set.Items {
// 		runExecution(execution, namespace, client)
// 	}
// }

func runTask(machine Machine, namespace string, client *containerd.Client) (containerd.Task, error) {
	ctx := namespaces.WithNamespace(context.Background(), namespace)

	image, err := client.Pull(ctx, machine.Remote, containerd.WithPullUnpack)
	if err != nil {
		return nil, fmt.Errorf("Error pulling image: %s", err)
	}

	container, err := client.NewContainer(
		ctx,
		machine.GUID,
		containerd.WithImage(image),
		containerd.WithNewSnapshot(machine.SnapshotGUID(), image),
		containerd.WithNewSpec(containerd.WithImageConfig(image), withHostNetworkNamespace),
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

// adapted from WithHtop https://github.com/containerd/containerd/blob/a6ce1ef2a140d79856a8647e1d1ae5ac9ab581eb/docs/client-opts.md
func withHostNetworkNamespace(context context.Context, client *containerd.Client, container *containers.Container, s *specs.Spec) error {
	// make sure we are in the host network namespace
	if err := containerd.WithHostNamespace(specs.NetworkNamespace)(context, client, container, s); err != nil {
		return err
	}
	return nil
}

// func stopExecutionSet(execution_set ExecutionSet, namespace string, client containerd.Client) []error {
// 	errs := []error{}
// 	for _, execution := range execution_set.Items {
// 		if err := stopExecution(execution, namespace, client); err != nil {
// 			errs = append(errs, err)
// 		}
// 	}
// 	return errs
// }

// func stopExecution(execution Execution, namespace string, client containerd.Client) error {
// 	ctx := namespaces.WithNamespace(context.Background(), namespace)

// 	if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
// 		return err
// 	}

// 	status := <-exitStatusC
// 	_, _, err := status.Result()
// 	if err != nil {
// 		return err
// 	}

// 	if container, err := client.LoadContainer(ctx, execution.ContainerGUID); err != nil {
// 		return err
// 	}
// 	if err = container.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
// 		return err
// 	}

// 	return nil
// }
