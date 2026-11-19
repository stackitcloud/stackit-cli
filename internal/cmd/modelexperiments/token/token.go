package token

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/token/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/token/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/token/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/token/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/modelexperiments/token/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Provides functionality for AI Model Experiments instance tokens",
		Long:  "Provides functionality for AI Model Experiments instance tokens.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	cmd.AddCommand(list.NewCmd(params))
	cmd.AddCommand(create.NewCmd(params))
	cmd.AddCommand(delete.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
	cmd.AddCommand(update.NewCmd(params))
	return cmd
}
