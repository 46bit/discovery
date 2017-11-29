package rainbow

import (
	"fmt"
	"github.com/46bit/discovery/rainbow/executor"
)

type Deployment struct {
	Name string `json:"name"`
	Jobs []Job  `json:"jobs"`
}

func (d *Deployment) Namespace() string {
	return "deployment." + d.Name
}

type Job struct {
	Name      string `json:"name"`
	Remote    string `json:"remote"`
	Instances uint   `json:"instances"`
}

func (j *Job) ContainerID(instanceNumber uint) string {
	return fmt.Sprintf("%s.%d", j.Name, instanceNumber)
}

func (j *Job) ContainerDescs(namespace string) []executor.ContainerDesc {
	cs := []executor.ContainerDesc{}
	for i := uint(0); i < j.Instances; i++ {
		cs = append(cs, executor.ContainerDesc{
			ID:        j.ContainerID(i),
			Remote:    j.Remote,
			Namespace: namespace,
		})
	}
	return cs
}
