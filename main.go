package main

import (
	"os"

	"github.com/stackitcloud/stackit-cli/internal/cmd"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

// These values are overwritten by GoReleaser at build time
var (
	version = "DEV"
	date    = "UNKNOWN"
)

func main() {
	// Set up configuration files
	config.InitConfig()

	printer := print.NewPrinter(
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)
	params := types.CmdParams{
		Printer:    printer,
		CliVersion: version,
		Date:       date,
		Fs:         os.DirFS("/"),
		Args:       os.Args[1:],
	}
	if !cmd.Execute(&params) {
		os.Exit(1)
	}
}
