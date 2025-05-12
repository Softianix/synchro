// The main issue with stdlib flag package is that it doesn't support
// short and long flag names out of the box as well as handling of optional / required flags.
//
// It's also not very extensible and doesn't allow to easily add new flags.
// Moreover solution for testing the CLI with different arguments had to be found.
//
// Issues mentioned above were solved in this solution, while still
// keeping project stdlib only - without any third party dependencies.
//
// The solution is extensible and allows to easily add new flags in the future.
package flags

import (
	"flag"
	"fmt"
	"os"

	"github.com/Devourian/synchro/internal/enums"
)

// Flag is a generic struct to define command line flags.
// It contains all the necessary information and data about the flag.
type Flag[T any] struct {
	Usage       string
	LongName    string
	ShortName   string
	Required    bool
	Default     T
	ParsedValue *T
	Validators  []func(T) bool
}

// RegisterFlags registers the flags using the flagSet and adds them to the flagRegistry.
// Both long and short flag names registering works.
// Selects proper registration method based on the type of the flag / dispatch using type switch.
// For now bool and string types are supported, but it can be easily extended to support other types.
func RegisterFlags(
	flagSet *flag.FlagSet,
	flagRegistry map[string]any,
	flags ...any,
) {
	for _, untypedFlag := range flags {
		var longName string
		switch flg := untypedFlag.(type) {

		case *Flag[string]:
			flagSet.StringVar(
				flg.ParsedValue,
				flg.LongName,
				flg.Default,
				flg.Usage,
			)
			flagSet.StringVar(
				flg.ParsedValue,
				flg.ShortName,
				flg.Default,
				fmt.Sprintf("shorthand for --%s", flg.LongName),
			)
			longName = flg.LongName

		case *Flag[bool]:
			flagSet.BoolVar(
				flg.ParsedValue,
				flg.LongName,
				flg.Default,
				flg.Usage,
			)
			flagSet.BoolVar(
				flg.ParsedValue,
				flg.ShortName,
				flg.Default,
				fmt.Sprintf("shorthand for --%s", flg.LongName),
			)
			longName = flg.LongName

		default:
			panic(fmt.Sprintf("registerFlags: unsupported type %T", flg))
		}

		// Add flag to the registry
		flagRegistry[longName] = untypedFlag
	}
}

// CheckRequiredFlagsWereSet checks if all required/non-optional flags were set.
func CheckRequiredFlagsWereSet(
	flagRegistry map[string]any,
	flagSet *flag.FlagSet,
) int {
	var missing bool

	for _, untypedFlag := range flagRegistry {
		switch flg := untypedFlag.(type) {

		case *Flag[string]:
			if flg.Required && *flg.ParsedValue == flg.Default {
				fmt.Fprintf(os.Stderr, "Error: --%s/--%s argument is required\n", flg.LongName, flg.ShortName)
				missing = true
			}

		case *Flag[bool]:
			if flg.Required && *flg.ParsedValue == flg.Default {
				fmt.Fprintf(os.Stderr, "Error: --%s/--%s argument is required\n", flg.LongName, flg.ShortName)
				missing = true
			}
		}
	}

	if missing {
		fmt.Fprintf(os.Stderr, "\n")
		flagSet.Usage()
		return enums.ExitCodeFlagError.Value()
	}

	return enums.ExitCodeSuccess.Value()
}

// ValidateAllFlags validates all flags using the validators provided in the Flag struct.
func ValidateAllFlags(
	flagRegistry map[string]any,
) {
	for _, untypedFlag := range flagRegistry {
		switch flg := untypedFlag.(type) {

		case *Flag[string]:
			validateFlag(flg)

		case *Flag[bool]:
			validateFlag(flg)
		}
	}
}

// validateFlag validates a single flag using the validators provided in the Flag struct.
func validateFlag[T any](
	flg *Flag[T],
) bool {
	value := *flg.ParsedValue

	for _, validator := range flg.Validators {
		if !validator(value) {
			fmt.Fprintf(os.Stderr, "Error: --%s/--%s argument is invalid\n", flg.LongName, flg.ShortName)
			os.Exit(enums.ExitCodeFlagError.Value())
		}
	}

	return true
}
