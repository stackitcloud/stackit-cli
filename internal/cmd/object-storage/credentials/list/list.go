package list

import (
	"context"
	"encoding/json"
	"fmt"

	objectStorageUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"

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
	limitFlag            = "limit"
	credentialsGroupFlag = "credentials-group"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialsGroupId string
	Limit              *int64
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all credentials for an Object Storage credentials group",
		Long:  "Lists all credentials for a credentials group.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all credentials for a credentials group`,
				"$ stackit object-storage credentials list --credentials-group xxx"),
			examples.NewExample(
				`List all credentials for a credentials group in JSON format`,
				"$ stackit object-storage credentials list --credentials-group xxx --output-format json"),
			examples.NewExample(
				`List up to 10 credentials for a credentials group`,
				"$ stackit object-storage credentials list --credentials-group xxx --limit 10"),
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
				return fmt.Errorf("list Object Storage credentials: %w", err)
			}
			credentials := *resp.AccessKeys
			if len(credentials) == 0 {
				credentialsGroupLabel, err := objectStorageUtils.GetCredentialsGroupName(ctx, apiClient, model.ProjectId, model.CredentialsGroupId)
				if err != nil {
					credentialsGroupLabel = model.CredentialsGroupId
				}

				cmd.Printf("No credentials found for your credentials group %q\n", credentialsGroupLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(credentials) > int(*model.Limit) {
				credentials = credentials[:*model.Limit]
			}
			return outputResult(cmd, model.OutputFormat, credentials)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), credentialsGroupFlag, "Credentials Group ID")

	err := flags.MarkFlagsRequired(cmd, credentialsGroupFlag)
	cobra.CheckErr(err)
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
		GlobalFlagModel:    globalFlags,
		CredentialsGroupId: flags.FlagToStringValue(cmd, credentialsGroupFlag),
		Limit:              limit,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiListAccessKeysRequest {
	req := apiClient.ListAccessKeys(ctx, model.ProjectId)
	req = req.CredentialsGroup(model.CredentialsGroupId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, credentials []objectstorage.AccessKey) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		details, err := json.MarshalIndent(credentials, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Object Storage credentials list: %w", err)
		}
		cmd.Println(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("CREDENTIAL ID", "ACCESS KEY ID", "EXPIRES AT")
		for i := range credentials {
			c := credentials[i]

			expiresAt := "Never"
			if c.Expires != nil {
				expiresAt = *c.Expires
			}
			table.AddRow(*c.KeyId, *c.DisplayName, expiresAt)
		}
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
