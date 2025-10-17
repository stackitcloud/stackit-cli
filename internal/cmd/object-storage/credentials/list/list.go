package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	objectStorageUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	limitFlag              = "limit"
	credentialsGroupIdFlag = "credentials-group-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialsGroupId string
	Limit              *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all credentials for an Object Storage credentials group",
		Long:  "Lists all credentials for an Object Storage credentials group.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all credentials for a credentials group with ID "xxx"`,
				"$ stackit object-storage credentials list --credentials-group-id xxx"),
			examples.NewExample(
				`List all credentials for a credentials group with ID "xxx" in JSON format`,
				"$ stackit object-storage credentials list --credentials-group-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 credentials for a credentials group with ID "xxx"`,
				"$ stackit object-storage credentials list --credentials-group-id xxx --limit 10"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
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
				credentialsGroupLabel, err := objectStorageUtils.GetCredentialsGroupName(ctx, apiClient, model.ProjectId, model.CredentialsGroupId, model.Region)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get credentials group name: %v", err)
					credentialsGroupLabel = model.CredentialsGroupId
				}

				params.Printer.Info("No credentials found for credentials group %q\n", credentialsGroupLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(credentials) > int(*model.Limit) {
				credentials = credentials[:*model.Limit]
			}
			return outputResult(params.Printer, model.OutputFormat, credentials)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.UUIDFlag(), credentialsGroupIdFlag, "Credentials Group ID")

	err := flags.MarkFlagsRequired(cmd, credentialsGroupIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		CredentialsGroupId: flags.FlagToStringValue(p, cmd, credentialsGroupIdFlag),
		Limit:              limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiListAccessKeysRequest {
	req := apiClient.ListAccessKeys(ctx, model.ProjectId, model.Region)
	req = req.CredentialsGroup(model.CredentialsGroupId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, credentials []objectstorage.AccessKey) error {
	return p.OutputResult(outputFormat, credentials, func() error {
		table := tables.NewTable()
		table.SetHeader("CREDENTIALS ID", "ACCESS KEY ID", "EXPIRES AT")
		for i := range credentials {
			c := credentials[i]

			expiresAt := utils.PtrStringDefault(c.Expires, "Never")
			table.AddRow(utils.PtrString(c.KeyId), utils.PtrString(c.DisplayName), expiresAt)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
