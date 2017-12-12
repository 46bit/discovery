package executor

import (
	"github.com/46bit/discovery/rainbow/instance"
	"github.com/containerd/containerd"
)

type Executor struct {
	CmdChan   chan Cmd
	EventChan chan Event
	exitChan  chan string
	client    *containerd.Client
	instances map[string]*instance.Instance
}

func NewExecutor(client *containerd.Client) *Executor {
	return &Executor{
		CmdChan:   make(chan Cmd),
		EventChan: make(chan Event),
		exitChan:  make(chan string),
		client:    client,
		instances: map[string]*instance.Instance{},
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
			exitedInstance := e.instances[exitedInstanceID]
			e.remove(exitedInstance.Namespace, exitedInstanceID)
		}
	}
}

func (e *Executor) execute(namespace, instanceID, instanceRemote string) {
	e.instances[instanceID] = instance.NewInstance(namespace, instanceID, instanceRemote)
	e.instances[instanceID].Create(e.client)
	e.instances[instanceID].Task()
	instanceExitChan, _ := e.instances[instanceID].Start()
	go func(instanceID string, instanceExitChan <-chan containerd.ExitStatus, exitChan chan<- string) {
		<-instanceExitChan
		exitChan <- instanceID
	}(instanceID, instanceExitChan, e.exitChan)
	e.EventChan <- NewStartEvent(namespace, instanceID)
}

func (e *Executor) kill(instanceID string) {
	if _, ok := e.instances[instanceID]; !ok {
		return
	}
	e.instances[instanceID].Stop()
}

func (e *Executor) remove(namespace, instanceID string) {
	if _, ok := e.instances[instanceID]; !ok {
		return
	}
	e.instances[instanceID].Untask()
	e.instances[instanceID].Delete()
	delete(e.instances, instanceID)
	e.EventChan <- NewStopEvent(namespace, instanceID, e.instances[instanceID].Remote)
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
