package delete

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
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
	Region     string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", instanceIdArg),
		Short: "Deletes an AI Model Experiments instance",
		Long:  "Deletes an AI Model Experiments instance from a STACKIT project.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete an AI Model Experiments instance with ID "xxx"`,
				`$ stackit ai-model-experiments instance delete xxx`,
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

			prompt := fmt.Sprintf(
				"Are you sure you want to delete AI Model Experiments instance %q?",
				model.InstanceId,
			)

			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildDeleteInstanceRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("delete AI Model Experiments instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	_ = flags.MarkFlagsRequired(cmd, globalflags.RegionFlag)
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
		Region:          flags.FlagToStringValue(p, cmd, globalflags.RegionFlag),
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildDeleteInstanceRequest(
	ctx context.Context,
	model *inputModel,
	apiClient *modelexperiments.APIClient,
) modelexperiments.ApiDeleteInstanceRequest {
	return apiClient.DefaultAPI.DeleteInstance(
		ctx,
		model.ProjectId,
		model.Region,
		model.InstanceId,
	)
}

func outputResult(
	p *print.Printer,
	outputFormat string,
	resp *modelexperiments.DeleteInstanceResponse,
) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(outputFormat, resp.Instance, func() error {
		p.Outputf(
			"Deleted AI Model Experiments instance.\n",
		)
		return nil
	})
}
