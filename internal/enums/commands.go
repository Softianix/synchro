package enums

type Command string

const (
	CmdHelp Command = "help"
	CmdSync Command = "sync"
)

var AllCommands = []Command{
	CmdHelp,
	CmdSync,
}

func (cmd Command) String() string {
	return string(cmd)
}
