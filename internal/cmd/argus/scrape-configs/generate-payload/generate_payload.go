package generatepayload

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

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
)

const (
	jobNameFlag    = "job-name"
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	JobName    *string
	InstanceId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-payload",
		Short: "Generates a payload to create/update Scrape Configurations for an Argus instance ",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Generates a JSON payload with values to be used as --payload input for Scrape Configurations creation or update.",
			"If --job-name is set, an Update payload will be generated with the current state of the given configuration. If unset, a Create payload will be generated with default values.",
			"See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_create for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit argus scrape-configs generate-payload > ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit argus scrape-configs create my-config --payload @./payload.json`),
			examples.NewExample(
				`Generate a payload with values of an existing configuration, and adapt it with custom values for the different configuration options`,
				`$ stackit argus scrape-configs generate-payload --job-name my-config > ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit argus scrape-configs update my-config --payload @./payload.json`),
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

			if model.JobName == nil {
				createPayload, err := argusUtils.GetDefaultCreateScrapeConfigPayload(ctx, apiClient)
				if err != nil {
					return err
				}
				return outputCreateResult(p, createPayload)
			} else {
				req := buildRequest(ctx, model, apiClient)
				resp, err := req.Execute()
				if err != nil {
					return fmt.Errorf("read SKE cluster: %w", err)
				}

				payload := argusUtils.MapToUpdateScrapeConfigPayload(resp)

				return outputUpdateResult(p, payload)
			}

		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().StringP(jobNameFlag, "n", "", "If set, generates the payload with the current state of the given cluster. If unset, generates the payload with default values")
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)

	jobName := flags.FlagToStringPointer(cmd, jobNameFlag)
	// If jobName is provided, projectId and instanceId are needed as well
	if jobName != nil {
		err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
		cobra.CheckErr(err)

		if globalFlags.ProjectId == "" {
			return nil, &errors.ProjectIdError{}
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		JobName:         jobName,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiGetScrapeConfigRequest {
	req := apiClient.GetScrapeConfig(ctx, model.InstanceId, *model.JobName, model.ProjectId)
	return req
}

func outputCreateResult(p *print.Printer, payload *argus.CreateScrapeConfigPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	p.Outputln(string(payloadBytes))

	return nil
}

func outputUpdateResult(p *print.Printer, payload *argus.UpdateScrapeConfigPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	p.Outputln(string(payloadBytes))

	return nil
}
