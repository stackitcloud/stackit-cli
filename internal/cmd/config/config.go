package config

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/config/inspect"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/unset"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "config",
	Short:   "CLI configuration options",
	Long:    "CLI configuration options",
	Example: fmt.Sprintf("%s\n%s\n%s", set.Cmd.Example, inspect.Cmd.Example, unset.Cmd.Example),
}

func init() {
	Cmd.AddCommand(inspect.Cmd)
	Cmd.AddCommand(set.Cmd)
	Cmd.AddCommand(unset.Cmd)
}
