package auth

import (
	activateserviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/auth/activate-service-account"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/api"
	getaccesstoken "github.com/stackitcloud/stackit-cli/internal/cmd/auth/get-access-token"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/login"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/logout"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticates the STACKIT CLI",
		Long:  "Authenticates in the STACKIT CLI.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(login.NewCmd(params))
	cmd.AddCommand(logout.NewCmd(params))
	cmd.AddCommand(activateserviceaccount.NewCmd(params))
	cmd.AddCommand(getaccesstoken.NewCmd(params))
	cmd.AddCommand(api.NewCmd(params))
}
