package beta

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/logs"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/edge"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/intake"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/kms"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sfs"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "beta",
		Short: "Contains beta STACKIT CLI commands",
		Long: fmt.Sprintf("%s\n%s",
			"Contains beta STACKIT CLI commands.",
			"The commands under this group are still in a beta state, and functionality may be incomplete or have breaking changes."),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
		Example: examples.Build(
			examples.NewExample(
				"See the currently available beta commands",
				"$ stackit beta --help"),
			examples.NewExample(
				"Execute a beta command",
				"$ stackit beta MY_COMMAND"),
		),
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(sqlserverflex.NewCmd(params))
	cmd.AddCommand(sfs.NewCmd(params))
	cmd.AddCommand(alb.NewCmd(params))
	cmd.AddCommand(edge.NewCmd(params))
	cmd.AddCommand(intake.NewCmd(params))
	cmd.AddCommand(kms.NewCmd(params))
	cmd.AddCommand(logs.NewCmd(params))
}
