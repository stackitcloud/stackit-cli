package list

import (
	"context"
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	resourcemanager "github.com/stackitcloud/stackit-sdk-go/services/resourcemanager/v0api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
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

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists STACKIT folders",
		Long:  "Lists all STACKIT folders that match certain criteria.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all STACKIT folders that are children of the configured project`,
				"$ stackit folder list"),
			examples.NewExample(
				`List all STACKIT folders that are children of a specific parent`,
				"$ stackit folder list --parent-id xxx"),
			examples.NewExample(
				`List all STACKIT folders that match the given project IDs, located under the same parent resource`,
				"$ stackit folder list --project-id-like xxx,yyy,zzz"),
			examples.NewExample(
				`List all STACKIT folders that a certain user is a member of`,
				"$ stackit folder list --member example@email.com"),
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

			// Fetch folders
			folders, err := fetchFolders(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				return err
			}

			return outputResult(params.Printer, model.OutputFormat, folders)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(parentIdFlag, "", "Filter by parent identifier")
	cmd.Flags().Var(flags.UUIDSliceFlag(), projectIdLikeFlag, "Filter by project identifier. Multiple project IDs can be provided, but they need to belong to the same parent resource")
	cmd.Flags().String(memberFlag, "", "Filter by member. The list of folders of which the member is part of will be shown")
	cmd.Flags().String(creationTimeAfterFlag, "", "Filter by creation timestamp, in a date-time with the RFC3339 layout format, e.g. 2023-01-01T00:00:00Z. The list of folders that were created after the given timestamp will be shown")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Int64(pageSizeFlag, pageSizeDefault, "Number of items fetched in each API call. Does not affect the number of items in the command output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient resourceManagerClient, offset int) (resourcemanager.ApiListFoldersRequest, error) {
	req := apiClient.ListFolders(ctx)
	if model.ParentId != nil {
		req = req.ContainerParentId(*model.ParentId)
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
	req = req.Limit(float32(model.PageSize))
	req = req.Offset(float32(offset))
	return req, nil
}

type resourceManagerClient interface {
	ListFolders(ctx context.Context) resourcemanager.ApiListFoldersRequest
}

func fetchFolders(ctx context.Context, model *inputModel, apiClient resourceManagerClient) ([]resourcemanager.ListFoldersResponseItemsInner, error) {
	if model.Limit != nil && *model.Limit < model.PageSize {
		model.PageSize = *model.Limit
	}

	offset := 0
	folders := []resourcemanager.ListFoldersResponseItemsInner{}
	for {
		// Call API
		req, err := buildRequest(ctx, model, apiClient, offset)
		if err != nil {
			return nil, fmt.Errorf("build list folders request: %w", err)
		}
		resp, err := req.Execute()
		if err != nil {
			return nil, fmt.Errorf("get folders: %w", err)
		}
		respFolders := resp.GetItems()
		if len(respFolders) == 0 {
			break
		}
		folders = append(folders, respFolders...)
		// Stop if no more pages
		if len(respFolders) < int(model.PageSize) {
			break
		}

		// Stop and truncate if limit is reached
		if model.Limit != nil && len(folders) >= int(*model.Limit) {
			folders = folders[:*model.Limit]
			break
		}
		offset += int(model.PageSize)
	}
	return folders, nil
}

func outputResult(p *print.Printer, outputFormat string, folders []resourcemanager.ListFoldersResponseItemsInner) error {
	return p.OutputResult(outputFormat, folders, func() error {
		if len(folders) == 0 {
			p.Outputf("No folders found matching the criteria\n")
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "PARENT ID")
		for i := range folders {
			f := folders[i]

			table.AddRow(
				f.FolderId,
				f.Name,
				f.Parent.Id,
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
