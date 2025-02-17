package auth

import (
	activateserviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/auth/activate-service-account"
	getaccesstoken "github.com/stackitcloud/stackit-cli/internal/cmd/auth/get-access-token"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/login"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/logout"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticates the STACKIT CLI",
		Long:  "Authenticates in the STACKIT CLI.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(login.NewCmd(p))
	cmd.AddCommand(logout.NewCmd(p))
	cmd.AddCommand(activateserviceaccount.NewCmd(p))
	cmd.AddCommand(getaccesstoken.NewCmd(p))
}
