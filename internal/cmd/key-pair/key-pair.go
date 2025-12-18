package keypair

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/cmd/key-pair/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/key-pair/delete"
	"github.com/stackitcloud/stackit-cli/internal/cmd/key-pair/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/key-pair/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/key-pair/update"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "key-pair",
		Short: "Provides functionality for SSH key pairs",
		Long:  "Provides functionality for SSH key pairs",
		Args:  cobra.NoArgs,
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
