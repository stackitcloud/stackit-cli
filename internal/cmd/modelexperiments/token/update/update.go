package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

const (
	tokenIDArg      = "TOKEN_ID"
	instanceIDFlag  = "instance-id"
	nameFlag        = "name"
	descriptionFlag = "description"
	labelsFlag      = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TokenID     string
	InstanceID  string
	Name        *string
	Description *string
	Labels      *map[string]*string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", tokenIDArg),
		Short: "Updates an AI Model Experiments instance token",
		Long:  "Partially updates an auth token for an AI Model Experiments instance.",
		Args:  args.SingleArg(tokenIDArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the name of an auth token with ID "xxx"`,
				`$ stackit model-experiments token update xxx --instance-id yyy --name updated-token`,
			),
			examples.NewExample(
				`Update the description of an auth token`,
				`$ stackit model-experiments token update xxx --instance-id yyy --description "updated token"`,
			),
			examples.NewExample(
				`Update labels on an auth token`,
				`$ stackit model-experiments token update xxx --instance-id yyy --labels env=prod,team=ml`,
			),
			examples.NewExample(
				`Update multiple fields of an auth token`,
				`$ stackit model-experiments token update xxx --instance-id yyy --name updated-token --description "updated token" --labels env=prod,team=ml`,
			),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, inputArgs)
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
			tokenLabel, err := meUtils.GetTokenName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceID, model.TokenID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get token name: %v", err)
				tokenLabel = model.TokenID
			}
			prompt := fmt.Sprintf("Are you sure you want to update token %q for instance %q in project %q?", tokenLabel, instanceLabel, projectLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("update instance token: %w", err)
			}
			return outputResult(params.Printer, model, resp)
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
	err := flags.MarkFlagsRequired(cmd, instanceIDFlag)
	cobra.CheckErr(err)
	cmd.MarkFlagsOneRequired(nameFlag, descriptionFlag, labelsFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	labels := flags.FlagToStringToAny(p, cmd, labelsFlag)
	var labelsPtr *map[string]*string
	if len(labels) > 0 {
		convertedLabels := make(map[string]*string, len(labels))
		for key, value := range labels {
			valueCopy := fmt.Sprint(value)
			convertedLabels[key] = &valueCopy
		}
		labelsPtr = &convertedLabels
	}
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		TokenID:         inputArgs[0],
		InstanceID:      flags.FlagToStringValue(p, cmd, instanceIDFlag),
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          labelsPtr,
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiPartialUpdateInstanceTokenRequest {
	payload := modelexperiments.PartialUpdateInstanceTokenPayload{
		Name:        model.Name,
		Description: model.Description,
		Labels:      model.Labels,
	}
	return apiClient.DefaultAPI.PartialUpdateInstanceToken(ctx, model.ProjectId, model.Region, model.TokenID, model.InstanceID).PartialUpdateInstanceTokenPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, resp *modelexperiments.PartialUpdateInstanceTokenResponse) error {
	if resp == nil {
		return fmt.Errorf("response token is nil")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Updated instance token %q (ID: %s)\n", resp.Token.Name, resp.Token.Id)
		return nil
	})
}
