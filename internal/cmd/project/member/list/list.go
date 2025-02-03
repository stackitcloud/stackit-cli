package list

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
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
	subjectFlag = "subject"
	limitFlag   = "limit"
	sortByFlag  = "sort-by"

	projectResourceType = "project"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Subject *string
	Limit   *int64
	SortBy  string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists members of a project",
		Long:  "Lists members of a project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all members of a project`,
				"$ stackit project member list --project-id xxx"),
			examples.NewExample(
				`List all members of a project, sorted by role`,
				"$ stackit project member list --project-id xxx --sort-by role"),
			examples.NewExample(
				`List up to 10 members of a project`,
				"$ stackit project member list --project-id xxx --limit 10"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list members: %w", err)
			}
			members := *resp.Members
			if len(members) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
				if err != nil {
					p.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				p.Info("No members found for project %q\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(members) > int(*model.Limit) {
				members = members[:*model.Limit]
			}

			return outputResult(p, model, members)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	sortByFlagOptions := []string{"subject", "role"}

	cmd.Flags().String(subjectFlag, "", "Filter by subject (the identifier of a user, service account or client). This is usually the email address (for users) or name (for clients)")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.EnumFlag(false, "subject", sortByFlagOptions...), sortByFlag, fmt.Sprintf("Sort entries by a specific field, one of %q", sortByFlagOptions))
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
		Subject:         flags.FlagToStringPointer(p, cmd, subjectFlag),
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
		SortBy:          flags.FlagWithDefaultToStringValue(p, cmd, sortByFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiListMembersRequest {
	req := apiClient.ListMembers(ctx, projectResourceType, model.GlobalFlagModel.ProjectId)
	if model.Subject != nil {
		req = req.Subject(*model.Subject)
	}
	return req
}

func outputResult(p *print.Printer, model *inputModel, members []authorization.Member) error {
	sortFn := func(i, j int) bool {
		switch model.SortBy {
		case "subject":
			return *members[i].Subject < *members[j].Subject
		case "role":
			return *members[i].Role < *members[j].Role
		default:
			return false
		}
	}
	sort.SliceStable(members, sortFn)

	switch model.OutputFormat {
	case print.JSONOutputFormat:
		// Show details
		details, err := json.MarshalIndent(members, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal members: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(members, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal members: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("SUBJECT", "ROLE")
		for i := range members {
			m := members[i]
			// If the previous item differs from the current item on the element to sort by, add a separator between the rows to help readability
			if i > 0 && sortFn(i-1, i) {
				table.AddSeparator()
			}
			table.AddRow(utils.PtrString(m.Subject), utils.PtrString(m.Role))
		}

		if model.SortBy == "subject" {
			table.EnableAutoMergeOnColumns(1)
		} else if model.SortBy == "role" {
			table.EnableAutoMergeOnColumns(2)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
