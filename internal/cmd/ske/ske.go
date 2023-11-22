package ske

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/describe"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ske",
		Short:   "Provides functionality for SKE",
		Long:    "Provides functionality for STACKIT Kubernetes engine (SKE)",
		Example: fmt.Sprintf("%s\n%s", describe.NewCmd().Example, cluster.NewCmd().Example),
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(cluster.NewCmd())
}
