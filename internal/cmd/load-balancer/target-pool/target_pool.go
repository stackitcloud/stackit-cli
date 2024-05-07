package targetpool

import (
	addtarget "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/target-pool/add-target"
	"github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/target-pool/describe"
	removetarget "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer/target-pool/remove-target"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target-pool",
		Short: "Provides functionality for target pools",
		Long:  "Provides functionality for target pools.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(addtarget.NewCmd(p))
	cmd.AddCommand(removetarget.NewCmd(p))
	cmd.AddCommand(describe.NewCmd(p))
}
