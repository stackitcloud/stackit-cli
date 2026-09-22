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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/modelexperiments/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

const (
	instanceIDFlag  = "instance-id"
	nameFlag        = "name"
	descriptionFlag = "description"
	labelFlag       = "label"
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
	cmd := &cobra.Command{Use: "create", Short: "Creates an instance token", Long: "Creates an auth token for an AI Model Experiments instance.", Args: args.NoArgs, Example: examples.Build(examples.NewExample("Create an auth token", "$ stackit ai-model-experiments token create --instance-id xxx --name my-token")), RunE: func(cmd *cobra.Command, _ []string) error {
		model, err := parseInput(params.Printer, cmd, nil)
		if err != nil {
			return err
		}
		apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
		if err != nil {
			return err
		}
		resp, err := buildRequest(context.Background(), model, apiClient).Execute()
		if err != nil {
			return fmt.Errorf("create instance token: %w", err)
		}
		return outputResult(params.Printer, model.OutputFormat, resp)
	}}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIDFlag, "ID of the instance")
	cmd.Flags().String(nameFlag, "", "Token name")
	cmd.Flags().String(descriptionFlag, "", "Token description")
	cmd.Flags().StringToString(labelFlag, nil, `Labels as key-value pairs, e.g. "--label env=prod"`)
	cmd.Flags().String(ttlDurationFlag, "", "Token time to live duration")
	_ = flags.MarkFlagsRequired(cmd, instanceIDFlag, nameFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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
	model := &inputModel{GlobalFlagModel: globalFlags, InstanceID: flags.FlagToStringValue(p, cmd, instanceIDFlag), Name: flags.FlagToStringValue(p, cmd, nameFlag), Description: flags.FlagToStringPointer(p, cmd, descriptionFlag), Labels: labelsPtr, TTLDuration: flags.FlagToStringPointer(p, cmd, ttlDurationFlag)}
	if model.Name == "" {
		return nil, fmt.Errorf("%s flag is required", nameFlag)
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiCreateInstanceTokenRequest {
	payload := modelexperiments.CreateInstanceTokenPayload{Name: model.Name, Description: model.Description, Labels: model.Labels, TtlDuration: model.TTLDuration}
	return apiClient.DefaultAPI.CreateInstanceToken(ctx, model.ProjectId, model.GlobalFlagModel.Region, model.InstanceID).CreateInstanceTokenPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, resp *modelexperiments.CreateInstanceTokenResponse) error {
	if resp == nil || resp.Token.Name == "" {
		return fmt.Errorf("response token is nil")
	}
	return p.OutputResult(outputFormat, resp.Token, func() error {
		p.Outputf(
			"Created instance token. ID: %s\nToken: %s\n",
			resp.Token.Id,
			resp.Token.Content,
		)
		return nil
	})
}
