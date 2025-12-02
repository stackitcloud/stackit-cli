package secretsmanager

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager/user"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets-manager",
		Short: "Provides functionality for Secrets Manager",
		Long:  "Provides functionality for Secrets Manager.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(user.NewCmd(params))
}
