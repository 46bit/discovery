package executor

import (
	"github.com/46bit/discovery/deployer/deployer"
	cd "github.com/containerd/containerd"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type container struct {
	desc      deployer.Container
	state     state
	container *cd.Container
	task      *task
}

func newContainer(api cdApi, desc deployer.Container) (*container, error) {
	c := container{desc: desc}
	if err := c.create(api); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *container) create(api cdApi) error {
	image, err := api.client.Pull(api.context, c.desc.Remote)
	if err != nil {
		return err
	}
	spec, err := containerSpec(api.context, c.desc.ID)
	if err != nil {
		return err
	}
	cdContainer, err := api.client.NewContainer(
		api.context,
		c.desc.ID,
		cd.WithImage(image),
		cd.WithNewSnapshot("snapshot-"+c.desc.ID, image),
		cd.WithSpec(spec, cd.WithImageConfig(image), cd.WithHostNamespace(specs.NetworkNamespace)),
	)
	if err != nil {
		return err
	}
	c.container = &cdContainer
	c.state = created
	return nil
}

func (c *container) delete(api cdApi) error {
	err := (*c.container).Delete(api.context, cd.WithSnapshotCleanup)
	if err != nil {
		return err
	}
	c.state = deleted
	return nil
}
