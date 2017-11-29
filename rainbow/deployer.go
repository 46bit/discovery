package rainbow

import (
	"github.com/46bit/discovery/rainbow/executor"
)

type Deployer struct {
	Runtime     *executor.Runtime
	Deployments map[string]Deployment
	Add         chan Deployment
	Remove      chan string
}

func NewDeployer(runtime *executor.Runtime) *Deployer {
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
		containerDescs := job.ContainerDescs(deployment.Namespace())
		for _, containerDesc := range containerDescs {
			d.Runtime.Add <- containerDesc
		}
	}
}

func (d *Deployer) remove(name string) {
	deployment := d.Deployments[name]
	for _, job := range deployment.Jobs {
		containerDescs := job.ContainerDescs(deployment.Namespace())
		for _, containerDesc := range containerDescs {
			d.Runtime.Remove <- containerDesc.ID
		}
	}
}
