package singlesignon

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana/single-sign-on/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana/single-sign-on/enable"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "single-sign-on",
		Aliases: []string{"sso"},
		Short:   "Enable or disable single sign-on for Grafana in Observability instances",
		Long: fmt.Sprintf("%s\n%s",
			"Enable or disable single sign-on for Grafana in Observability instances.",
			"When enabled for an instance, overwrites the generic OAuth2 authentication and configures STACKIT single sign-on for that instance.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(enable.NewCmd(p))
	cmd.AddCommand(disable.NewCmd(p))
}
