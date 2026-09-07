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
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
)

const (
	nameFlag        = "name"
	descriptionFlag = "description"
	labelFlag       = "label"
	retentionFlag   = "deleted-experiment-retention"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Name        string
	Description *string
	Labels      *map[string]string
	Retention   *string
	Region      string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an AI Model Experiments instance",
		Long:  "Creates an AI Model Experiments (MLflow) instance in your STACKIT project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an AI Model Experiments instance with name "my-tracking"`,
				`$ stackit ai-model-experiments instance create --name my-tracking`),
			examples.NewExample(
				`Create an instance with a description and labels`,
				`$ stackit ai-model-experiments instance create --name my-tracking --description "team tracking server" --label env=prod`),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create an AI Model Experiments instance for project %q?", projectLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildCreateInstanceRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create AI Model Experiments instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Instance name")
	cmd.Flags().String(descriptionFlag, "", "Instance description")
	cmd.Flags().StringToString(labelFlag, nil, `Labels as key-value pairs, e.g. "--label env=prod"`)
	cmd.Flags().String(retentionFlag, "", `Retention period for deleted experiments before permanent purge, e.g. "30d" (min 1d, max 90d)`)
	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          labelsPtr,
		Retention:       flags.FlagToStringPointer(p, cmd, retentionFlag),
		Region:          flags.FlagToStringValue(p, cmd, globalflags.RegionFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildCreateInstanceRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiCreateInstanceRequest {
	req := apiClient.DefaultAPI.CreateInstance(ctx, model.ProjectId, model.Region)

	payload := modelexperiments.CreateInstancePayload{
		Name: model.Name,
	}
	if model.Description != nil {
		payload.Description = model.Description
	}
	if model.Labels != nil {
		payload.Labels = model.Labels
	}
	if model.Retention != nil {
		payload.DeletedExperimentRetention = model.Retention
	}

	return req.CreateInstancePayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *modelexperiments.CreateInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("response instance is nil")
	}

	return p.OutputResult(outputFormat, resp.Instance, func() error {
		p.Outputf("Creating AI Model Experiments instance for project %q. Instance ID: %s\n", projectLabel, resp.Instance.Id)
		return nil
	})
}
