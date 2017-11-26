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

func snapshotID(containerID string) string {
	return "snapshot." + containerID
}
