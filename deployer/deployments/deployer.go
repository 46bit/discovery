package deployments

import (
	"github.com/46bit/discovery/deployer/containers"
)

type Deployer struct {
	Runtime     *containers.Runtime
	Deployments map[string]Deployment
	Add         chan Deployment
	Remove      chan string
}

func NewDeployer(runtime *containers.Runtime) *Deployer {
	return &Deployer{
		Runtime:     runtime,
		Deployments: map[string]Deployment{},
		Add:         make(chan Deployment),
		Remove:      make(chan string),
	}
}

func (d *Deployer) Run() {
	for {
		select {
		case deployment := <-d.Add:
			d.Deployments[deployment.Name] = deployment
			d.add(deployment.Name)
		case deploymentName := <-d.Remove:
			d.remove(deploymentName)
			delete(d.Deployments, deploymentName)
		}
	}
}

func (d *Deployer) add(name string) {
	deployment := d.Deployments[name]
	for _, job := range deployment.Jobs {
		containers := job.ContainerDescs(deployment.Namespace())
		for _, container := range containers {
			d.Runtime.Add <- container
		}
	}
}

func (d *Deployer) remove(name string) {
	deployment := d.Deployments[name]
	for _, job := range deployment.Jobs {
		containers := job.ContainerDescs(deployment.Namespace())
		for _, container := range containers {
			d.Runtime.Remove <- container.ID
		}
	}
}
