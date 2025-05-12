package main

import (
	"os"

	"github.com/Devourian/synchro/internal/core"
	"github.com/Devourian/synchro/internal/logging"
)

var log = logging.GetLogger()

// main just runs the Run function and passes the command line arguments to it.
func main() {
	log.Info("synchro started")
	os.Exit(core.Run(os.Args))
}
