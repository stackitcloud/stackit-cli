package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/client"
	scautils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sca/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	applicationIDArg  = "APPLICATION_ID"
	environmentIDFlag = "environment-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	EnvironmentID string
	ApplicationID string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a SCA application",
		Long:  "Show details of a STACKIT Kubernetes Engine (SCA) application.",
		Args:  args.SingleArg(applicationIDArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a SCA application with ID "xxx" from an environment with ID "yyy"`,
				"$ stackit sca application describe xxx --environment-id yyy"),
			examples.NewExample(
				`Get details of all SCA application with ID "xxx" from an environment with ID "yyy" in JSON format`,
				"$ stackit sca application describe xxx --environment-id yyy --output-format json"),
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

			globalFlags := globalflags.Parse(params.Printer, cmd)
			if globalFlags.ProjectId == "" {
				return &errors.ProjectIdError{}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe SCA application: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func outputResult(p *print.Printer, outputFormat string, application *sca.Application) error {
	return p.OutputResult(outputFormat, application, func() error {
		if application == nil {
			p.Outputf("No application found")
			return nil
		}

		table := tables.NewTable()
		table.SetTitle("Application")
		table.AddRow("ID", utils.PtrString(application.Id))
		table.AddSeparator()
		table.AddRow("NAME", application.DisplayName)
		table.AddSeparator()
		table.AddRow("STATUS", scautils.ApplicationStatusToStr(application.RuntimeStatus.GetCurrentStatus()))
		table.AddSeparator()
		table.AddRow("STATE", scautils.ApplicationStateToStr(application.GetStopped()))
		table.AddSeparator()
		table.AddRow("INSTANCES", len(application.RuntimeStatus.Instances))

		containersTable := tables.NewTable()
		containersTable.SetTitle("Application Containers")
		containersTable.SetHeader("NAME", "IMAGE", "CPU", "MEMORY")
		for _, c := range application.Containers {
			containersTable.AddRow(
				c.Name,
				c.Image,
				*c.Cpu,
				*c.Memory,
			)
		}
		err := tables.DisplayTables(p, []tables.Table{table, containersTable})
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), environmentIDFlag, "Environment ID (if not set uses default environment)")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	applicationID := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	environmentID := flags.FlagToStringValue(p, cmd, environmentIDFlag)
	if environmentID == "" {
		environmentID = globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		EnvironmentID:   environmentID,
		ApplicationID:   applicationID,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sca.APIClient) sca.ApiGetApplicationRequest {
	return apiClient.DefaultAPI.GetApplication(ctx, model.ProjectId, model.EnvironmentID, model.ApplicationID)
}
