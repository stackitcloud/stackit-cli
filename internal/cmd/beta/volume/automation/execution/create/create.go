package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	automation "github.com/stackitcloud/stackit-sdk-go/services/automation/v1betaapi"
	"github.com/stackitcloud/stackit-sdk-go/services/automation/v1betaapi/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/automation/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	automationIdFlag = "automation-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AutomationId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new Volume Automation Execution",
		Long:  "Creates a new Volume Automation Execution.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new Volume Automation Execution for Automation with ID "xxx"`,
				"$ stackit beta volume automation execution create --automation-id xxx"),
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

			// Call API
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("create volume automation execution: %w", err)
			}

			if !model.Async {
				if resp == nil {
					return fmt.Errorf("create volume automation execution response is empty")
				}
				err = spinner.Run(params.Printer, "Creating automation execution", func() error {
					_, err := wait.CreateVolumeExecutionWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.AutomationId, resp.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for automation execution: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, model.AutomationId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), automationIdFlag, "ID of the automation which should be executed")

	err := flags.MarkFlagsRequired(cmd, automationIdFlag)
	cobra.CheckErr(err)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *automation.APIClient) automation.ApiCreateVolumeExecutionRequest {
	return apiClient.DefaultAPI.CreateVolumeExecution(ctx, model.ProjectId, model.Region, model.AutomationId)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		AutomationId:    flags.FlagToStringValue(p, cmd, automationIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, async bool, automationId string, execution *automation.VolumeExecutionResponse) error {
	return p.OutputResult(outputFormat, execution, func() error {
		if execution == nil {
			return fmt.Errorf("create volume automation execution response is empty")
		}

		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s volume automation execution for automation %q. Execution ID %q.\n", operationState, automationId, execution.Id)
		return nil
	})
}
