package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/authorization/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
)

const (
	limitFlag = "limit"

	projectResourceType = "project"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Limit *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists roles and permissions of a project",
		Long:  "Lists roles and permissions of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all roles and permissions of a project`,
				"$ stackit project role list --project-id xxx"),
			examples.NewExample(
				`List all roles and permissions of a project in JSON format`,
				"$ stackit project role list --project-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 roles and permissions of a project`,
				"$ stackit project role list --project-id xxx --limit 10"),
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
				return fmt.Errorf("get project roles: %w", err)
			}
			roles := *resp.Roles
			if len(roles) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No roles found for project %q\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(roles) > int(*model.Limit) {
				roles = roles[:*model.Limit]
			}

			return outputRolesResult(params.Printer, model.OutputFormat, roles)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
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
		GlobalFlagModel: globalFlags,
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiListRolesRequest {
	return apiClient.ListRoles(ctx, projectResourceType, model.GlobalFlagModel.ProjectId)
}

func outputRolesResult(p *print.Printer, outputFormat string, roles []authorization.Role) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		// Show details
		details, err := json.MarshalIndent(roles, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal roles: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(roles, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
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
				table.AddRow(
					utils.PtrString(r.Name),
					utils.PtrString(r.Description),
					utils.PtrString(p.Name),
					utils.PtrString(p.Description),
				)
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
