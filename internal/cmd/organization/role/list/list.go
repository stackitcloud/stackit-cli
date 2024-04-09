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
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/authorization/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
)

const (
	organizationIdFlag = "organization-id"
	limitFlag          = "limit"

	organizationResourceType = "organization"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	OrganizationId *string
	Limit          *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists roles and permissions of an organization",
		Long:  "Lists roles and permissions of an organization.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all roles and permissions of an organization`,
				"$ stackit organization role list --organization-id xxx"),
			examples.NewExample(
				`List all roles and permissions of an organization in JSON format`,
				"$ stackit organization role list --organization-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 roles and permissions of an organization`,
				"$ stackit organization role list --organization-id xxx --limit 10"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get organization roles: %w", err)
			}
			roles := *resp.Roles
			if len(roles) == 0 {
				p.Info("No roles found for organization with ID %q\n", *model.OrganizationId)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(roles) > int(*model.Limit) {
				roles = roles[:*model.Limit]
			}

			return outputRolesResult(p, model.OutputFormat, roles)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(organizationIdFlag, "", "Organization ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)

	limit := flags.FlagToInt64Pointer(cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(cmd, organizationIdFlag),
		Limit:           flags.FlagToInt64Pointer(cmd, limitFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiListRolesRequest {
	return apiClient.ListRoles(ctx, organizationResourceType, *model.OrganizationId)
}

func outputRolesResult(p *print.Printer, outputFormat string, roles []authorization.Role) error {
	switch outputFormat {
	case globalflags.JSONOutputFormat:
		// Show details
		details, err := json.MarshalIndent(roles, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal roles: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ROLE NAME", "ROLE DESCRIPTION", "PERMISSION NAME", "PERMISSION DESCRIPTION")
		for i := range roles {
			r := roles[i]
			for j := range *r.Permissions {
				p := (*r.Permissions)[j]
				table.AddRow(*r.Name, *r.Description, *p.Name, *p.Description)
			}
			table.AddSeparator()
		}
		table.EnableAutoMergeOnColumns(1, 2)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
