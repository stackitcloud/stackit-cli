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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
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
			model, err := parseInput(params.Printer, cmd, args)
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
				return fmt.Errorf("list Object Storage credentials groups: %w", err)
			}
			credentialsGroups := resp.GetCredentialsGroups()

			// Truncate output
			if model.Limit != nil && len(credentialsGroups) > int(*model.Limit) {
				credentialsGroups = credentialsGroups[:*model.Limit]
			}
			return outputResult(params.Printer, model.OutputFormat, credentialsGroups)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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
		GlobalFlagModel: globalFlags,
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiListCredentialsGroupsRequest {
	req := apiClient.ListCredentialsGroups(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, credentialsGroups []objectstorage.CredentialsGroup) error {
	return p.OutputResult(outputFormat, credentialsGroups, func() error {
		if len(credentialsGroups) == 0 {
			p.Outputf("No credentials groups found for your project")
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "URN")
		for i := range credentialsGroups {
			c := credentialsGroups[i]
			table.AddRow(
				utils.PtrString(c.CredentialsGroupId),
				utils.PtrString(c.DisplayName),
				utils.PtrString(c.Urn),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
