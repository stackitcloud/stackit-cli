package describe

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
		Use:   "describe",
		Short: "Describe a CDN distribution",
		Long:  "Describe a CDN distribution by its ID.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get details of a CDN distribution with ID "xxx"`,
				`$ stackit beta cdn distribution describe --distribution-id xxx`,
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
	cmd.Flags().String(distributionIDFlag, "", "The ID of the CDN distribution to describe")
	err := cmd.MarkFlagRequired(distributionIDFlag)
	cobra.CheckErr(err)
}
