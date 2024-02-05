package token

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/token/revoke"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Provides functionality regarding service account tokens",
		Long:  "Provides functionality regarding service account tokens.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(create.NewCmd())
	cmd.AddCommand(list.NewCmd())
	cmd.AddCommand(revoke.NewCmd())
}
