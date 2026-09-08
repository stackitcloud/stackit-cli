package patch

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
	tokenIDArg      = "TOKEN_ID"
	instanceIDFlag  = "instance-id"
	regionFlag      = "region"
	nameFlag        = "name"
	descriptionFlag = "description"
	labelFlag       = "label"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TokenID     string
	InstanceID  string
	Region      string
	Name        *string
	Description *string
	Labels      *map[string]string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use: fmt.Sprintf("patch %s", tokenIDArg), Short: "Updates an instance token", Long: "Partially updates an auth token for an AI Model Experiments instance.", Args: args.SingleArg(tokenIDArg, utils.ValidateUUID),
		Example: examples.Build(examples.NewExample(`Update an auth token with ID "xxx"`, `$ stackit ai-model-experiments token patch xxx --instance-id yyy --name updated-token`)),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return err
			}
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}
			resp, err := buildRequest(context.Background(), model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("update instance token: %w", err)
			}
			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIDFlag, "ID of the instance")
	cmd.Flags().String(nameFlag, "", "Token name")
	cmd.Flags().String(descriptionFlag, "", "Token description")
	cmd.Flags().StringToString(labelFlag, nil, `Labels as key-value pairs, e.g. "--label env=prod"`)
	_ = flags.MarkFlagsRequired(cmd, instanceIDFlag, regionFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}
	labels, err := cmd.Flags().GetStringToString(labelFlag)
	if err != nil {
		return nil, fmt.Errorf("parse %q flag: %w", labelFlag, err)
	}
	var labelsPtr *map[string]string
	if len(labels) > 0 {
		labelsPtr = &labels
	}
	model := &inputModel{GlobalFlagModel: globalFlags, TokenID: inputArgs[0], InstanceID: flags.FlagToStringValue(p, cmd, instanceIDFlag), Region: flags.FlagToStringValue(p, cmd, regionFlag), Name: flags.FlagToStringPointer(p, cmd, nameFlag), Description: flags.FlagToStringPointer(p, cmd, descriptionFlag), Labels: labelsPtr}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiPartialUpdateInstanceTokenRequest {
	labels := map[string]*string(nil)
	if model.Labels != nil {
		labels = make(map[string]*string, len(*model.Labels))
		for key, value := range *model.Labels {
			labels[key] = &value
		}
	}
	payload := modelexperiments.PartialUpdateInstanceTokenPayload{Name: model.Name, Description: model.Description, Labels: &labels}
	return apiClient.DefaultAPI.PartialUpdateInstanceToken(ctx, model.ProjectId, model.Region, model.TokenID, model.InstanceID).PartialUpdateInstanceTokenPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, resp *modelexperiments.PartialUpdateInstanceTokenResponse) error {
	if resp == nil || resp.Token.Name == "" {
		return fmt.Errorf("response token is nil")
	}
	return p.OutputResult(outputFormat, resp.Token, func() error { return nil })
}
