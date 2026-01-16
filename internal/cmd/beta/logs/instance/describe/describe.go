package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	argInstanceID = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceID string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", argInstanceID),
		Short: "Shows details of a Logs instance",
		Long:  "Shows details of a Logs instance",
		Args:  args.SingleArg(argInstanceID, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Logs instance with ID "xxx"`,
				`$ stackit beta logs instance describe xxx`,
			),
			examples.NewExample(
				`Get details of a Logs instance with ID "xxx" in JSON format`,
				"$ stackit beta logs instance describe xxx --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	model := &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceID:      inputArgs[0],
	}
	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logs.APIClient) logs.ApiGetLogsInstanceRequest {
	return apiClient.GetLogsInstance(ctx, model.ProjectId, model.Region, model.InstanceID)
}

func outputResult(p *print.Printer, outputFormat string, instance *logs.LogsInstance) error {
	if instance == nil {
		return fmt.Errorf("instance response is empty")
	}
	return p.OutputResult(outputFormat, instance, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(instance.Id))
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", utils.PtrString(instance.DisplayName))
		table.AddSeparator()
		table.AddRow("RETENTION DAYS", utils.PtrString(instance.RetentionDays))
		table.AddSeparator()
		table.AddRow("ACL IP RANGES", utils.PtrString(instance.Acl))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display table: %w", err)
		}
		return nil
	})
}
