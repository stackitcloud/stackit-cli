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

var Cmd = &cobra.Command{
	Use:     "record-set",
	Short:   "Provides functionality for DNS record set",
	Long:    "Provides functionality for DNS record set",
	Example: fmt.Sprintf("%s\n%s", list.Cmd.Example, create.Cmd.Example),
}

func init() {
	// Add all direct child commands
	Cmd.AddCommand(list.Cmd)
	Cmd.AddCommand(create.Cmd)
	Cmd.AddCommand(describe.Cmd)
	Cmd.AddCommand(delete.Cmd)
	Cmd.AddCommand(update.Cmd)
}
