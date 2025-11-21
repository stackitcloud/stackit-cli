package create

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
)

const (
	flagDistributionID = "distribution-id"
	flagName           = "name"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a CDN domain",
		Long:  "Create a new CDN domain associated with a CDN distribution.",
		Args:  cobra.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a CDN domain named "example.com" for distribution with ID "xxx"`,
				`$ stackit beta cdn domain create --name example.com --distribution-id xxx`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	// TODO
}
