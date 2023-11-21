package recordset

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set/update"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "record-set",
		Short:   "Provides functionality for DNS record set",
		Long:    "Provides functionality for DNS record set",
		Example: fmt.Sprintf("%s\n%s", list.NewCmd().Example, create.NewCmd().Example),
	}
	addChilds(cmd)
	return cmd
}

func addChilds(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(delete.NewCmd())
	cmd.AddCommand(update.NewCmd())
}
