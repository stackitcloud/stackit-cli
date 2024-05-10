package config

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/config/list"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/profile"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config/unset"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Provides functionality for CLI configuration options",
		Long: fmt.Sprintf("%s\n%s\n\n%s\n%s\n%s",
			"Provides functionality for CLI configuration options.",
			`You can set and unset different configuration options via the "stackit config set" and "stackit config unset" commands.`,
			"Additionally, you can configure the CLI to use different profiles, each with its own configuration.",
			`Additional profiles can be configured via the "STACKIT_CLI_PROFILE" environment variable or using the "stackit config profile set PROFILE" and "stackit config profile unset" commands.`,
			"The environment variable takes precedence over what is set via the commands.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(list.NewCmd(p))
	cmd.AddCommand(set.NewCmd(p))
	cmd.AddCommand(unset.NewCmd(p))
	cmd.AddCommand(profile.NewCmd(p))
}
