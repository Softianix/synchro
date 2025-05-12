package validators

import (
	"fmt"
	"os"

	"github.com/Devourian/synchro/internal/enums"
)

// IsDir reports whether the given path exists and is a directory.
// Returns true if path exists and is a directory; false otherwise.
func IsDir(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(enums.ExitCodeFlagError.Value())
	}

	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "Error: %s is not a directory\n", dirPath)
		os.Exit(enums.ExitCodeFlagError.Value())
	}

	return true
}
