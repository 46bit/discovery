package executor

import (
	"github.com/46bit/discovery/rainbow/container"
	"github.com/containerd/containerd"
)

// Executor
// - takes commands to Execute and Kill a type of Container
// - emits events when a Container has started and has stopped

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
				e.execute(command.Execute.ID, command.Execute.Remote)
			case CmdKillVariant:
				e.kill(command.Kill.ID)
			}
		case exitedInstanceID := <-e.exitChan:
			e.remove(exitedInstanceID)
		}
	}
}

func (e *Executor) execute(id, remote string) {
	e.containers[id] = container.NewContainer(id, remote)
	e.containers[id].Create(e.client)
	e.containers[id].Task()
	containerExitChan, _ := e.containers[id].Start()
	go func(id string, containerExitChan <-chan containerd.ExitStatus, exitChan chan<- string) {
		<-containerExitChan
		exitChan <- id
	}(id, containerExitChan, e.exitChan)
	e.EventChan <- NewStartEvent(id)
}

func (e *Executor) kill(id string) {
	if _, ok := e.containers[id]; !ok {
		return
	}
	e.containers[id].Stop()
}

func (e *Executor) remove(id string) {
	if _, ok := e.containers[id]; !ok {
		return
	}
	e.containers[id].Untask()
	e.containers[id].Delete()
	remote := e.containers[id].Remote
	delete(e.containers, id)
	e.EventChan <- NewStopEvent(id, remote)
}
