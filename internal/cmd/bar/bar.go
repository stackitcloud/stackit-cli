package bar

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(_ *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bar",
		Short: "Provides functionality for SKE",
		Long:  "Provides functionality for STACKIT Kubernetes Engine (SKE).",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	return cmd
}
