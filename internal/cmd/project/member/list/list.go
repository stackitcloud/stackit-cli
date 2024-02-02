package list

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/membership/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/membership"
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List members of a project",
		Long:  "List members of a project",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all members of a project`,
				"$ stackit project role list --project-id xxx"),
			examples.NewExample(
				`List all members of a project, sorted by role`,
				"$ stackit project role list --project-id xxx --sort-by role"),
			examples.NewExample(
				`List up to 10 members of a project`,
				"$ stackit project role list --project-id xxx --limit 10"),
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
				return fmt.Errorf("list members: %w", err)
			}
			members := *resp.Members
			if len(members) == 0 {
				projectLabel, err := projectname.GetProjectName(ctx, cmd)
				if err != nil {
					projectLabel = model.ProjectId
				}
				cmd.Printf("No members found for project %s\n", projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(members) > int(*model.Limit) {
				members = members[:*model.Limit]
			}

			return outputResult(cmd, model, members)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	sortByFlagOptions := []string{"subject", "role"}

	cmd.Flags().String(subjectFlag, "", "Filter by subject (Identifier of user, service account or client. Usually email address in case of users or name in case of clients)")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.EnumFlag(false, "subject", sortByFlagOptions...), sortByFlag, fmt.Sprintf("Sort entries by a specific field, one of %q", sortByFlagOptions))
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
		Subject:         flags.FlagToStringPointer(cmd, subjectFlag),
		Limit:           flags.FlagToInt64Pointer(cmd, limitFlag),
		SortBy:          flags.FlagWithDefaultToStringValue(cmd, sortByFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *membership.APIClient) membership.ApiListMembersRequest {
	req := apiClient.ListMembers(ctx, projectResourceType, model.GlobalFlagModel.ProjectId)
	if model.Subject != nil {
		req = req.Subject(*model.Subject)
	}
	return req
}

func outputResult(cmd *cobra.Command, model *inputModel, members []membership.Member) error {
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
	case globalflags.JSONOutputFormat:
		// Show details
		details, err := json.MarshalIndent(members, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal members: %w", err)
		}
		cmd.Println(string(details))

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
			table.AddRow(*m.Subject, *m.Role)
		}

		if model.SortBy == "subject" {
			table.EnableAutoMergeOnColumns(1)
		} else if model.SortBy == "role" {
			table.EnableAutoMergeOnColumns(2)
		}

		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
