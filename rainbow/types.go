package rainbow

import (
	"fmt"
)

type Deployment struct {
	Name string `json:"name"`
	Jobs []Job  `json:"jobs"`
}

func (d *Deployment) Instances() []Instance {
	instances := []Instance{}
	for _, j := range d.Jobs {
		for i := uint(0); i < j.InstanceCount; i++ {
			instances = append(instances, Instance{
				Index:          i,
				Remote:         j.Remote,
				JobName:        j.Name,
				DeploymentName: d.Name,
			})
		}
	}
	return instances
}

type Job struct {
	Name          string `json:"name"`
	Remote        string `json:"remote"`
	InstanceCount uint   `json:"instance_count"`
}

type Instance struct {
	Index          uint
	Remote         string
	JobName        string
	DeploymentName string
}

func (i *Instance) ID() string {
	return fmt.Sprintf("%s.%s.i%d", i.DeploymentName, i.JobName, i.Index)
}
