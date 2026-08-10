package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

// NewCmd https://aip.stackit.cloud/aip/general/0121/
// We have decided to eliminate the usage of display name flag
// To be the AIP compliant, and align with the standard CLI implementation, we will use the InstanceID arg
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists Edge Cloud instances",
		Long:  "Lists STACKIT Edge Cloud (STEC) instances of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all Edge Cloud instances of a given project`,
				`$ stackit beta edge-cloud instance list`),
			examples.NewExample(
				`Lists all Edge Cloud instances of a given project and limits the output to two instances`,
				`$ stackit beta edge-cloud instance list --limit 2`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Call API
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("get Edge Cloud instances: %w", err)
			}
			instances := resp.Instances

			// Truncate output
			if model.Limit != nil && len(instances) > int(*model.Limit) {
				instances = instances[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, instances)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &cliErr.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *edge.APIClient) edge.ApiListInstancesRequest {
	return apiClient.DefaultAPI.ListInstances(ctx, model.ProjectId, model.Region)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, instances []edge.Instance) error {
	return p.OutputResult(outputFormat, instances, func() error {
		if len(instances) == 0 {
			p.Outputf("No instances found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "UI", "STATE")
		for _, instance := range instances {
			table.AddRow(
				instance.Id,
				instance.DisplayName,
				instance.FrontendUrl,
				instance.Status)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
