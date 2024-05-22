package main

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

// These values are dynamically overridden by GoReleaser
var (
	version = "DEV"
	date    = "UNKNOWN"
)

func main() {
	// Set up configuration files
	config.InitConfig()

	cmd.Execute(version, date)
}
