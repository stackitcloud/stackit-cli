package create

import (
	"context"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/cdn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

const (
	regionsFlag   = "regions"
	originURLFlag = "origin-url"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Regions   []cdn.Region
	OriginURL string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a CDN distribution",
		Long:  "Create a CDN distribution for a given originUrl in multiple regions.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a distribution for regions EU and AF`,
				`$ stackit beta cdn distribution create --regions=EU,AF --origin-url=https://example.com`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a CDN distribution for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create CDN distribution: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.EnumSliceFlag(false, []string{}, sdkUtils.EnumSliceToStringSlice(cdn.AllowedRegionEnumValues)...), regionsFlag, fmt.Sprintf("Regions in which content should be cached, multiple of: %q", cdn.AllowedRegionEnumValues))
	cmd.Flags().String(originURLFlag, "", "The origin of the content that should be made available through the CDN. Note that the path and query parameters are ignored. Ports are allowed. If no protocol is provided, `https` is assumed. So `www.example.com:1234/somePath?q=123` is normalized to `https://www.example.com:1234`")
	err := flags.MarkFlagsRequired(cmd, regionsFlag, originURLFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	regionStrings := flags.FlagToStringSliceValue(p, cmd, regionsFlag)
	regions := make([]cdn.Region, 0, len(regionStrings))
	for _, regionStr := range regionStrings {
		regions = append(regions, cdn.Region(regionStr))
	}

	originUrlString := flags.FlagToStringValue(p, cmd, originURLFlag)
	_, err := url.Parse(originUrlString)
	if err != nil {
		return nil, fmt.Errorf("invalid originUrl: '%s' (%w)", originUrlString, err)
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Regions:         regions,
		OriginURL:       originUrlString,
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *cdn.APIClient) cdn.ApiCreateDistributionRequest {
	req := apiClient.CreateDistribution(ctx, model.ProjectId)
	payload := cdn.NewCreateDistributionPayload(
		model.OriginURL,
		model.Regions,
	)
	return req.CreateDistributionPayload(*payload)
}

func outputResult(p *print.Printer, outputFormat string, projectLabel string, resp *cdn.CreateDistributionResponse) error {
	if resp == nil {
		return fmt.Errorf("create distribution response is nil")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created CDN distribution for %q. Id: %s\n", projectLabel, utils.PtrString(resp.Distribution.Id))
		return nil
	})
}
