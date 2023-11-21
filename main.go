package main

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
)

func main() {
	// Set up configuration files
	config.InitConfig()

	cmd.Execute()
}
