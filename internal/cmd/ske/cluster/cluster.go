package cluster

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/list"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "cluster",
	Short:   "Provides functionality for SKE cluster",
	Long:    "Provides functionality for SKE cluster",
	Example: fmt.Sprintf("%s\n%s", create.Cmd.Example, list.Cmd.Example),
}

func init() {
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(delete.Cmd)
	Cmd.AddCommand(describe.Cmd)
	Cmd.AddCommand(list.Cmd)
}
