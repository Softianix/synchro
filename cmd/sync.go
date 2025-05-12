// Handling of the `sync` subcommand.
package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/Devourian/synchro/internal/core/sync"
	"github.com/Devourian/synchro/internal/enums"
	"github.com/Devourian/synchro/internal/flags"
	"github.com/Devourian/synchro/internal/validators"
)

// Definitions of CLI flags for `sync` subcommand.
var (
	srcFlag = flags.Flag[string]{
		Usage:       "Path to source directory to copy from, absolute or relative",
		LongName:    "source_dir",
		ShortName:   "src",
		Required:    true,
		Default:     "",
		ParsedValue: new(string),
		Validators:  []func(string) bool{validators.IsDir},
	}
	dstFlag = flags.Flag[string]{
		Usage:       "Path to destination directory to copy to, absolute or relative",
		LongName:    "destination_dir",
		ShortName:   "dst",
		Required:    true,
		Default:     "",
		ParsedValue: new(string),
		Validators:  []func(string) bool{validators.IsDir},
	}
	delFlag = flags.Flag[bool]{
		Usage:       "Delete files in destination directory not present in source",
		LongName:    "delete_missing",
		ShortName:   "d",
		Required:    false,
		Default:     false,
		ParsedValue: new(bool),
		Validators:  []func(bool) bool{},
	}
)

// syncFlagRegistry is a map to keep track of all registered flags for `sync` subcommand.
// Can be used to retrieve the flag by name with O(1) time complexity.
// Also allows to iterate through all flags.
var syncFlagRegistry = make(map[string]any)

// Handler for the `sync` subcommand.
func SyncHandler(subArgs []string) int {

	// Register the CLI flags
	syncFlagSet := flag.NewFlagSet("sync", flag.ContinueOnError)
	flags.RegisterFlags(
		syncFlagSet,
		syncFlagRegistry,
		&srcFlag,
		&dstFlag,
		&delFlag,
	)

	// Parse subcommand arguments
	if err := syncFlagSet.Parse(subArgs); err != nil {
		return enums.ExitCodeFlagError.Value()
	}

	// Check if the flags were set correctly
	if syncFlagSet.Parsed() {
		result := flags.CheckRequiredFlagsWereSet(syncFlagRegistry, syncFlagSet)
		if result != 0 {
			return result
		}
	} else {
		fmt.Fprintf(os.Stderr, "Error: flags were not parsed correctly\n\n")
		syncFlagSet.Usage()
		return enums.ExitCodeFlagError.Value()
	}

	// Validate all flags
	flags.ValidateAllFlags(syncFlagRegistry)

	// Run the business logic for `sync` subcommand
	err := sync.Sync(
		*srcFlag.ParsedValue,
		*dstFlag.ParsedValue,
		*delFlag.ParsedValue,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return enums.ExitCodeError.Value()
	}

	return enums.ExitCodeSuccess.Value()
}
