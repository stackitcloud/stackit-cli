package main

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

// These values are overwritten by GoReleaser at build time
var (
	version = "DEV"
	date    = "UNKNOWN"
)

func main() {
	// Set up configuration files
	config.InitConfig()

	cmd.Execute(version, date)
}
