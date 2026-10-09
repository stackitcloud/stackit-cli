package alpha

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/sca"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/vpc"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alpha",
		Short: "Contains alpha STACKIT CLI commands",
		Long: fmt.Sprintf("%s\n%s",
			"Contains alpha STACKIT CLI commands.",
			"The commands under this group are still in an alpha state, and functionality may be incomplete or have breaking changes."),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
		Example: examples.Build(
			examples.NewExample(
				"See the currently available alpha commands",
				"$ stackit alpha --help"),
			examples.NewExample(
				"Execute an alpha command",
				"$ stackit alpha MY_COMMAND"),
		),
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(sca.NewCmd(params))
	cmd.AddCommand(vpc.NewCmd(params))
}
