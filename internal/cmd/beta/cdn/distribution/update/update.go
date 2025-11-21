package update

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

const (
	distributionIDFlag = "distribution-id"
	regionsFlag        = "regions"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	DistributionID string
	Regions        []cdn.Region
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a CDN distribution",
		Long:  "Update a CDN distribution by its ID, allowing replacement of its regions.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Update a CDN distribution with ID "xxx" to be cached in reions "EU" and "AF"`,
				`$ stackit beta cdn distribution update --distribution-id xxx --regions EU,AF`,
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
