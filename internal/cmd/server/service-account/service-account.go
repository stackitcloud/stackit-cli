package serviceaccount

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/server/service-account/attach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/service-account/detach"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server/service-account/list"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service-account",
		Short: "Allows attaching/detaching service accounts to servers",
		Long:  "Allows attaching/detaching service accounts to servers",
		Args:  cobra.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(attach.NewCmd(p))
	cmd.AddCommand(detach.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
}
