package delete

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
)

const distributionIDFlag = "distribution-id"

type inputModel struct {
	*globalflags.GlobalFlagModel
	DistributionID string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a CDN distribution",
		Long:  "Delete a CDN distribution by its ID.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Delete a CDN distribution with ID "xxx"`,
				`$ stackit beta cdn distribution delete --distribution-id xxx`,
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
	//TODO
}
