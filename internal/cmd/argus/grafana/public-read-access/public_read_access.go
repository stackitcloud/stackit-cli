package publicreadaccess

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/grafana/public-read-access/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/grafana/public-read-access/enable"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "public-read-access",
		Short: "Enable or disable public read access for Grafana in Argus instances",
		Long: fmt.Sprintf("%s\n%s",
			"Enable or disable public read access for Grafana in Argus instances.",
			"When enabled, anyone can access the Grafana dashboards of the instance without logging in. Otherwise, a login is required.",
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
