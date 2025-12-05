package list

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/cdn/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SortBy string
	Limit  *int32
}

const (
	sortByFlag  = "sort-by"
	limitFlag   = ""
	maxPageSize = int32(100)
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List CDN distributions",
		Long:  "List all CDN distributions in your account.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all CDN distributions`,
				`$ stackit beta cdn distribution list`,
			),
			examples.NewExample(
				`List all CDN distributions sorted by id`,
				`$ stackit beta cdn distribution list --sort-by=id`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background() // should this be cancellable?

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			distributions, err := fetchDistributions(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("fetch distributions: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, distributions)
		},
	}

	configureFlags(cmd)
	return cmd
}

var sortByFlagOptions = []string{"id", "createdAt", "updatedAt", "originUrl", "status", "originUrlRelated"}

func configureFlags(cmd *cobra.Command) {
	// same default as apiClient
	cmd.Flags().Var(flags.EnumFlag(false, "createdAt", sortByFlagOptions...), sortByFlag, fmt.Sprintf("Sort entries by a specific field, one of %q", sortByFlagOptions))
	cmd.Flags().Int64(limitFlag, 0, "Limit the output to the first n elements")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt32Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SortBy:          flags.FlagWithDefaultToStringValue(p, cmd, sortByFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *cdn.APIClient, nextPageID cdn.ListDistributionsResponseGetNextPageIdentifierAttributeType, pageLimit int32) cdn.ApiListDistributionsRequest {
	req := apiClient.ListDistributions(ctx, model.GlobalFlagModel.ProjectId)
	req = req.SortBy(model.SortBy)
	req = req.PageSize(pageLimit)
	if nextPageID != nil {
		req = req.PageIdentifier(*nextPageID)
	}
	return req
}

func outputResult(p *print.Printer, outputFormat string, distributions []cdn.Distribution) error {
	if distributions == nil {
		distributions = make([]cdn.Distribution, 0) // otherwise prints null in json output
	}
	return p.OutputResult(outputFormat, distributions, func() error {
		if len(distributions) == 0 {
			p.Outputln("No CDN distributions found")
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "REGIONS", "STATUS")
		for i := range distributions {
			d := &distributions[i]
			joinedRegions := strings.Join(sdkUtils.EnumSliceToStringSlice(*d.Config.Regions), ", ")
			table.AddRow(
				utils.PtrString(d.Id),
				joinedRegions,
				utils.PtrString(d.Status),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}

func fetchDistributions(ctx context.Context, model *inputModel, apiClient *cdn.APIClient) ([]cdn.Distribution, error) {
	var nextPageID cdn.ListDistributionsResponseGetNextPageIdentifierAttributeType
	var distributions []cdn.Distribution
	received := int32(0)
	limit := int32(math.MaxInt32)
	if model.Limit != nil {
		limit = min(limit, *model.Limit)
	}
	for {
		want := min(maxPageSize, limit-received)
		request := buildRequest(ctx, model, apiClient, nextPageID, want)
		response, err := request.Execute()
		if err != nil {
			return nil, fmt.Errorf("list distributions: %w", err)
		}
		if response.Distributions != nil {
			distributions = append(distributions, *response.Distributions...)
		}
		nextPageID = response.NextPageIdentifier
		received += want
		if nextPageID == nil || received >= limit {
			break
		}
	}
	return distributions, nil
}
