package executor

import (
	"github.com/46bit/discovery/rainbow"
)

type CmdVariant uint

const (
	CmdExecuteVariant CmdVariant = iota
	CmdKillVariant
)

type Cmd struct {
	Variant CmdVariant  `json:"variant"`
	Execute *CmdExecute `json:"execute"`
	Kill    *CmdKill    `json:"kill"`
}

func NewExecuteCmd(namespace string, instance rainbow.Instance) Cmd {
	return Cmd{
		Variant: CmdExecuteVariant,
		Execute: &CmdExecute{Namespace: namespace, Instance: instance},
	}
}

func NewKillCmd(namespace, instanceID string) Cmd {
	return Cmd{
		Variant: CmdKillVariant,
		Kill:    &CmdKill{Namespace: namespace, InstanceID: instanceID},
	}
}

type CmdExecute struct {
	Namespace string           `json:"namespace"`
	Instance  rainbow.Instance `json:"instance"`
}

type CmdKill struct {
	Namespace  string `json:"namespace"`
	InstanceID string `json:"instance_id"`
}
