package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/git/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/git"
	"github.com/stackitcloud/stackit-sdk-go/services/git/wait"
)

const (
	nameFlag = "name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Id   *string
	Name string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates STACKIT Git instance",
		Long:  "Create an STACKIT Git instance by name.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an instance with name 'my-new-instance'`,
				`$ stackit git create --name my-new-instance`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create the instance %q?", model.Name)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			result, err := request.Execute()
			if err != nil {
				return fmt.Errorf("create stackit git instance: %w", err)
			}
			model.Id = result.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Creating stackit git instance")
				_, err = wait.CreateGitInstanceWaitHandler(ctx, apiClient, model.ProjectId, *model.Id).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for stackit git Instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p, model, result)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "The name of the instance.")
	if err := flags.MarkFlagsRequired(cmd, nameFlag); err != nil {
		cobra.CheckErr(err)
	}
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	name := flags.FlagToStringValue(p, cmd, nameFlag)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            name,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *git.APIClient) git.ApiCreateInstanceRequest {
	return apiClient.CreateInstance(ctx, model.ProjectId).CreateInstancePayload(createPayload(model))
}

func createPayload(model *inputModel) git.CreateInstancePayload {
	return git.CreateInstancePayload{
		Name: &model.Name,
	}
}

func outputResult(p *print.Printer, model *inputModel, resp *git.Instance) error {
	if model == nil {
		return fmt.Errorf("input model is nil")
	}
	var outputFormat string
	if model.GlobalFlagModel != nil {
		outputFormat = model.OutputFormat
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal iminstanceage: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created instance %q with id %s\n", model.Name, utils.PtrString(model.Id))
		return nil
	}
}
