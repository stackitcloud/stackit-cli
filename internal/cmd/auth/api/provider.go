package api

import (
	"github.com/spf13/cobra"
	getaccesstoken "github.com/stackitcloud/stackit-cli/internal/cmd/auth/api/get-access-token"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/api/login"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/api/logout"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth/api/status"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Manages authentication for the STACKIT Terraform Provider and SDK",
		Long: `Manages authentication for the STACKIT Terraform Provider and SDK.

These commands allow you to authenticate with your personal STACKIT account
and share the credentials with the STACKIT Terraform Provider and SDK.
This provides an alternative to using service accounts for local development.

Tokens are stored separately from tokens received from "stackit auth login", this allows using separate accounts when using the cli directly vs. using it to manage credentials for the SDK and TF provider.

Tokens are stored in the OS keychain with a fallback to local storage.`,
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *types.CmdParams) {
	cmd.AddCommand(login.NewCmd(p))
	cmd.AddCommand(logout.NewCmd(p))
	cmd.AddCommand(getaccesstoken.NewCmd(p))
	cmd.AddCommand(status.NewCmd(p))
}
