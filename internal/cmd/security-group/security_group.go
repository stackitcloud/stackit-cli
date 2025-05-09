package security_group

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/rule"
	"github.com/stackitcloud/stackit-cli/internal/cmd/security-group/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "security-group",
		Short: "Manage security groups",
		Long:  "Manage the lifecycle of security groups and rules.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *params.CmdParams) {
	cmd.AddCommand(
		rule.NewCmd(params),
		create.NewCmd(params),
		delete.NewCmd(params),
		describe.NewCmd(params),
		list.NewCmd(params),
		update.NewCmd(params),
	)
}
