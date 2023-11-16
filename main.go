package main

import (
	"stackit/internal/cmd"
	"stackit/internal/pkg/config"
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
