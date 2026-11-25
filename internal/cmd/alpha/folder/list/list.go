package list

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/core/experimental/paginate"
	resourcemanager "github.com/stackitcloud/stackit-sdk-go/services/resourcemanager/v0api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	internalerrors "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	creationTimeAfter, err := flags.FlagToDateTimePointer(p, cmd, creationTimeAfterFlag, creationTimeAfterFormat)
	if err != nil {
		return nil, &internalerrors.FlagValidationError{
			Flag:    creationTimeAfterFlag,
			Details: err.Error(),
		}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &internalerrors.FlagValidationError{
			Flag:    limitFlag,
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
		PageSize:          pageSizeDefault,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient resourcemanager.DefaultAPI) (resourcemanager.ApiListFoldersRequest, error) {
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
	if model.ProjectIdLike != nil {
		req = req.ContainerIds(model.ProjectIdLike)
	}

	if model.ParentId == nil && model.ProjectIdLike == nil && model.Member == nil {
		email, err := auth.GetAuthEmail()
		if err != nil {
			return req, fmt.Errorf("get email of authenticated user: %w", err)
		}
		req = req.Member(email)
	}
	req = req.Limit(float32(model.PageSize))
	return req, nil
}

func fetchFolders(ctx context.Context, model *inputModel, apiClient resourcemanager.DefaultAPI) ([]resourcemanager.ListFoldersResponseItemsInner, error) {
	opts := []paginate.Option{paginate.WithPageSize(pageSizeDefault)}
	if model.Limit != nil {
		opts = append(opts, paginate.WithLimit(int(*model.Limit)))
	}

	req, err := buildRequest(ctx, model, apiClient)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	items := make([]resourcemanager.ListFoldersResponseItemsInner, 0)
	for item, err := range paginate.Items(&listFoldersAdapter{request: req}, opts...) {
		if err != nil {
			return nil, fmt.Errorf("list items: %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func outputResult(p *print.Printer, outputFormat string, folders []resourcemanager.ListFoldersResponseItemsInner) error {
	return p.OutputResult(outputFormat, folders, func() error {
		if len(folders) == 0 {
			p.Outputf("No folders found matching the criteria\n")
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "PARENT ID")
		for _, f := range folders {
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

var (
	_ paginate.Request[*listFoldersAdapter, listFoldersAdapterResponse] = (*listFoldersAdapter)(nil)
	_ paginate.Response[resourcemanager.ListFoldersResponseItemsInner]  = listFoldersAdapterResponse{}
)

// listFoldersAdapter maps the token-based paginate API onto the
// offset/limit-based ListFolders endpoint. The page token is the
// decimal offset of the next page.
type listFoldersAdapter struct {
	request  resourcemanager.ApiListFoldersRequest
	pageSize int32
	offset   int64
	err      error
}

func (a *listFoldersAdapter) PageSize(pageSize int32) *listFoldersAdapter {
	a.pageSize = pageSize
	a.request = a.request.Limit(float32(pageSize))
	return a
}

func (a *listFoldersAdapter) PageToken(pageToken string) *listFoldersAdapter {
	offset, err := strconv.ParseInt(pageToken, 10, 64)
	if err != nil || offset < 0 {
		a.err = fmt.Errorf("invalid page token %q", pageToken)
		return a
	}
	a.offset = offset
	a.request = a.request.Offset(float32(offset))
	return a
}

func (a *listFoldersAdapter) Execute() (listFoldersAdapterResponse, error) {
	if a.err != nil {
		return listFoldersAdapterResponse{}, a.err
	}
	resp, err := a.request.Execute()
	if err != nil {
		return listFoldersAdapterResponse{}, err
	}
	if resp == nil {
		return listFoldersAdapterResponse{}, errors.New("empty response from API")
	}
	return listFoldersAdapterResponse{
		ListFoldersResponse: resp,
		offset:              a.offset,
		pageSize:            a.pageSize,
	}, nil
}

type listFoldersAdapterResponse struct {
	*resourcemanager.ListFoldersResponse
	offset   int64
	pageSize int32
}

func (r listFoldersAdapterResponse) GetNextPageToken() string {
	n := int64(len(r.GetItems()))
	if n == 0 || n < int64(r.pageSize) {
		return ""
	}
	return strconv.FormatInt(r.offset+n, 10)
}
