package list

import (
	"context"
	"fmt"
	"math"

	"github.com/spf13/cobra"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/registry/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	limitFlag   = "limit"
	maxPageSize = 100
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Container Registry artifactories",
		Long:  "Lists all Container Registry artifactories.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all artifactories`,
				"$ stackit registry artifactory list"),
			examples.NewExample(
				`List all artifactories in JSON format`,
				"$ stackit registry artifactory list --output-format json"),
			examples.NewExample(
				`List up to 10 artifactories`,
				"$ stackit registry artifactory list --limit 10"),
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
			artifactories, err := fetchArtifactories(ctx, model, apiClient)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, artifactories)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
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
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *registry.APIClient, pageToken string, pageSize int32) registry.ApiListArtifactoriesRequest {
	req := apiClient.DefaultAPI.ListArtifactories(ctx, model.ProjectId, model.Region)
	if pageToken != "" {
		req = req.NextPageToken(pageToken)
	}
	if pageSize > 0 {
		req = req.PageSize(pageSize)
	}
	return req
}

func fetchArtifactories(ctx context.Context, model *inputModel, apiClient *registry.APIClient) ([]registry.Artifactory, error) {
	var pageToken string
	artifactories := make([]registry.Artifactory, 0)
	received := int64(0)
	limit := int64(math.MaxInt64)
	if model.Limit != nil {
		limit = *model.Limit
	}
	for {
		want := min(int64(maxPageSize), limit-received)
		request := buildRequest(ctx, model, apiClient, pageToken, int32(want))
		response, err := request.Execute()
		if err != nil {
			return nil, fmt.Errorf("get artifactories: %w", err)
		}
		if response == nil {
			return nil, fmt.Errorf("get artifactories: empty response")
		}
		if response.Artifactories != nil {
			artifactories = append(artifactories, response.Artifactories...)
		}
		pageToken = ""
		if response.NextPageToken != nil {
			pageToken = *response.NextPageToken
		}
		received += want
		if pageToken == "" || received >= limit {
			break
		}
	}
	if model.Limit != nil && len(artifactories) > int(*model.Limit) {
		artifactories = artifactories[:*model.Limit]
	}
	return artifactories, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, artifactories []registry.Artifactory) error {
	return p.OutputResult(outputFormat, artifactories, func() error {
		if len(artifactories) == 0 {
			p.Outputf("No artifactories found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATE", "URL")
		for i := range artifactories {
			artifactory := artifactories[i]
			table.AddRow(
				artifactory.Id,
				artifactory.Name,
				string(artifactory.State),
				artifactory.Url,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
