package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of an AI Model Experiments instance",
		Long:  "Shows details of an AI Model Experiments instance in a STACKIT project.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of an AI Model Experiments instance with ID "xxx"`,
				`$ stackit model-experiments instance describe xxx`,
			),
			examples.NewExample(
				`Get details of an AI Model Experiments instance with ID "xxx" in JSON format`,
				`$ stackit ai-model-experiments instance get xxx --output-format json`,
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
				return fmt.Errorf("describe AI Model Experiments instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	return cmd
}

func parseInput(
	p *print.Printer,
	cmd *cobra.Command,
	inputArgs []string,
) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      inputArgs[0],
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildRequest(
	ctx context.Context,
	model *inputModel,
	apiClient *modelexperiments.APIClient,
) modelexperiments.ApiGetInstanceRequest {
	return apiClient.DefaultAPI.GetInstance(
		ctx,
		model.ProjectId,
		model.Region,
		model.InstanceId,
	)
}

func outputResult(
	p *print.Printer,
	outputFormat string,
	resp *modelexperiments.GetInstanceResponse,
) error {
	if resp == nil {
		return fmt.Errorf("response instance is nil")
	}

	return p.OutputResult(outputFormat, resp.Instance, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "REGION", "STATUS")
		table.AddRow(
			resp.Instance.Id,
			resp.Instance.Name,
			utils.PtrString(resp.Instance.Region),
			resp.Instance.State,
		)
		return table.Display(p)
	})
}
