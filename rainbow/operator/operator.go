package operator

import (
	"github.com/46bit/discovery/rainbow"
	"github.com/46bit/discovery/rainbow/executor"
	"log"
)

const (
	NAMESPACE = "default"
)

type Operator struct {
	CmdChan     chan<- executor.Cmd
	EventChan   <-chan executor.Event
	Deployments map[string]rainbow.Deployment
}

func NewOperator(cmdChan chan<- executor.Cmd, eventChan <-chan executor.Event) *Operator {
	return &Operator{
		CmdChan:     cmdChan,
		EventChan:   eventChan,
		Deployments: map[string]rainbow.Deployment{},
	}
}

func (o *Operator) Run() {
	for {
		select {
		case event := <-o.EventChan:
			log.Println(event)
		}
	}
}

func (o *Operator) Add(deployment rainbow.Deployment) {
	o.Deployments[deployment.Name] = deployment
	for _, instance := range deployment.Instances() {
		o.CmdChan <- executor.NewExecuteCmd(NAMESPACE, instance)
	}
}

func (o *Operator) Remove(name string) {
	deployment := o.Deployments[name]
	delete(o.Deployments, name)
	for _, instance := range deployment.Instances() {
		o.CmdChan <- executor.NewKillCmd(NAMESPACE, instance.ID())
	}
}
