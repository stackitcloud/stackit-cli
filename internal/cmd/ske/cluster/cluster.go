package cluster

import (

	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/list"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "cluster",
	Short:   "Provides functionality for SKE cluster",
	Long:    "Provides functionality for SKE cluster",
	Example: list.Cmd.Example,
}

func init() {
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(describe.Cmd)
	Cmd.AddCommand(delete.Cmd)
}
