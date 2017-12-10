package executor

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

func NewExecuteCmd(namespace, instanceID, instanceRemote string) Cmd {
	return Cmd{
		Variant: CmdExecuteVariant,
		Execute: &CmdExecute{Namespace: namespace, InstanceID: instanceID, InstanceRemote: instanceRemote},
	}
}

func NewKillCmd(namespace, instanceID string) Cmd {
	return Cmd{
		Variant: CmdKillVariant,
		Kill:    &CmdKill{Namespace: namespace, InstanceID: instanceID},
	}
}

type CmdExecute struct {
	Namespace      string `json:"namespace"`
	InstanceID     string `json:"instance_id"`
	InstanceRemote string `json:"instance_remote"`
}

type CmdKill struct {
	Namespace  string `json:"namespace"`
	InstanceID string `json:"instance_id"`
}
