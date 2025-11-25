package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	intakeIdArg = "INTAKE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	IntakeId string
}

func NewCmd(p *params.CmdParams) *cobra.Command {
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

func outputResult(p *print.Printer, outputFormat string, intake *intake.IntakeResponse) error {
	if intake == nil {
		return fmt.Errorf("received nil response, could not display details")
	}

	return p.OutputResult(outputFormat, intake, func() error {
		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")

		table.AddRow("ID", intake.GetId())
		table.AddRow("Name", intake.GetDisplayName())
		table.AddRow("State", intake.GetState())
		table.AddRow("Runner ID", intake.GetIntakeRunnerId())
		table.AddRow("Created", intake.GetCreateTime())
		table.AddRow("Labels", intake.GetLabels())

		if description := intake.GetDescription(); description != "" {
			table.AddRow("Description", description)
		}

		if failureMessage := intake.GetFailureMessage(); failureMessage != "" {
			table.AddRow("Failure Message", failureMessage)
		}

		table.AddSeparator()
		table.AddRow("Ingestion URI", intake.GetUri())
		table.AddRow("Topic", intake.GetTopic())
		table.AddRow("Dead Letter Topic", intake.GetDeadLetterTopic())
		table.AddRow("Undelivered Messages", intake.GetUndeliveredMessageCount())

		table.AddSeparator()
		catalog := intake.GetCatalog()
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
