package main

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
)

var AvailableSubcommands = []string{
	"sync",
	"help",
}

type ParsedArgs struct {
	Subcommand string
}

type SyncArgs struct {
	SourceDirPath string `cardinality:"required"`
	DestDirPath   string `cardinality:"required"`
	DeleteMissing bool   `cardinality:"optional"`
}

// realMain is the main function that does all the work.
// It takes the command line arguments and returns an exit code.
// It was moved to a separate function to make it easier to test.
func realMain(args []string) int {
	// args[0] is the program name
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: synchro <subcommand> [<args>]")
		return 2
	}

	// args[1] is the subcommand, e.g. "sync"
	// Check if the subcommand is valid
	subcommand := args[1]
	if !slices.Contains(AvailableSubcommands, subcommand) {
		slog.Error("unknown subcommand", "subcommand", subcommand)
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n", subcommand)
		fmt.Fprintln(os.Stderr, "available subcommands:", AvailableSubcommands)
		return 2
	}

	if subcommand == "sync" {
		return syncCommand(args)
	}

	fmt.Printf("Hello, %s!\n", args[1])
	slog.Info("synchro started", "subcommand", subcommand)
	return 0
}

// Main is the entry point for the program.
// It just calls realMain and exits with the returned code.
func main() {
	os.Exit(realMain(os.Args))
}
