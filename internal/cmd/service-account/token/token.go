package token

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/revoke"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Provides functionality for service account tokens",
		Long:  "Provides functionality for service account tokens.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(revoke.NewCmd(params))
}
