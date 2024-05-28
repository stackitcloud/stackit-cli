package profile

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/create"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile/unset"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Manage the CLI configuration profiles",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Manage the CLI configuration profiles.",
			`The profile to be used can be managed via the "STACKIT_CLI_PROFILE" environment variable or using the "stackit config profile set PROFILE" and "stackit config profile unset" commands.`,
			"The environment variable takes precedence over what is set via the commands.",
			"When no profile is set, the default profile is used.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(set.NewCmd(p))
	cmd.AddCommand(unset.NewCmd(p))
	cmd.AddCommand(create.NewCmd(p))
	cmd.AddCommand(list.NewCmd(p))
}
