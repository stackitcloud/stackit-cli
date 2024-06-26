package grafana

import (
	"github.com/stackitcloud/stackit-cli/internal/cmd/argus/grafana/describe"
	publicreadaccess "github.com/stackitcloud/stackit-cli/internal/cmd/argus/grafana/public-read-access"
	singlesignon "github.com/stackitcloud/stackit-cli/internal/cmd/argus/grafana/single-sign-on"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grafana",
		Short: "Provides functionality for the Grafana configuration of Argus instances",
		Long:  "Provides functionality for the Grafana configuration of Argus instances.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, p)
	return cmd
}

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(describe.NewCmd(p))
	cmd.AddCommand(publicreadaccess.NewCmd(p))
	cmd.AddCommand(singlesignon.NewCmd(p))
}
