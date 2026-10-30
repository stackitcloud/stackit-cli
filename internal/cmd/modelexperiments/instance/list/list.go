package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

const limitFlag = "limit"

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all AI Model Experiments instances",
		Long:  "Lists all AI Model Experiments instances in a STACKIT project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List AI Model Experiments instances in a project`,
				`$ stackit model-experiments instance list`,
			),
			examples.NewExample(
				`List AI Model Experiments instances in a project in region eu02`,
				`$ stackit model-experiments instance list --region eu02`,
			),
			examples.NewExample(
				`List AI Model Experiments instances in JSON format`,
				`$ stackit model-experiments instance list --output-format json`,
			),
			examples.NewExample(
				`List up to 10 AI Model Experiments instances`,
				`$ stackit model-experiments instance list --limit 10`,
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

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list AI Model Experiments instances: %w", err)
			}

			applyLimit(resp, model.Limit)

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}

	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	return cmd
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiListInstancesRequest {
	return apiClient.DefaultAPI.ListInstances(ctx, model.ProjectId, model.Region)
}

func applyLimit(resp *modelexperiments.ListInstancesResponse, limit *int64) {
	if resp == nil || limit == nil || int64(len(resp.Instances)) <= *limit {
		return
	}

	resp.Instances = resp.Instances[:int(*limit)]
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *modelexperiments.ListInstancesResponse) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(outputFormat, resp.Instances, func() error {
		if len(resp.Instances) == 0 {
			p.Outputf("No instances found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "ARTIFACT STORAGE NAME", "REGION", "STATUS")

		for _, instance := range resp.Instances {
			table.AddRow(
				instance.Id,
				instance.Name,
				instance.GetBucketName(),
				utils.PtrString(instance.Region),
				instance.State,
			)
		}

		return table.Display(p)
	})
}
