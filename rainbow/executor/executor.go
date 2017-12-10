package executor

import (
	"github.com/46bit/discovery/rainbow/instance"
	"github.com/containerd/containerd"
)

type Executor struct {
	CmdChan   chan Cmd
	EventChan chan Event
	client    *containerd.Client
	instances map[string]*instance.Instance
}

func NewExecutor(client *containerd.Client) *Executor {
	return &Executor{
		CmdChan:   make(chan Cmd),
		EventChan: make(chan Event),
		client:    client,
		instances: map[string]*instance.Instance{},
	}
}

func (e *Executor) Run() {
	for {
		command := <-e.CmdChan
		switch command.Variant {
		case CmdExecuteVariant:
			e.execute(*command.Execute)
		case CmdKillVariant:
			e.kill(*command.Kill)
		}
	}
}

func (e *Executor) execute(c CmdExecute) {
	e.instances[c.InstanceID] = instance.NewInstance(c.Namespace, c.InstanceID, c.InstanceRemote)
	e.instances[c.InstanceID].Create(e.client)
	e.instances[c.InstanceID].Task()
	e.instances[c.InstanceID].Start()
	e.EventChan <- NewStartEvent(c.Namespace, instanceID)
}

func (e *Executor) kill(c CmdKill) {
	if _, ok := e.instances[c.InstanceID]; !ok {
		return
	}
	e.instances[c.InstanceID].Stop()
	e.instances[c.InstanceID].Untask()
	e.instances[c.InstanceID].Delete()
	delete(e.instances, c.InstanceID)
	e.EventChan <- NewStopEvent(c.Namespace, c.InstanceID)
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
