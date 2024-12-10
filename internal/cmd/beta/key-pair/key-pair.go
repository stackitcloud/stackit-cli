package keypair

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/key-pair/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/key-pair/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/key-pair/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/key-pair/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/key-pair/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-pair",
		Short: "Provides functionality for SSH Key Pairs",
		Long:  "Provides functionality for SSH Key Pairs",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(delete.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(update.NewCmd(p))
}
