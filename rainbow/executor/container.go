package executor

import (
	cd "github.com/containerd/containerd"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/pkg/errors"
)

type container struct {
	desc      ContainerDesc
	state     state
	container *cd.Container
	task      *task
}

func newContainer(api cdApi, desc ContainerDesc) (*container, error) {
	c := container{desc: desc}
	if err := c.create(api); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *container) create(api cdApi) error {
	image, err := api.client.Pull(api.context, c.desc.Remote, cd.WithPullUnpack)
	if err != nil {
		return errors.Wrap(err, "Error pulling image")
	}
	spec, err := containerSpec(api.context, c.desc.ID)
	if err != nil {
		return errors.Wrap(err, "Error generating container spec")
	}
	cdContainer, err := api.client.NewContainer(
		api.context,
		c.desc.ID,
		cd.WithImage(image),
		cd.WithNewSnapshot("snapshot-"+c.desc.ID, image),
		cd.WithSpec(spec, cd.WithImageConfig(image), cd.WithHostNamespace(specs.NetworkNamespace)),
	)
	if err != nil {
		return errors.Wrap(err, "Error creating new container")
	}
	c.container = &cdContainer
	c.state = created
	return nil
}

func (c *container) delete(api cdApi) error {
	err := (*c.container).Delete(api.context, cd.WithSnapshotCleanup)
	if err != nil {
		return errors.Wrap(err, "Error deleting container")
	}
	c.state = deleted
	return nil
}
