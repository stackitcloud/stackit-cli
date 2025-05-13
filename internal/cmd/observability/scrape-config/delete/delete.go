package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
	"github.com/stackitcloud/stackit-sdk-go/services/observability/wait"
)

const (
	jobNameArg = "JOB_NAME"

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	JobName    string
	InstanceId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", jobNameArg),
		Short: "Deletes a scrape configuration from an Observability instance",
		Long:  "Deletes a scrape configuration from an Observability instance.",
		Args:  args.SingleArg(jobNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete a scrape configuration job with name "my-config" from Observability instance "xxx"`,
				"$ stackit observability scrape-config delete my-config --instance-id xxx"),
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

			instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete scrape configuration %q on Observability instance %q? (This cannot be undone)", model.JobName, instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete scrape configuration: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Deleting scrape config")
				_, err = wait.DeleteScrapeConfigWaitHandler(ctx, apiClient, model.InstanceId, model.JobName, model.ProjectId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for scrape config deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			params.Printer.Info("%s scrape configuration with name %q for Observability instance %q\n", operationState, model.JobName, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	jobName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		JobName:         jobName,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiDeleteScrapeConfigRequest {
	req := apiClient.DeleteScrapeConfig(ctx, model.InstanceId, model.JobName, model.ProjectId)
	return req
}
