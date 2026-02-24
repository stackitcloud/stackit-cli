package list

import (
	"context"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit  *int64
	Member string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all organizations",
		Long:  "Lists all organizations.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists organizations for your user`,
				"$ stackit organization list",
			),
			examples.NewExample(
				`Lists organizations for the user with the email foo@bar`,
				"$ stackit organization list --member foo@bar",
			),
			examples.NewExample(
				`Lists the first 10 organizations`,
				"$ stackit organization list --limit 10",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			model.Member = auth.GetProfileEmail(config.DefaultProfileName)

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return err
			}

			if resp == nil {
				params.Printer.Outputln("list organizations: empty response")
				return nil
			}

			return outputResult(params.Printer, model.OutputFormat, utils.PtrValue(resp.Items))
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list (default 50)")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

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

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiListOrganizationsRequest {
	req := apiClient.ListOrganizations(ctx)
	req = req.Member(model.Member)
	if model.Limit != nil {
		req = req.Limit(float32(*model.Limit))
	}
	return req
}

func outputResult(p *print.Printer, outputFormat string, organizations []resourcemanager.ListOrganizationsResponseItemsInner) error {
	return p.OutputResult(outputFormat, organizations, func() error {
		if len(organizations) == 0 {
			p.Outputln("No organizations found")
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "CONTAINER ID")

		for _, organization := range organizations {
			table.AddRow(
				utils.PtrString(organization.OrganizationId),
				utils.PtrString(organization.Name),
				utils.PtrString(organization.ContainerId),
			)
			table.AddSeparator()
		}

		p.Outputln(table.Render())
		return nil
	})
}
