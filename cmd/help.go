// Handling of the `sync` subcommand.
package cmd

import "fmt"

func HelpHandler(subArgs []string) int {
	helpMessage := `
		Usage: synchro [command] [flags]
		Available commands:
		  sync   Synchronize files between two directories
		  help   Show this help message
		Use 'synchro [command] --help' for more information about a command.`

	fmt.Println(helpMessage)

	return 0
}
