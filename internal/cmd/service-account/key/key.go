package key

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/service-account/key/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key",
		Short: "Provides functionality for service account keys",
		Long:  "Provides functionality for service account keys.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
}
