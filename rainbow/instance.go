package rainbow

import (
	"fmt"
)

type InstanceState uint

const (
	InstanceStopped InstanceState = iota
	InstanceStarted
)

type InstanceMeta struct {
	DeploymentName string `json:"deployment_name"`
	JobName        string `json:"job_name"`
	Index          uint   `json:"index"`
}

func (iid *InstanceMeta) ID() string {
	return fmt.Sprintf("%s.%s.i%d", iid.DeploymentName, iid.JobName, iid.Index)
}

type Instance struct {
	InstanceMeta
	ID     string        `json:"id"`
	Remote string        `json:"remote"`
	State  InstanceState `json:"state,omitempty"`
}

func NewInstance(deploymentName, jobName string, index uint, remote string) Instance {
	instanceMeta := InstanceMeta{
		DeploymentName: deploymentName,
		JobName:        jobName,
		Index:          index,
	}
	return Instance{
		InstanceMeta: instanceMeta,
		ID:           instanceMeta.ID(),
		Remote:       remote,
		State:        InstanceStopped,
	}
}

type Instances struct {
}

func (i *Instances) AddDeployment() {}

func (i *Instances) RemoveDeployment() {}

func (i *Instances) GetByDeployment() {}

func (i *Instances) GetByID() {}
