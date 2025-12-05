package list

import (
	"cmp"
	"context"
	"fmt"
	"path"
	"slices"
	"sync"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"golang.org/x/sync/errgroup"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	authclient "github.com/stackitcloud/stackit-cli/internal/pkg/services/authorization/client"
	resourceclient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
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

func NewCmd(params *types.CmdParams) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			resourceClient, err := resourceclient.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			authClient, err := authclient.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Fetch projects
			projects, err := fetchProjects(cmd.Context(), model, resourceClient, authClient)
			if err != nil {
				return err
			}
			if len(projects) == 0 {
				params.Printer.Info("No projects found matching the criteria\n")
				return nil
			}

			return outputResult(params.Printer, model.OutputFormat, projects)
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

type project struct {
	Name         string
	ID           string
	Organization string
	Folder       []string
}

func (p project) FolderPath() string {
	return path.Join(p.Folder...)
}

func getProjects(ctx context.Context, parent *node, org string, projChan chan<- project) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, child := range parent.children {
		g.Go(func() error {
			if child.typ != resourceTypeProject {
				return getProjects(ctx, child, org, projChan)
			}
			parent := child.parent
			folderName := []string{}
			for parent != nil {
				if parent.typ == resourceTypeFolder {
					folderName = append([]string{parent.name}, folderName...)
				}
				parent = parent.parent
			}
			projChan <- project{
				Name:         child.name,
				ID:           child.resourceID,
				Organization: org,
				Folder:       folderName,
			}
			return nil
		})
	}
	return g.Wait()
}

type resourceManagerClient interface {
	ListProjects(ctx context.Context) resourcemanager.ApiListProjectsRequest
}

func fetchProjects(ctx context.Context, model *inputModel, resourceClient *resourcemanager.APIClient, authClient *authorization.APIClient) ([]project, error) {
	tree, err := newResourceTree(resourceClient, authClient, model)
	if err != nil {
		return nil, err
	}

	if err := tree.Fill(ctx); err != nil {
		return nil, err
	}

	var projs []project
	projChan := make(chan project)

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		for p := range projChan {
			i, _ := slices.BinarySearchFunc(projs, p, func(e project, target project) int {
				if orgCmp := cmp.Compare(e.Organization, target.Organization); orgCmp != 0 {
					return orgCmp
				}
				return cmp.Compare(e.FolderPath(), p.FolderPath())
			})
			projs = slices.Insert(projs, i, p)
		}
	}()

	for _, root := range tree.roots {
		if err := getProjects(ctx, root, root.name, projChan); err != nil {
			return nil, err
		}
	}
	close(projChan)
	wg.Wait()
	return projs, nil
}

func outputResult(p *print.Printer, outputFormat string, projects []project) error {
	return p.OutputResult(outputFormat, projects, func() error {
		table := tables.NewTable()
		table.SetHeader("ORGANIZATION", "FOLDER", "NAME", "ID")
		for i := range projects {
			p := projects[i]
			table.AddRow(
				p.Organization,
				p.FolderPath(),
				p.Name,
				p.ID,
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
