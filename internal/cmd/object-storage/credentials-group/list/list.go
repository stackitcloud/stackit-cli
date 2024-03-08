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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
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
		Short: "Lists all credentials groups that hold Object Storage access credentials",
		Long:  "Lists all credentials groups that hold Object Storage access credentials.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all credentials groups`,
				"$ stackit object-storage credentials-group list"),
			examples.NewExample(
				`List all credentials groups in JSON format`,
				"$ stackit object-storage credentials-group list --output-format json"),
			examples.NewExample(
				`List up to 10 credentials groups`,
				"$ stackit object-storage credentials-group list --limit 10"),
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
				return fmt.Errorf("list Object Storage credentials groups: %w", err)
			}
			credentialsGroups := *resp.CredentialsGroups
			if len(credentialsGroups) == 0 {
				cmd.Println("No credentials groups found for your project")
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(credentialsGroups) > int(*model.Limit) {
				credentialsGroups = credentialsGroups[:*model.Limit]
			}
			return outputResult(cmd, model.OutputFormat, credentialsGroups)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiListCredentialsGroupsRequest {
	req := apiClient.ListCredentialsGroups(ctx, model.ProjectId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, credentialsGroups []objectstorage.CredentialsGroup) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(credentialsGroups, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Object Storage credentials group list: %w", err)
		}
		cmd.Println(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "URN")
		for i := range credentialsGroups {
			c := credentialsGroups[i]
			table.AddRow(*c.CredentialsGroupId, *c.DisplayName, *c.Urn)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
