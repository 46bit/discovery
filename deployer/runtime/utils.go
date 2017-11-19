package runtime

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func createContainer(client *containerd.Client, ctx context.Context, id string, remote string) (containerd.Container, error) {
	image, err := client.Pull(ctx, remote, containerd.WithPullUnpack)
	if err != nil {
		return nil, fmt.Errorf("Error pulling image for %s: %s", id, err)
	}

	spec, err := containerSpec(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Error building container %s spec: %s", id, err)
	}

	container, err := client.NewContainer(
		ctx,
		id,
		containerd.WithImage(image),
		containerd.WithNewSnapshot(snapshotID(id), image),
		containerd.WithSpec(spec, containerd.WithImageConfig(image), containerd.WithHostNamespace(specs.NetworkNamespace)),
	)
	if err != nil {
		return nil, fmt.Errorf("Error creating container %s: %s", id, err)
	}

	return container, nil
}

func createTask(ctx context.Context, container containerd.Container) (containerd.Task, <-chan containerd.ExitStatus, error) {
	task, err := container.NewTask(ctx, containerd.Stdio)
	if err != nil {
		container.Delete(ctx, containerd.WithSnapshotCleanup)
		return nil, nil, fmt.Errorf("Error creating task %s: %s", container.ID(), err)
	}

	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("Error waiting for task %s: %s", container.ID(), err)
	}

	err = task.Start(ctx)
	if err != nil {
		task.Delete(ctx)
		return nil, nil, fmt.Errorf("Error starting task %s: %s", container.ID(), err)
	}

	return task, exitStatusC, nil
}

func snapshotID(containerID string) string {
	return "snapshot." + containerID
}
