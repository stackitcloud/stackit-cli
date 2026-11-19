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
	instanceIdArg = "INSTANCE_ID"

	nameFlag        = "name"
	descriptionFlag = "description"
	labelsFlag      = "labels"
	retentionFlag   = "deleted-experiment-retention"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string

	Name        *string
	Description *string
	Labels      *map[string]*string
	Retention   *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates an AI Model Experiments instance",
		Long:  "Partially updates an AI Model Experiments instance in a STACKIT project.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the name of an AI Model Experiments instance with ID "xxx"`,
				`$ stackit model-experiments instance update xxx --name my-new-name`,
			),
			examples.NewExample(
				`Update the description of an AI Model Experiments instance`,
				`$ stackit model-experiments instance update xxx --description "team tracking server"`,
			),
			examples.NewExample(
				`Update labels on an AI Model Experiments instance`,
				`$ stackit model-experiments instance update xxx --labels env=prod,team=ml`,
			),
			examples.NewExample(
				`Update multiple fields of an AI Model Experiments instance`,
				`$ stackit model-experiments instance update xxx --name my-new-name --description "team tracking server" --labels env=prod,team=ml --deleted-experiment-retention 30d`,
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			instanceLabel, err := meUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to update AI Model Experiments instance %q in project %q?", instanceLabel, projectLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update AI Model Experiments instance: %w", err)
			}

			return outputResult(params.Printer, model, resp)
		},
	}

	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Instance name")
	cmd.Flags().String(descriptionFlag, "", "Instance description")
	cmd.Flags().StringToString(labelsFlag, nil, `Labels are comma-separated key-value string pairs, e.g. "--labels env=prod,team=ml"`)
	cmd.Flags().String(retentionFlag, "", `Retention period for deleted experiments before permanent purge, e.g. "30d" (min 1d, max 90d)`)
	cmd.MarkFlagsOneRequired(nameFlag, descriptionFlag, labelsFlag, retentionFlag)
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

	model := inputModel{
		GlobalFlagModel: globalFlags,

		InstanceId: inputArgs[0],

		Name:        flags.FlagToStringPointer(p, cmd, nameFlag),
		Description: flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:      labelsPtr,
		Retention:   flags.FlagToStringPointer(p, cmd, retentionFlag),
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiPartialUpdateInstanceRequest {
	req := apiClient.DefaultAPI.PartialUpdateInstance(ctx, model.ProjectId, model.Region, model.InstanceId)

	payload := modelexperiments.PartialUpdateInstancePayload{}

	if model.Name != nil {
		payload.Name = model.Name
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

	return req.PartialUpdateInstancePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, resp *modelexperiments.PartialUpdateInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("response instance is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Updated AI Model Experiments instance. Instance ID: %s\n", resp.Instance.Id)
		return nil
	})
}
