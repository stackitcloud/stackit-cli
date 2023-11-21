package ske

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ske",
		Short:   "Provides functionality for SKE",
		Long:    "Provides functionality for STACKIT Kubernetes engine (SKE)",
		Example: cluster.NewCmd().Example,
	}
	addChilds(cmd)
	return cmd
}

func addChilds(cmd *cobra.Command) {
	cmd.AddCommand(cluster.NewCmd())
}
