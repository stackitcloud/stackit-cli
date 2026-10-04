package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/client"
	scautils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	limitFlag         = "limit"
	environmentIDFlag = "environment-id"
)

type listApplicationsRequest interface {
	Execute() (*sca.ListApplicationsResponse, error)
}

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID string
	Limit         *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all SCA applications",
		Long:  "Lists all STACKIT Container Applications (SCA) applications.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all SCA applications`,
				"$ stackit alpha sca application list"),
			examples.NewExample(
				`List all SCA applications from environment with ID "xxx"`,
				"$ stackit alpha sca application list --environment-id xxx"),
			examples.NewExample(
				`List all SCA applications in JSON format`,
				"$ stackit alpha sca application list --output-format json"),
			examples.NewExample(
				`List up to 10 SCA applications`,
				"$ stackit alpha sca application list --limit 10"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, nil)
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
				return fmt.Errorf("list SCA applications: %w", err)
			}

			applications := resp.Items

			// Truncate output
			if model.Limit != nil && len(applications) > int(*model.Limit) {
				applications = applications[:*model.Limit]
			}

			projectLabel := model.ProjectId
			if len(applications) == 0 {
				projectLabel, err = projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, applications)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), environmentIDFlag, "Optional environment ID to filter applications")
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
		EnvironmentID:   flags.FlagToStringValue(p, cmd, environmentIDFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) listApplicationsRequest {
	// If environment ID is not defined, return all applications from the project
	if model.EnvironmentID == "" {
		return apiClient.DefaultAPI.ListProjectApplications(ctx, model.ProjectId)
	}

	return apiClient.DefaultAPI.ListApplications(ctx, model.ProjectId, model.EnvironmentID)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, applications []sca.ApplicationSummary) error {
	return p.OutputResult(outputFormat, applications, func() error {
		if len(applications) == 0 {
			p.Outputf("No applications found for project %q\n", projectLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "INSTANCES", "ENVIRONMENT", "ENVIRONMENT ID", "URL", "STATUS", "STATE")
		for _, a := range applications {
			table.AddRow(
				a.GetId(),
				a.GetDisplayName(),
				a.GetInstances(),
				a.GetEnvironmentName(),
				a.GetEnvironmentId(),
				a.GetUrl(),
				scautils.ApplicationStatusToStr(a.Status),
				scautils.ApplicationStateToStr(a.GetStopped()),
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
