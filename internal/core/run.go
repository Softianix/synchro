package core

import (
	"fmt"
	"os"

	"github.com/Devourian/synchro/cmd"
	"github.com/Devourian/synchro/internal/enums"
	"github.com/Devourian/synchro/internal/logging"
)

var log = logging.GetLogger()

// Run function does all the actual work.
// It was extracted from main function to make the solution easy to test.
// Extracting it + using FlagSets allows to run the function with different
// arguments, simulating different user inputs for testing purposes.
func Run(cliArgs []string) int {
	if len(cliArgs) < 2 {
		fmt.Fprintf(
			os.Stderr,
			"Error: no subcommand provided\n"+
				"Available subcommands: %s\n",
			enums.AllCommands,
		)
		return enums.ExitCodeFlagError.Value()
	}

	// Get subcommand name and arguments
	subcommand := cliArgs[1]
	subArgs := []string{}
	if len(cliArgs) > 2 {
		subArgs = cliArgs[2:]
	}

	return subcommandDispatch(
		enums.Command(subcommand),
		subArgs,
	)
}

// Dispatch program flow to the appropriate
// subcommand handler, if valid subcommand given.
func subcommandDispatch(
	subcommand enums.Command,
	subArgs []string,
) int {
	switch subcommand {

	case enums.CmdHelp:
		return cmd.HelpHandler(subArgs)

	case enums.CmdSync:
		return cmd.SyncHandler(subArgs)

	default:
		fmt.Fprintf(
			os.Stderr,
			"Error: unknown subcommand %s\n"+
				"Available subcommands: %s\n",
			subcommand,
			enums.AllCommands,
		)
		return enums.ExitCodeFlagError.Value()
	}
}
