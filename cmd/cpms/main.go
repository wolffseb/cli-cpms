// Command cpms is a terminal charge point management system: it speaks OCPP to
// a charging station and exposes a minimal OCPI 2.3.0 CPO interface on the LAN.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/wolffseb/cli-cpms/internal/cli"
)

func main() {
	if err := cli.NewRootCommand().Execute(); err != nil {
		// ErrSilent means the command already printed a full report.
		if !errors.Is(err, cli.ErrSilent) {
			fmt.Fprintln(os.Stderr, "cpms: "+err.Error())
		}
		os.Exit(1)
	}
}
