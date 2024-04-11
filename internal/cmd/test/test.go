package test

import (
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "test",
		Long:  "test",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	return cmd
}
