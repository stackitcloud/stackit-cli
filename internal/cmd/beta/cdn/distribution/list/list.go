package list

import (
	"context"
	"fmt"
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
	"github.com/stackitcloud/stackit-sdk-go/services/cdn"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SortBy string
}

const (
	sortByFlag = "sort-by"
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
				`$ stackit beta dns distribution list --sort-by=id`,
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

var sortByFlagOptions = []string{"id", "created", "updated", "origin-url", "status"}

func configureFlags(cmd *cobra.Command) {
	// same default as apiClient
	cmd.Flags().Var(flags.EnumFlag(false, "created", sortByFlagOptions...), sortByFlag, fmt.Sprintf("Sort entries by a specific field, one of %q", sortByFlagOptions))
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SortBy:          flags.FlagWithDefaultToStringValue(p, cmd, sortByFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *cdn.APIClient, nextPageID cdn.ListDistributionsResponseGetNextPageIdentifierAttributeType) cdn.ApiListDistributionsRequest {
	req := apiClient.ListDistributions(ctx, model.GlobalFlagModel.ProjectId)
	req = req.SortBy(toAPISortBy(model.SortBy))
	req = req.PageSize(100)
	if nextPageID != nil {
		req = req.PageIdentifier(*nextPageID)
	}
	return req
}

func toAPISortBy(sortBy string) string {
	switch sortBy {
	case "id":
		return "id"
	case "created":
		return "createdAt"
	case "updated":
		return "updatedAt"
	case "origin-url":
		return "originUrl"
	case "status":
		return "status"
	default:
		panic("invalid sortBy value, programmer error")
	}
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
			regions := make([]string, 0, len(*d.Config.Regions))
			for _, r := range *d.Config.Regions {
				regions = append(regions, string(r))
			}
			joinedRegions := strings.Join(regions, ", ")
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
	for {
		request := buildRequest(ctx, model, apiClient, nextPageID)
		response, err := request.Execute()
		if err != nil {
			return nil, fmt.Errorf("list distributions: %w", err)
		}
		nextPageID = response.NextPageIdentifier
		if response.Distributions != nil {
			distributions = append(distributions, *response.Distributions...)
		}
		if nextPageID == nil {
			break
		}
	}
	return distributions, nil
}
