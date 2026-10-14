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
	instanceIdArg = "INSTANCE_ID"

	nameFlag        = "name"
	descriptionFlag = "description"
	labelFlag       = "label"
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
		Use:   fmt.Sprintf("patch %s", instanceIdArg),
		Short: "Updates an AI Model Experiments instance",
		Long:  "Partially updates an AI Model Experiments instance in a STACKIT project.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the name of an AI Model Experiments instance with ID "xxx"`,
				`$ stackit ai-model-experiments instance patch xxx --name my-new-name`,
			),
			examples.NewExample(
				`Update the description of an AI Model Experiments instance`,
				`$ stackit ai-model-experiments instance patch xxx --description "team tracking server"`,
			),
			examples.NewExample(
				`Update labels on an AI Model Experiments instance`,
				`$ stackit ai-model-experiments instance patch xxx --label env=prod`,
			),
			examples.NewExample(
				`Update multiple fields of an AI Model Experiments instance`,
				`$ stackit ai-model-experiments instance patch xxx --name my-new-name --description "team tracking server" --label env=prod --deleted-experiment-retention 30d`,
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

			req := buildPatchInstanceRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update AI Model Experiments instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
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

	labels, err := cmd.Flags().GetStringToString(labelFlag)
	if err != nil {
		return nil, fmt.Errorf("parse %q flag: %w", labelFlag, err)
	}

	var labelsPtr *map[string]*string

	if len(labels) > 0 {
		convertedLabels := make(map[string]*string, len(labels))

		for key, value := range labels {
			valueCopy := value
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

	if model.Name == nil &&
		model.Description == nil &&
		model.Labels == nil &&
		model.Retention == nil {
		return nil, fmt.Errorf("at least one update flag must be provided")
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildPatchInstanceRequest(
	ctx context.Context,
	model *inputModel,
	apiClient *modelexperiments.APIClient,
) modelexperiments.ApiPartialUpdateInstanceRequest {
	req := apiClient.DefaultAPI.PartialUpdateInstance(
		ctx,
		model.ProjectId,
		model.Region,
		model.InstanceId,
	)

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

func outputResult(
	p *print.Printer,
	outputFormat string,
	resp *modelexperiments.PartialUpdateInstanceResponse,
) error {
	if resp == nil {
		return fmt.Errorf("response instance is nil")
	}

	return p.OutputResult(outputFormat, resp.Instance, func() error {
		p.Outputf(
			"Updated AI Model Experiments instance. Instance ID: %s\n",
			resp.Instance.Id,
		)
		return nil
	})
}
