package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"
	"github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	displayNameFlag = "name"
	descriptionFlag = "description"
	planIdFlag      = "plan-id" // UUID
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	DisplayName string
	Description *string
	PlanId      string
}

// NewCmd https://aip.stackit.cloud/aip/general/0121/
// We have decided to eliminate the usage of display name flag
// To be the AIP compliant, and align with the standard CLI implementation, we will use the InstanceID arg
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an Edge Cloud instance",
		Long:  "Creates a STACKIT Edge Cloud (STEC) instance. The instance will take a moment to become fully functional.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Creates an Edge Cloud instance with the name "xxx" and plan-id "yyy"`,
				`$ stackit beta edge-cloud instance create --name "xxx" --plan-id "yyy"`),
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

			prompt := fmt.Sprintf("Are you sure you want to create a new edge instance for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			resp, err := request.Execute()
			if err != nil {
				return fmt.Errorf("create edge cloud instance: %w", err)
			}

			if !model.Async {
				err := spinner.Run(params.Printer, "Creating instance", func() error {
					_, err = wait.CreateOrUpdateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, resp.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for edge instance creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(displayNameFlag, "n", "", "The displayed name to distinguish multiple instances.")
	cmd.Flags().StringP(descriptionFlag, "d", "", "A user chosen description to distinguish multiple instances.")
	cmd.Flags().String(planIdFlag, "", "Service plan configures the size of the Instance.")

	err := flags.MarkFlagsRequired(cmd, displayNameFlag, planIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	displayNameValue := flags.FlagToStringValue(p, cmd, displayNameFlag)
	planIdValue := flags.FlagToStringValue(p, cmd, planIdFlag)
	descriptionValue := flags.FlagToStringPointer(p, cmd, descriptionFlag)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DisplayName:     displayNameValue,
		Description:     descriptionValue,
		PlanId:          planIdValue,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *edge.APIClient) edge.ApiCreateInstanceRequest {
	req := apiClient.DefaultAPI.CreateInstance(ctx, model.ProjectId, model.Region)

	payload := edge.CreateInstancePayload{
		DisplayName: model.DisplayName,
		Description: model.Description,
		PlanId:      model.PlanId,
	}

	return req.CreateInstancePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, instance *edge.Instance) error {
	return p.OutputResult(model.OutputFormat, instance, func() error {
		if instance == nil {
			return fmt.Errorf("instance response is empty")
		}

		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %q.\n", operationState, projectLabel, instance.Id)
		return nil
	})
}
