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

var Cmd = &cobra.Command{
	Use:     "zone",
	Short:   "Provides functionality for DNS zone",
	Long:    "Provides functionality for DNS zone",
	Example: fmt.Sprintf("%s\n%s", list.Cmd.Example, create.Cmd.Example),
}

func init() {
	// Add all direct child commands
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(describe.Cmd)
	Cmd.AddCommand(update.Cmd)
	Cmd.AddCommand(delete.Cmd)
}
