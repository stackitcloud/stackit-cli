package targetpool

import (
	addtarget "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/target-pool/add-target"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/target-pool/describe"
	removetarget "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/target-pool/remove-target"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target-pool",
		Short: "Provides functionality for target pools",
		Long:  "Provides functionality for target pools.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(addtarget.NewCmd(params))
	cmd.AddCommand(removetarget.NewCmd(params))
	cmd.AddCommand(describe.NewCmd(params))
}
