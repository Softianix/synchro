package enums

type ExitCode int

const (
	ExitCodeSuccess   ExitCode = iota // 0 - Success
	ExitCodeError                     // 1 - General error
	ExitCodeFlagError                 // 2 - Error with command line flags
)

// ExitCode.Value returns the integer value of the ExitCode,
// analogous to Python’s Enum.value.
func (exitCode ExitCode) Value() int {
	return int(exitCode)
}
