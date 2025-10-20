package list

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/authorization/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
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
				return fmt.Errorf("list members: %w", err)
			}
			members := *resp.Members
			if len(members) == 0 {
				params.Printer.Info("No members found for organization with ID %q\n", *model.OrganizationId)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(members) > int(*model.Limit) {
				members = members[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, model.SortBy, members)
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
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		Subject:         flags.FlagToStringPointer(p, cmd, subjectFlag),
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
		SortBy:          flags.FlagWithDefaultToStringValue(p, cmd, sortByFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiListMembersRequest {
	req := apiClient.ListMembers(ctx, organizationResourceType, *model.OrganizationId)
	if model.Subject != nil {
		req = req.Subject(*model.Subject)
	}
	return req
}

func outputResult(p *print.Printer, outputFormat, sortBy string, members []authorization.Member) error {
	sortFn := func(i, j int) bool {
		switch sortBy {
		case "subject":
			return *members[i].Subject < *members[j].Subject
		case "role":
			return *members[i].Role < *members[j].Role
		default:
			return false
		}
	}
	sort.SliceStable(members, sortFn)

	return p.OutputResult(outputFormat, members, func() error {
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

		if sortBy == "subject" {
			table.EnableAutoMergeOnColumns(1)
		} else if sortBy == "role" {
			table.EnableAutoMergeOnColumns(2)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
