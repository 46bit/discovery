package rainbow

import (
	"github.com/46bit/discovery/rainbow/executor"
	"github.com/davecgh/go-spew/spew"
	"log"
	"time"
)

type Operator struct {
	CmdChan     chan<- executor.Cmd
	EventChan   <-chan executor.Event
	Deployments map[string]Deployment
	Instances   map[string]Instance
}

func NewOperator(cmdChan chan<- executor.Cmd, eventChan <-chan executor.Event) *Operator {
	return &Operator{
		CmdChan:     cmdChan,
		EventChan:   eventChan,
		Deployments: map[string]Deployment{},
		Instances:   map[string]Instance{},
	}
}

func (o *Operator) Run() {
	for {
		select {
		case event := <-o.EventChan:
			log.Printf("event received by operator: %s\n", spew.Sdump(event))
			if event.Variant == executor.EventStartVariant {
				if i, ok := o.Instances[event.Start.ID]; ok {
					i.State = InstanceStarted
				}
			} else if event.Variant == executor.EventStopVariant {
				// Resume containers only if their deployment is still registered.
				if i, ok := o.Instances[event.Stop.ID]; ok {
					i.State = InstanceStopped
					cmd := executor.NewExecuteCmd(event.Stop.ID, event.Stop.Remote)
					// Executor does not properly handle starting to execute something before it is
					// entirely deleted. Until this behaviour is more solid, this hackily delays the
					// re-execute command.
					go func(cmdChan chan<- executor.Cmd, cmd executor.Cmd) {
						<-time.After(10 * time.Second)
						cmdChan <- cmd
					}(o.CmdChan, cmd)
				}
			}
		}
	}
}

func (o *Operator) Add(deployment Deployment) {
	o.Deployments[deployment.Name] = deployment
	for _, jobInstances := range deployment.Instances() {
		for _, i := range jobInstances {
			o.CmdChan <- executor.NewExecuteCmd(i.ID, i.Remote)
			o.Instances[i.ID] = i
		}
	}
}

func (o *Operator) Remove(name string) {
	deployment := o.Deployments[name]
	delete(o.Deployments, name)
	// Stop resuming all containers before potentially-blocking channel operations.
	for _, jobInstances := range deployment.Instances() {
		for _, i := range jobInstances {
			delete(o.Instances, i.ID)
		}
	}
	for _, jobInstances := range deployment.Instances() {
		for _, i := range jobInstances {
			o.CmdChan <- executor.NewKillCmd(i.ID)
		}
	}
}
