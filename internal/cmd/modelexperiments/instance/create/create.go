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
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	modelexperiments "github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/modelexperiments/v1api/wait"
)

const (
	nameFlag        = "name"
	descriptionFlag = "description"
	labelsFlag      = "labels"
	retentionFlag   = "deleted-experiment-retention"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Name        string
	Description *string
	Labels      *map[string]string
	Retention   *string
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
				`$ stackit model-experiments instance create --name my-tracking`),
			examples.NewExample(
				`Create an instance with a description and labels`,
				`$ stackit model-experiments instance create --name my-tracking --description "team tracking server" --labels env=prod,team=ml`),
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
			}

			prompt := fmt.Sprintf("Are you sure you want to create an AI Model Experiments instance for project %q?", projectLabel)
			if err := params.Printer.PromptForConfirmation(prompt); err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create AI Model Experiments instance: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("create AI Model Experiments instance: response instance is nil")
			}

			if !model.Async {
				err := spinner.Run(params.Printer, "Creating AI Model Experiments instance", func() error {
					_, err = wait.CreateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.Region, model.ProjectId, resp.Instance.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for AI Model Experiments instance creation: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Instance name")
	cmd.Flags().String(descriptionFlag, "", "Instance description")
	cmd.Flags().StringToString(labelsFlag, nil, `Labels as comma-separated key-value pairs, e.g. "--labels env=prod,team=ml"`)
	cmd.Flags().String(retentionFlag, "", `Retention period for deleted experiments before permanent purge, e.g. "30d" (min 1d, max 90d)`)
	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		Retention:       flags.FlagToStringPointer(p, cmd, retentionFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *modelexperiments.APIClient) modelexperiments.ApiCreateInstanceRequest {
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

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, resp *modelexperiments.CreateInstanceResponse) error {
	var instance any
	if resp != nil {
		instance = resp.Instance
	}

	return p.OutputResult(outputFormat, instance, func() error {
		if resp == nil {
			return fmt.Errorf("response instance is nil")
		}

		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s AI Model Experiments instance for project %q. Instance ID: %s\n", operationState, projectLabel, resp.Instance.Id)
		return nil
	})
}
