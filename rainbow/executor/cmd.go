package executor

type CmdVariant uint

const (
	CmdExecuteVariant CmdVariant = iota
	CmdKillVariant
)

type Cmd struct {
	Variant CmdVariant
	Execute *CmdExecute
	Kill    *CmdKill
}

func NewExecuteCmd(id, remote string) Cmd {
	return Cmd{
		Variant: CmdExecuteVariant,
		Execute: &CmdExecute{ID: id, Remote: remote},
	}
}

func NewKillCmd(id string) Cmd {
	return Cmd{
		Variant: CmdKillVariant,
		Kill:    &CmdKill{ID: id},
	}
}

type CmdExecute struct {
	ID     string
	Remote string
}

type CmdKill struct {
	ID string
}
