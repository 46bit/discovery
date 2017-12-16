package executor

import (
	"github.com/46bit/discovery/rainbow/container"
	"github.com/containerd/containerd"
)

type Executor struct {
	CmdChan    chan Cmd
	EventChan  chan Event
	exitChan   chan string
	client     *containerd.Client
	containers map[string]*container.Container
}

func NewExecutor(client *containerd.Client) *Executor {
	return &Executor{
		CmdChan:    make(chan Cmd),
		EventChan:  make(chan Event),
		exitChan:   make(chan string),
		client:     client,
		containers: map[string]*container.Container{},
	}
}

func (e *Executor) Run() {
	for {
		select {
		case command := <-e.CmdChan:
			switch command.Variant {
			case CmdExecuteVariant:
				e.execute(command.Execute.Namespace, command.Execute.InstanceID, command.Execute.InstanceRemote)
			case CmdKillVariant:
				e.kill(command.Kill.InstanceID)
			}
		case exitedInstanceID := <-e.exitChan:
			exitedInstance := e.containers[exitedInstanceID]
			e.remove(exitedInstance.Namespace, exitedInstanceID)
		}
	}
}

func (e *Executor) execute(namespace, containerID, containerRemote string) {
	e.containers[containerID] = container.NewContainer(namespace, containerID, containerRemote)
	e.containers[containerID].Create(e.client)
	e.containers[containerID].Task()
	containerExitChan, _ := e.containers[containerID].Start()
	go func(containerID string, containerExitChan <-chan containerd.ExitStatus, exitChan chan<- string) {
		<-containerExitChan
		exitChan <- containerID
	}(containerID, containerExitChan, e.exitChan)
	e.EventChan <- NewStartEvent(namespace, containerID)
}

func (e *Executor) kill(containerID string) {
	if _, ok := e.containers[containerID]; !ok {
		return
	}
	e.containers[containerID].Stop()
}

func (e *Executor) remove(namespace, containerID string) {
	if _, ok := e.containers[containerID]; !ok {
		return
	}
	e.containers[containerID].Untask()
	e.containers[containerID].Delete()
	containerRemote := e.containers[containerID].Remote
	delete(e.containers, containerID)
	e.EventChan <- NewStopEvent(namespace, containerID, containerRemote)
}

// func (e *containerExecutor) runTowardsStarted(client *cd.Client) error {
// 	var err error
// 	switch *e.container.State() {
// 	case run.ContainerDescribed, run.ContainerDeleted:
// 		err = e.container.Create(client)
// 	case run.ContainerCreated:
// 		s, err := e.container.Task.State()
// 		if err != nil {
// 			return err
// 		}
// 		switch *s {
// 		case run.TaskDescribed, run.TaskDeleted:
// 			err = e.container.Task.Create(client, e.container.Container)
// 		case run.TaskCreated:
// 			err = e.container.Task.Start(client)
// 		case run.TaskStarted:
// 		case run.TaskStopped:
// 			err = e.container.Task.Delete(client)
// 		}
// 	}
// 	return err
// }

// func (e *containerExecutor) runTowardsDeleted(client *cd.Client) error {
// 	var err error
// 	switch *e.container.State() {
// 	case run.ContainerDescribed, run.ContainerDeleted:
// 	case run.ContainerCreated:
// 		s, err := e.container.Task.State()
// 		if err != nil {
// 			return err
// 		}
// 		switch *s {
// 		case run.TaskDescribed, run.TaskDeleted:
// 			err = e.container.Delete(client)
// 		case run.TaskCreated:
// 			err = e.container.Task.Delete(client)
// 		case run.TaskStarted:
// 			err = e.container.Task.Stop(client)
// 		case run.TaskStopped:
// 			err = e.container.Task.Delete(client)
// 		}
// 	}
// 	return err
// }

// func (e *containerExecutor) targetReached() (*bool, error) {
// 	if e.targetState == containerStarted {
// 		s, err := e.container.Task.State()
// 		if err != nil {
// 			return nil, err
// 		}
// 		reached := *s == run.TaskStarted
// 		return &reached, nil
// 	} else if e.targetState == containerDeleted {
// 		reached := *e.container.State() == run.ContainerDeleted
// 		return &reached, nil
// 	} else {
// 		return nil, fmt.Errorf("Unknown target state %s.", e.targetState)
// 	}
// }
