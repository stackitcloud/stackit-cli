package update

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

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

const (
	jobNameArg = "JOB_NAME"

	instanceIdFlag = "instance-id"
	payloadFlag    = "payload"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	JobName    string
	InstanceId string
	Payload    observability.UpdateScrapeConfigPayload
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", jobNameArg),
		Short: "Updates a scrape configuration of an Observability instance",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Updates a scrape configuration of an Observability instance.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_update for information regarding the payload structure.",
		),
		Args: args.SingleArg(jobNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update a scrape configuration with name "my-config" from Observability instance "xxx", using an API payload sourced from the file "./payload.json"`,
				"$ stackit observability scrape-config update my-config --payload @./payload.json --instance-id xxx"),
			examples.NewExample(
				`Update an scrape configuration with name "my-config" from Observability instance "xxx", using an API payload provided as a JSON string`,
				`$ stackit observability scrape-config update my-config --payload "{...}" --instance-id xxx`),
			examples.NewExample(
				`Generate a payload with the current values of a scrape configuration, and adapt it with custom values for the different configuration options`,
				`$ stackit observability scrape-config generate-payload --job-name my-config > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit observability scrape-configs update my-config --payload @./payload.json`),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update scrape configuration %q?", model.JobName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update scrape config: %w", err)
			}

			// The API has no status to wait on, so async mode is default
			params.Printer.Info("Updated Observability scrape configuration with name %q\n", model.JobName)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.ReadFromFileFlag(), payloadFlag, `Request payload (JSON). Can be a string or a file path, if prefixed with "@". Example: @./payload.json`)
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, payloadFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	clusterName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadString := flags.FlagToStringValue(p, cmd, payloadFlag)
	var payload observability.UpdateScrapeConfigPayload
	err := json.Unmarshal([]byte(payloadString), &payload)
	if err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		JobName:         clusterName,
		Payload:         payload,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiUpdateScrapeConfigRequest {
	req := apiClient.UpdateScrapeConfig(ctx, model.InstanceId, model.JobName, model.ProjectId)

	req = req.UpdateScrapeConfigPayload(model.Payload)
	return req
}
