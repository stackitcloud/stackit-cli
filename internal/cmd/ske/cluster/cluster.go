package cluster

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/list"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "zone",
	Short:   "Provides functionality for SKE cluster",
	Long:    "Provides functionality for SKE cluster",
	Example: list.Cmd.Example,
}

func init() {
	Cmd.AddCommand(list.Cmd)
}
