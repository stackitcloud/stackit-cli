package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
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
	sca "github.com/stackitcloud/stackit-sdk-go/services/sca/v1alphaapi"
)

const (
	limitFlag         = "limit"
	environmentIDFlag = "environment-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID string
	Limit         *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all SCA applications",
		Long:  "Lists all STACKIT Kubernetes Engine (SCA) applications.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all SCA applications`,
				"$ stackit sca application list"),
			examples.NewExample(
				`List all SCA applications from enviroment with ID "xxx"`,
				"$ stackit sca application list --environment-id xxx"),
			examples.NewExample(
				`List all SCA applications in JSON format`,
				"$ stackit sca application list --output-format json"),
			examples.NewExample(
				`List up to 10 SCA applications`,
				"$ stackit sca application list --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			applications, err := makeRequest(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("list SCA applications: %w", err)
			}

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

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		EnvironmentID:   flags.FlagToStringValue(p, cmd, environmentIDFlag),
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func makeRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) ([]sca.ApplicationSummary, error) {
	// If environment ID is not defined, return all applications from the project
	if model.EnvironmentID == "" {
		resp, err := apiClient.DefaultAPI.ListProjectApplications(ctx, model.ProjectId).Execute()
		if err != nil {
			return nil, err
		}

		return resp.Items, err
	}

	resp, err := apiClient.DefaultAPI.ListApplications(ctx, model.ProjectId, model.EnvironmentID).Execute()
	if err != nil {
		return nil, err
	}

	return resp.Items, err
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
				scautils.ApplicationStatusToStr(a.GetStatus()),
				scautils.ApplicationStateToStr(a.GetStopped()),
				// state,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
