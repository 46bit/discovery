package rainbow

type Deployment struct {
	Name string `json:"name"`
	Jobs []Job  `json:"jobs"`
}

func (d *Deployment) Instances() map[string][]Instance {
	instances := map[string][]Instance{}
	for _, j := range d.Jobs {
		instances[j.Name] = []Instance{}
		for i := uint(0); i < j.InstanceCount; i++ {
			i := NewInstance(d.Name, j.Name, i, j.Remote)
			instances[j.Name] = append(instances[j.Name], i)
		}
	}
	return instances
}

type Job struct {
	Name          string `json:"name"`
	Remote        string `json:"remote"`
	InstanceCount uint   `json:"instance_count"`
}
