package create

import (
	"context"
	"encoding/json"
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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
	"github.com/stackitcloud/stackit-sdk-go/services/observability/wait"
)

const (
	payloadFlag    = "payload"
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Payload    *observability.CreateScrapeConfigPayload
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a scrape configuration for an Observability instance",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Creates a scrape configuration job for an Observability instance.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"If no payload is provided, a default payload will be used.",
			"See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_create for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a scrape configuration on Observability instance "xxx" using default configuration`,
				"$ stackit observability scrape-config create"),
			examples.NewExample(
				`Create a scrape configuration on Observability instance "xxx" using an API payload sourced from the file "./payload.json"`,
				"$ stackit observability scrape-config create --payload @./payload.json --instance-id xxx"),
			examples.NewExample(
				`Create a scrape configuration on Observability instance "xxx" using an API payload provided as a JSON string`,
				`$ stackit observability scrape-config create --payload "{...}" --instance-id xxx`),
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit observability scrape-config generate-payload > ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit observability scrape-config create --payload @./payload.json --instance-id xxx`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			// Fill in default payload, if needed
			if model.Payload == nil {
				defaultPayload := observabilityUtils.DefaultCreateScrapeConfigPayload
				if err != nil {
					return fmt.Errorf("get default payload: %w", err)
				}
				model.Payload = &defaultPayload
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create scrape configuration %q on Observability instance %q?", *model.Payload.JobName, instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("create scrape configuration: %w", err)
			}

			jobName := model.Payload.JobName

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating scrape config")
				_, err = wait.CreateScrapeConfigWaitHandler(ctx, apiClient, model.InstanceId, *jobName, model.ProjectId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for scrape configuration creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			params.Printer.Outputf("%s scrape configuration with name %q for Observability instance %q\n", operationState, utils.PtrString(jobName), instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit observability scrape-config generate-payload")`)
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadValue := flags.FlagToStringPointer(p, cmd, payloadFlag)
	var payload *observability.CreateScrapeConfigPayload
	if payloadValue != nil {
		payload = &observability.CreateScrapeConfigPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("encode payload: %w", err)
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Payload:         payload,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiCreateScrapeConfigRequest {
	req := apiClient.CreateScrapeConfig(ctx, model.InstanceId, model.ProjectId)

	req = req.CreateScrapeConfigPayload(*model.Payload)
	return req
}
