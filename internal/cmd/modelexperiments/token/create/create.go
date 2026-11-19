package create

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
	meUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api/wait"
)

const (
	instanceIDFlag  = "instance-id"
	nameFlag        = "name"
	descriptionFlag = "description"
	labelsFlag      = "labels"
	ttlDurationFlag = "ttl-duration"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceID  string
	Name        string
	Description *string
	Labels      *map[string]string
	TTLDuration *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an AI Model Experiments instance token",
		Long:  "Creates an auth token for an AI Model Experiments instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				"Create an auth token",
				"$ stackit model-experiments token create --instance-id xxx --name my-token",
			),
			examples.NewExample(
				"Create an auth token with a description, labels, and time to live",
				`$ stackit model-experiments token create --instance-id xxx --name my-token --description "service token" --labels env=prod,team=ml --ttl-duration 30d`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, nil)
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
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}
			instanceLabel, err := meUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceID
			}
			prompt := fmt.Sprintf("Are you sure you want to create token %q for instance %q in project %q?", model.Name, instanceLabel, projectLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("create instance token: %w", err)
			}
			if resp == nil || resp.Token.Id == "" {
				return fmt.Errorf("create instance token: response token is nil")
			}
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating instance token", func() error {
					_, err = wait.CreateInstanceTokenWaitHandler(ctx, apiClient.DefaultAPI, model.Region, model.ProjectId, model.InstanceID, resp.Token.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for instance token creation: %w", err)
				}
			}
			return outputResult(params.Printer, model, projectLabel, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIDFlag, "ID of the instance")
	cmd.Flags().String(nameFlag, "", "Token name")
	cmd.Flags().String(descriptionFlag, "", "Token description")
	cmd.Flags().StringToString(labelsFlag, nil, `Labels are comma-separated key-value string pairs, e.g. "--labels env=prod,team=ml"`)
	cmd.Flags().String(ttlDurationFlag, "", `Maximum allowed duration is 6 months. If not set, defaults to 6 months. Format example: "720h" or "30d"`)
	err := flags.MarkFlagsRequired(cmd, instanceIDFlag, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	labels := flags.FlagToStringToStringPointer(p, cmd, labelsFlag)
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceID:      flags.FlagToStringValue(p, cmd, instanceIDFlag),
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          labels,
		TTLDuration:     flags.FlagToStringPointer(p, cmd, ttlDurationFlag),
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiCreateInstanceTokenRequest {
	payload := modelexperiments.CreateInstanceTokenPayload{
		Name:        model.Name,
		Description: model.Description,
		Labels:      model.Labels,
		TtlDuration: model.TTLDuration,
	}
	return apiClient.DefaultAPI.CreateInstanceToken(ctx, model.ProjectId, model.Region, model.InstanceID).CreateInstanceTokenPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel, instanceLabel string, resp *modelexperiments.CreateInstanceTokenResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("response token is nil")
		}

		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s token %q for instance %q in project %q. Token ID: %s\nToken: %s\n", operationState, resp.Token.Name, instanceLabel, projectLabel, resp.Token.Id, resp.Token.Content)
		return nil
	})
}
