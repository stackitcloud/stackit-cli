package beta

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/alb"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/sqlserverflex"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
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
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(sqlserverflex.NewCmd(p))
	cmd.AddCommand(alb.NewCmd(p))
}
