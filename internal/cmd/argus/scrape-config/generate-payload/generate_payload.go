package generatepayload

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
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
	filePathFlag   = "file-path"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	JobName    *string
	InstanceId string
	FilePath   *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-payload",
		Short: "Generates a payload to create/update scrape configurations for an Argus instance ",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n",
			"Generates a JSON payload with values to be used as --payload input for scrape configurations creation or update.",
			"This command can be used to generate a payload to update an existing scrape config or to create a new scrape config job.",
			"To update an existing scrape config job, provide the job name and the instance ID of the Argus instance.",
			"To obtain a default payload to create a new scrape config job, run the command with no flags.",
			"Note that some of the default values provided, such as the job name, the metrics path and URL of the targets, should be adapted to your use case.",
			"See https://docs.api.stackit.cloud/documentation/argus/version/v1#tag/scrape-config/operation/v1_projects_instances_scrapeconfigs_create for information regarding the payload structure.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a Create payload with default values, and adapt it with custom values for the different configuration options`,
				`$ stackit argus scrape-config generate-payload --file-path ./payload.json`,
				`<Modify payload in file, if needed>`,
				`$ stackit argus scrape-config create my-config --payload @./payload.json`),
			examples.NewExample(
				`Generate an Update payload with the values of an existing configuration named "my-config" for Argus instance xxx, and adapt it with custom values for the different configuration options`,
				`$ stackit argus scrape-config generate-payload --job-name my-config --instance-id xxx --file-path ./payload.json`,
				`<Modify payload in file>`,
				`$ stackit argus scrape-config update my-config --payload @./payload.json`),
			examples.NewExample(
				`Generate an Update payload with the values of an existing configuration named "my-config" for Argus instance xxx, and preview it in the terminal`,
				`$ stackit argus scrape-config generate-payload --job-name my-config --instance-id xxx`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			if model.JobName == nil {
				createPayload := argusUtils.DefaultCreateScrapeConfigPayload
				return outputCreateResult(p, model.FilePath, &createPayload)
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read Argus scrape config: %w", err)
			}

			payload, err := argusUtils.MapToUpdateScrapeConfigPayload(resp)
			if err != nil {
				return fmt.Errorf("map update scrape config payloads: %w", err)
			}

			return outputUpdateResult(p, model.FilePath, payload)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().StringP(jobNameFlag, "n", "", "If set, generates an update payload with the current state of the given scrape config. If unset, generates a create payload with default values")
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the payload to the given file. If unset, writes the payload to the standard output")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	jobName := flags.FlagToStringPointer(p, cmd, jobNameFlag)
	instanceId := flags.FlagToStringValue(p, cmd, instanceIdFlag)

	if jobName != nil && (globalFlags.ProjectId == "" || instanceId == "") {
		return nil, fmt.Errorf("if a job-name is provided then instance-id and project-id must be provided")
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		JobName:         jobName,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		FilePath:        flags.FlagToStringPointer(p, cmd, filePathFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *argus.APIClient) argus.ApiGetScrapeConfigRequest {
	req := apiClient.GetScrapeConfig(ctx, model.InstanceId, *model.JobName, model.ProjectId)
	return req
}

func outputCreateResult(p *print.Printer, filePath *string, payload *argus.CreateScrapeConfigPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if filePath != nil {
		err = fileutils.FileOutput(*filePath, string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}

func outputUpdateResult(p *print.Printer, filePath *string, payload *argus.UpdateScrapeConfigPayload) error {
	payloadBytes, err := json.MarshalIndent(*payload, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	if filePath != nil {
		err = fileutils.FileOutput(*filePath, string(payloadBytes))
		if err != nil {
			return fmt.Errorf("write payload to the file: %w", err)
		}
	} else {
		p.Outputln(string(payloadBytes))
	}

	return nil
}
