package api

import (
	"github.com/spf13/cobra"
	getaccesstoken "github.com/stackitcloud/stackit-cli/internal/cmd/auth/api/get-access-token"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/api/login"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/api/logout"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/api/status"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Manages authentication for the STACKIT Terraform Provider and SDK",
		Long: `Manages authentication for the STACKIT Terraform Provider and SDK.

These commands allow you to authenticate with your personal STACKIT account
and share the credentials with the STACKIT Terraform Provider and SDK.
This provides an alternative to using service accounts for local development.`,
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(login.NewCmd(params))
	cmd.AddCommand(logout.NewCmd(params))
	cmd.AddCommand(getaccesstoken.NewCmd(params))
	cmd.AddCommand(status.NewCmd(params))
}
