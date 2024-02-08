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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/authorization/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
)

const (
	organizationIdFlag = "organization-id"
	subjectFlag        = "subject"
	limitFlag          = "limit"
	sortByFlag         = "sort-by"

	organizationResourceType = "organization"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	OrganizationId *string
	Subject        *string
	Limit          *int64
	SortBy         string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists members of an organization",
		Long:  "Lists members of an organization",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all members of an organization`,
				"$ stackit organization member list --organization-id xxx"),
			examples.NewExample(
				`List all members of an organization in JSON format`,
				"$ stackit organization member list --organization-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 members of an organization`,
				"$ stackit organization member list --organization-id xxx --limit 10"),
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
				cmd.Printf("No members found for organization with ID %s\n", *model.OrganizationId)
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

	cmd.Flags().String(organizationIdFlag, "", "The organization ID")
	cmd.Flags().String(subjectFlag, "", "Filter by subject (Identifier of user, service account or client. Usually email address in case of users or name in case of clients)")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Var(flags.EnumFlag(false, "subject", sortByFlagOptions...), sortByFlag, fmt.Sprintf("Sort entries by a specific field, one of %q", sortByFlagOptions))

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
		Subject:         flags.FlagToStringPointer(cmd, subjectFlag),
		Limit:           flags.FlagToInt64Pointer(cmd, limitFlag),
		SortBy:          flags.FlagWithDefaultToStringValue(cmd, sortByFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiListMembersRequest {
	req := apiClient.ListMembers(ctx, organizationResourceType, *model.OrganizationId)
	if model.Subject != nil {
		req = req.Subject(*model.Subject)
	}
	return req
}

func outputResult(cmd *cobra.Command, model *inputModel, members []authorization.Member) error {
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
