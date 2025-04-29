package list

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
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
	parentIdFlag          = "parent-id"
	projectIdLikeFlag     = "project-id-like"
	memberFlag            = "member"
	creationTimeAfterFlag = "creation-time-after"
	limitFlag             = "limit"
	pageSizeFlag          = "page-size"

	creationTimeAfterFormat = time.RFC3339
	pageSizeDefault         = 50
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ParentId          *string
	ProjectIdLike     []string
	Member            *string
	CreationTimeAfter *time.Time
	Limit             *int64
	PageSize          int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists STACKIT projects",
		Long:  "Lists all STACKIT projects that match certain criteria.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all STACKIT projects that the authenticated user or service account is a member of`,
				"$ stackit project list"),
			examples.NewExample(
				`List all STACKIT projects that are children of a specific parent`,
				"$ stackit project list --parent-id xxx"),
			examples.NewExample(
				`List all STACKIT projects that match the given project IDs, located under the same parent resource`,
				"$ stackit project list --project-id-like xxx,yyy,zzz"),
			examples.NewExample(
				`List all STACKIT projects that a certain user is a member of`,
				"$ stackit project list --member example@email.com"),
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

			// Fetch projects
			projects, err := fetchProjects(ctx, model, apiClient)
			if err != nil {
				return err
			}
			if len(projects) == 0 {
				p.Info("No projects found matching the criteria\n")
				return nil
			}

			return outputResult(p, model.OutputFormat, projects)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(parentIdFlag, "", "Filter by parent identifier")
	cmd.Flags().Var(flags.UUIDSliceFlag(), projectIdLikeFlag, "Filter by project identifier. Multiple project IDs can be provided, but they need to belong to the same parent resource")
	cmd.Flags().String(memberFlag, "", "Filter by member. The list of projects of which the member is part of will be shown")
	cmd.Flags().String(creationTimeAfterFlag, "", "Filter by creation timestamp, in a date-time with the RFC3339 layout format, e.g. 2023-01-01T00:00:00Z. The list of projects that were created after the given timestamp will be shown")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Int64(pageSizeFlag, pageSizeDefault, "Number of items fetched in each API call. Does not affect the number of items in the command output")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	creationTimeAfter, err := flags.FlagToDateTimePointer(p, cmd, creationTimeAfterFlag, creationTimeAfterFormat)
	if err != nil {
		return nil, &errors.FlagValidationError{
			Flag:    creationTimeAfterFlag,
			Details: err.Error(),
		}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	pageSize := flags.FlagWithDefaultToInt64Value(p, cmd, pageSizeFlag)
	if pageSize < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    pageSizeFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		ParentId:          flags.FlagToStringPointer(p, cmd, parentIdFlag),
		ProjectIdLike:     flags.FlagToStringSliceValue(p, cmd, projectIdLikeFlag),
		Member:            flags.FlagToStringPointer(p, cmd, memberFlag),
		CreationTimeAfter: creationTimeAfter,
		Limit:             limit,
		PageSize:          pageSize,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient resourceManagerClient, offset int) (resourcemanager.ApiListProjectsRequest, error) {
	req := apiClient.ListProjects(ctx)
	if model.ParentId != nil {
		req = req.ContainerParentId(*model.ParentId)
	}
	if model.ProjectIdLike != nil {
		req = req.ContainerIds(model.ProjectIdLike)
	}
	if model.Member != nil {
		req = req.Member(*model.Member)
	}
	if model.CreationTimeAfter != nil {
		req = req.CreationTimeStart(*model.CreationTimeAfter)
	}

	if model.ParentId == nil && model.ProjectIdLike == nil && model.Member == nil {
		email, err := auth.GetAuthEmail()
		if err != nil {
			return req, fmt.Errorf("get email of authenticated user: %w", err)
		}
		req = req.Member(email)
	}

	req = req.Limit(float32(*model.Limit))
	req = req.Offset(float32(offset))
	return req, nil
}

type resourceManagerClient interface {
	ListProjects(ctx context.Context) resourcemanager.ApiListProjectsRequest
}

func fetchProjects(ctx context.Context, model *inputModel, apiClient resourceManagerClient) ([]resourcemanager.Project, error) {
	if model.Limit != nil && *model.Limit < model.PageSize {
		model.PageSize = *model.Limit
	}

	offset := 0
	projects := []resourcemanager.Project{}
	for {
		// Call API
		req, err := buildRequest(ctx, model, apiClient, offset)
		if err != nil {
			return nil, fmt.Errorf("build list projects request: %w", err)
		}
		resp, err := req.Execute()
		if err != nil {
			return nil, fmt.Errorf("get projects: %w", err)
		}
		respProjects := *resp.Items
		if len(respProjects) == 0 {
			break
		}
		projects = append(projects, respProjects...)
		// Stop if no more pages
		if len(respProjects) < int(model.PageSize) {
			break
		}

		// Stop and truncate if limit is reached
		if model.Limit != nil && len(projects) >= int(*model.Limit) {
			projects = projects[:*model.Limit]
			break
		}
		offset += int(model.PageSize)
	}
	return projects, nil
}

func outputResult(p *print.Printer, outputFormat string, projects []resourcemanager.Project) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(projects, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal projects list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(projects, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal projects list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATE", "PARENT ID")
		for i := range projects {
			p := projects[i]

			var parentId *string
			if p.Parent != nil {
				parentId = p.Parent.Id
			}
			table.AddRow(
				utils.PtrString(p.ProjectId),
				utils.PtrString(p.Name),
				utils.PtrString(p.LifecycleState),
				utils.PtrString(parentId),
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
