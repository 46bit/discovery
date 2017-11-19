package runtime

import (
	"github.com/46bit/discovery/deployer/deployer"
	"github.com/containerd/containerd"
)

type Runtime struct {
	Client     *containerd.Client
	Containers map[string]deployer.Container
	Tasks      map[string]containerd.Task
	Add        chan deployer.Container
	Remove     chan string
	Exit       chan taskExit
	Shutdown   chan interface{}
}

func NewRuntime(client *containerd.Client) Runtime {
	return Runtime{
		Client:     client,
		Containers: map[string]deployer.Container{},
		Tasks:      map[string]containerd.Task{},
		Add:        make(chan deployer.Container),
		Remove:     make(chan string),
		Shutdown:   make(chan interface{}),
	}
}

func (r *Runtime) Run() {
	for {
		select {
		case container := <-r.Add:
			r.Containers[container.ID] = container
		case containerID := <-r.Remove:
			delete(r.Containers, containerID)
		case <-r.Exit:

		case <-r.Shutdown:
			break
		}
	}
}

type taskExit struct {
	ContainerID string
	ExitCode    containerd.ExitStatus
}

// go func(containerID string, exit chan taskExit, exitStatusC <-chan containerd.ExitStatus) {
//   exitStatus := <-exitStatusC
//   exit <- taskExit{
//     ContainerID:  containerID,
//     ExitCode:   exitStatus,
//   }
// }(container.ID(), Exit, exitStatusC)
