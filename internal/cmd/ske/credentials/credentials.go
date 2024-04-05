package credentials

import (
	completerotation "github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/complete-rotation"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/describe"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/rotate"
	startrotation "github.com/stackitcloud/stackit-cli/internal/cmd/ske/credentials/start-rotation"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credentials",
		Short: "Provides functionality for SKE credentials",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE) credentials.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd)
	return cmd
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(describe.NewCmd())
	cmd.AddCommand(rotate.NewCmd())
	cmd.AddCommand(startrotation.NewCmd())
	cmd.AddCommand(completerotation.NewCmd())
}
