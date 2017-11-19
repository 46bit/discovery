package deployer

import "fmt"

// deployment:
//   name: "senders-receiver"
//   jobs:
//   - name: "sender"
//     remote: "docker.io/46bit/sender:latest"
//     replicas: 3
//   - name: "receiver"
//     remote: "docker.io/46bit/receiver:latest"
//     replicas: 1

type Deployment struct {
	Name string `json:"name"`
	Jobs []Job  `json:"jobs"`
}

type Job struct {
	Name     string `json:"name"`
	Remote   string `json:"remote"`
	Replicas uint   `json:"replicas"`
}

func (j *Job) ContainerID(replicaNumber uint) string {
	return fmt.Sprintf("%s.%s", j.Name, replicaNumber)
}

func (j *Job) Containers(namespace string) []Container {
	containers := []Container{}
	for i := uint(0); i < j.Replicas; i++ {
		containers[i] = Container{
			ID:        j.ContainerID(i),
			Remote:    j.Remote,
			Namespace: namespace,
		}
	}
	return containers
}

type Container struct {
	ID        string
	Remote    string
	Namespace string
}
