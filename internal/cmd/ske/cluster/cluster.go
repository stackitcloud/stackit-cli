package cluster

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/update"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cluster",
		Short:   "Provides functionality for SKE cluster",
		Long:    "Provides functionality for SKE cluster",
		Example: fmt.Sprintf("%s\n%s", create.NewCmd().Example, list.NewCmd().Example),
	}
	addChilds(cmd)
	return cmd
}

func addChilds(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(update.NewCmd())
}
