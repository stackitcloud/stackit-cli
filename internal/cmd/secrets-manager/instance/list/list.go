package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"

	"github.com/spf13/cobra"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Secrets Manager instances",
		Long:  "Lists all Secrets Manager instances.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Secrets Manager instances`,
				"$ stackit secrets-manager instance list"),
			examples.NewExample(
				`List all Secrets Manager instances in JSON format`,
				"$ stackit secrets-manager instance list --output-format json"),
			examples.NewExample(
				`List up to 10 Secrets Manager instances`,
				"$ stackit secrets-manager instance list --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Secret Manager instances: %w", err)
			}

			if resp.Instances == nil || len(*resp.Instances) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, cmd)
				if err != nil {
					projectLabel = model.ProjectId
				}
				cmd.Printf("No instances found for project %q\n", projectLabel)
				return nil
			}
			instances := *resp.Instances

			// Truncate output
			if model.Limit != nil && len(instances) > int(*model.Limit) {
				instances = instances[:*model.Limit]
			}

			return outputResult(cmd, model.OutputFormat, instances)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiListInstancesRequest {
	req := apiClient.ListInstances(ctx, model.ProjectId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, instances []secretsmanager.Instance) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(instances, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Secret Manager instance list: %w", err)
		}
		cmd.Println(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATE", "SECRETS")
		// TODO: Add support for Number of Users (to match Portal)
		for i := range instances {
			instance := instances[i]
			table.AddRow(*instance.Id, *instance.Name, *instance.State, *instance.SecretCount)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
