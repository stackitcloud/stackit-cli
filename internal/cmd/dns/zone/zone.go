package zone

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone/update"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "zone",
		Short:   "Provides functionality for DNS zone",
		Long:    "Provides functionality for DNS zone",
		Example: fmt.Sprintf("%s\n%s", list.NewCmd().Example, create.NewCmd().Example),
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(update.NewCmd())
	cmd.AddCommand(delete.NewCmd())
}
