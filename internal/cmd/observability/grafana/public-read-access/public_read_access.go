package publicreadaccess

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana/public-read-access/disable"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability/grafana/public-read-access/enable"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "public-read-access",
		Short: "Enable or disable public read access for Grafana in Observability instances",
		Long: fmt.Sprintf("%s\n%s",
			"Enable or disable public read access for Grafana in Observability instances.",
			"When enabled, anyone can access the Grafana dashboards of the instance without logging in. Otherwise, a login is required.",
		),
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(enable.NewCmd(params))
	cmd.AddCommand(disable.NewCmd(params))
}
