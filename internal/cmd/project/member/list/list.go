package list

import (
	"context"
	"fmt"
	"sort"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	authorization "github.com/stackitcloud/stackit-sdk-go/services/authorization/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/authorization/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

const (
	subjectFlag = "subject"
	limitFlag   = "limit"

	projectResourceType = "project"
)

var sortByFlag = flags.StringEnumFlag(
	"sort-by",
	[]string{"subject", "role"},
	"Sort entries by a specific field,",
	flags.StringEnumDefaultValue("subject"),
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Subject *string
	Limit   *int64
	SortBy  string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
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
			members := resp.Members

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Truncate output
			if model.Limit != nil && len(members) > int(*model.Limit) {
				members = members[:*model.Limit]
			}

			return outputResult(params.Printer, *model, projectLabel, members)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(subjectFlag, "", "Filter by subject (the identifier of a user, service account or client). This is usually the email address (for users) or name (for clients)")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	sortByFlag.Register(cmd.Flags())
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
		Subject:         flags.FlagToStringPointer(p, cmd, subjectFlag),
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
		SortBy:          sortByFlag.Get(),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiListMembersRequest {
	req := apiClient.DefaultAPI.ListMembers(ctx, projectResourceType, model.ProjectId)
	if model.Subject != nil {
		req = req.Subject(*model.Subject)
	}
	return req
}

func outputResult(p *print.Printer, model inputModel, projectLabel string, members []authorization.Member) error {
	if model.GlobalFlagModel == nil {
		return fmt.Errorf("globalflags are empty")
	}
	sortFn := func(i, j int) bool {
		switch model.SortBy {
		case "subject":
			return members[i].Subject < members[j].Subject
		case "role":
			return members[i].Role < members[j].Role
		default:
			return false
		}
	}
	sort.SliceStable(members, sortFn)

	return p.OutputResult(model.OutputFormat, members, func() error {
		if len(members) == 0 {
			p.Outputf("No members found for project %q\n", projectLabel)
		}

		table := tables.NewTable()
		table.SetHeader("SUBJECT", "ROLE")
		for i := range members {
			m := members[i]
			// If the previous item differs from the current item on the element to sort by, add a separator between the rows to help readability
			if i > 0 && sortFn(i-1, i) {
				table.AddSeparator()
			}
			table.AddRow(m.Subject, m.Role)
		}

		switch model.SortBy {
		case "subject":
			table.EnableAutoMergeOnColumns(1)
		case "role":
			table.EnableAutoMergeOnColumns(2)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
