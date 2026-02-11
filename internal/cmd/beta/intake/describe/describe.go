package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	intakeIdArg = "INTAKE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", intakeIdArg),
		Short: "Shows details of an Intake",
		Long:  "Shows details of an Intake.",
		Args:  args.SingleArg(intakeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an Intake with ID "xxx"`,
				`$ stackit beta intake describe xxx`),
			examples.NewExample(
				`Get details of an Intake with ID "xxx" in JSON format`,
				`$ stackit beta intake describe xxx --output-format json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get Intake: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	intakeId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		IntakeId:        intakeId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiGetIntakeRequest {
	req := apiClient.GetIntake(ctx, model.ProjectId, model.Region, model.IntakeId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, intk *intake.IntakeResponse) error {
	if intk == nil {
		return fmt.Errorf("received nil response, could not display details")
	}

	return p.OutputResult(outputFormat, intk, func() error {
		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")

		table.AddRow("ID", intk.GetId())
		table.AddRow("Name", intk.GetDisplayName())
		table.AddRow("State", intk.GetState())
		table.AddRow("Runner ID", intk.GetIntakeRunnerId())
		table.AddRow("Created", intk.GetCreateTime())
		table.AddRow("Labels", intk.GetLabels())

		if description := intk.GetDescription(); description != "" {
			table.AddRow("Description", description)
		}

		if failureMessage := intk.GetFailureMessage(); failureMessage != "" {
			table.AddRow("Failure Message", failureMessage)
		}

		table.AddSeparator()
		table.AddRow("Ingestion URI", intk.GetUri())
		table.AddRow("Topic", intk.GetTopic())
		table.AddRow("Dead Letter Topic", intk.GetDeadLetterTopic())
		table.AddRow("Undelivered Messages", intk.GetUndeliveredMessageCount())

		table.AddSeparator()
		catalog := intk.GetCatalog()
		table.AddRow("Catalog URI", catalog.GetUri())
		table.AddRow("Catalog Warehouse", catalog.GetWarehouse())
		if namespace := catalog.GetNamespace(); namespace != "" {
			table.AddRow("Catalog Namespace", namespace)
		}
		if tableName := catalog.GetTableName(); tableName != "" {
			table.AddRow("Catalog Table Name", tableName)
		}
		table.AddRow("Catalog Partitioning", catalog.GetPartitioning())
		if partitionBy := catalog.GetPartitionBy(); partitionBy != nil && len(*partitionBy) > 0 {
			table.AddRow("Catalog Partition By", strings.Join(*partitionBy, ", "))
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
