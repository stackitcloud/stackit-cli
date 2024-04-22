package update

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

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
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
	Payload    argus.UpdateScrapeConfigPayload
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", jobNameArg),
		Short: "Updates a scrape configuration of an Argus instance",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Updates a scrape configuration of an Argus instance.",
			"The payload can be provided as a JSON string or a file path prefixed with \"@\".",
			"See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_partial_update for information regarding the payload structure.",
		),
		Args: args.SingleArg(jobNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update a scrape configuration from Argus instance "xxx", using an API payload sourced from the file "./payload.json"`,
				"$ stackit argus scrape-config update my-config --payload @./payload.json --instance-id xxx"),
			examples.NewExample(
				`Update an scrape configuration from Argus instance "xxx", using an API payload provided as a JSON string`,
				`$ stackit argus scrape-config update my-config --payload "{...}" --instance-id xxx`),
			examples.NewExample(
				`Generate a payload with the current values of a scrape configuration, and adapt it with custom values for the different configuration options`,
				`$ stackit argus scrape-config generate-payload --job-name my-config > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit argus scrape-configs update my-config --payload @./payload.json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update scrape configuration %q?", model.JobName)
				err = p.PromptForConfirmation(prompt)
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
			p.Info("Updated Argus scrape configuration with name %q\n", model.JobName)
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

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	clusterName := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	payloadString := flags.FlagToStringValue(cmd, payloadFlag)
	var payload argus.UpdateScrapeConfigPayload
	err := json.Unmarshal([]byte(payloadString), &payload)
	if err != nil {
		return nil, fmt.Errorf("encode payload: %w", err)
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		JobName:         clusterName,
		Payload:         payload,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiUpdateScrapeConfigRequest {
	req := apiClient.UpdateScrapeConfig(ctx, model.InstanceId, model.JobName, model.ProjectId)

	req = req.UpdateScrapeConfigPayload(model.Payload)
	return req
}
