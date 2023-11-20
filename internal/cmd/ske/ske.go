package ske

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "ske",
	Short:   "Provides functionality for SKE",
	Long:    "Provides functionality for STACKIT Kubernetes engine (SKE)",
	Example: cluster.Cmd.Example,
}

func init() {
	Cmd.AddCommand(cluster.Cmd)
}
