package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	argusUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
	"github.com/stackitcloud/stackit-sdk-go/services/argus/wait"
)

const (
	payloadFlag    = "payload"
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Payload    *argus.CreateScrapeConfigPayload
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a scrape configuration for an Argus instance",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s",
			"Creates a scrape configuration job for an Argus instance.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"If no payload is provided, a default payload will be used.",
			"See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_create for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a scrape configuration on Argus instance "xxx" using default configuration`,
				"$ stackit argus scrape-config create"),
			examples.NewExample(
				`Create a scrape configuration on Argus instance "xxx" using an API payload sourced from the file "./payload.json"`,
				"$ stackit argus scrape-config create --payload @./payload.json --instance-id xxx"),
			examples.NewExample(
				`Create a scrape configuration on Argus instance "xxx" using an API payload provided as a JSON string`,
				`$ stackit argus scrape-config create --payload "{...}" --instance-id xxx`),
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit argus scrape-config generate-payload > ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit argus scrape-config create --payload @./payload.json --instance-id xxx`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			instanceLabel, err := argusUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a scrape configuration on Argus instance %q?", instanceLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Fill in default payload, if needed
			if model.Payload == nil {
				defaultPayload := argusUtils.DefaultCreateScrapeConfigPayload
				if err != nil {
					return fmt.Errorf("get default payload: %w", err)
				}
				model.Payload = &defaultPayload
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
				s := spinner.New(p)
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
			p.Outputf("%s scrape configuration with name %q for Argus instance %q\n", operationState, *jobName, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@" (example: @./payload.json). If unset, will use a default payload (you can check it by running "stackit argus scrape-config generate-payload")`)
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadValue := flags.FlagToStringPointer(cmd, payloadFlag)
	var payload *argus.CreateScrapeConfigPayload
	if payloadValue != nil {
		payload = &argus.CreateScrapeConfigPayload{}
		err := json.Unmarshal([]byte(*payloadValue), payload)
		if err != nil {
			return nil, fmt.Errorf("encode payload: %w", err)
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Payload:         payload,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiCreateScrapeConfigRequest {
	req := apiClient.CreateScrapeConfig(ctx, model.InstanceId, model.ProjectId)

	req = req.CreateScrapeConfigPayload(*model.Payload)
	return req
}
