package operator

import (
	"github.com/46bit/discovery/rainbow"
	"github.com/46bit/discovery/rainbow/executor"
	"github.com/davecgh/go-spew/spew"
	"log"
	"time"
)

const (
	namespace = "default"
)

type Operator struct {
	CmdChan             chan<- executor.Cmd
	EventChan           <-chan executor.Event
	Deployments         map[string]rainbow.Deployment
	instanceDeployments map[string]string
}

func NewOperator(cmdChan chan<- executor.Cmd, eventChan <-chan executor.Event) *Operator {
	return &Operator{
		CmdChan:             cmdChan,
		EventChan:           eventChan,
		Deployments:         map[string]rainbow.Deployment{},
		instanceDeployments: map[string]string{},
	}
}

func (o *Operator) Run() {
	for {
		select {
		case event := <-o.EventChan:
			log.Printf("event received by operator: %s\n", spew.Sdump(event))
			if event.Variant == executor.EventStopVariant {
				// Resume containers only if their deployment is still registered.
				if _, ok := o.instanceDeployments[event.Stop.InstanceID]; ok {
					cmd := executor.NewExecuteCmd(namespace, event.Stop.InstanceID, event.Stop.InstanceRemote)
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

func (o *Operator) Add(deployment rainbow.Deployment) {
	o.Deployments[deployment.Name] = deployment
	for _, instance := range deployment.Instances() {
		o.CmdChan <- executor.NewExecuteCmd(namespace, instance.ID(), instance.Remote)
		o.instanceDeployments[instance.ID()] = deployment.Name
	}
}

func (o *Operator) Remove(name string) {
	deployment := o.Deployments[name]
	delete(o.Deployments, name)
	// Stop resuming all containers before potentially-blocking channel operations.
	for _, instance := range deployment.Instances() {
		delete(o.instanceDeployments, instance.ID())
	}
	for _, instance := range deployment.Instances() {
		o.CmdChan <- executor.NewKillCmd(namespace, instance.ID())
	}
}
