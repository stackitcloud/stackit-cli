package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of a PostgreSQL Flex instance",
		Long:  "Shows details of a PostgreSQL Flex instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a PostgreSQL Flex instance with ID "xxx"`,
				"$ stackit postgresflex instance describe xxx"),
			examples.NewExample(
				`Get details of a PostgreSQL Flex instance with ID "xxx" in JSON format`,
				"$ stackit postgresflex instance describe xxx --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read PostgreSQL Flex instance: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp.Item)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *postgresflex.APIClient) postgresflex.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *postgresflex.Instance) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(instance, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL Flex instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		aclsArray := *instance.Acl.Items
		acls := strings.Join(aclsArray, ",")

		instanceType, err := postgresflexUtils.GetInstanceType(*instance.Replicas)
		if err != nil {
			// Should never happen
			instanceType = ""
		}

		table := tables.NewTable()
		table.AddRow("ID", *instance.Id)
		table.AddSeparator()
		table.AddRow("NAME", *instance.Name)
		table.AddSeparator()
		table.AddRow("STATUS", cases.Title(language.English).String(*instance.Status))
		table.AddSeparator()
		table.AddRow("STORAGE SIZE", *instance.Storage.Size)
		table.AddSeparator()
		table.AddRow("VERSION", *instance.Version)
		table.AddSeparator()
		table.AddRow("ACL", acls)
		table.AddSeparator()
		table.AddRow("FLAVOR DESCRIPTION", *instance.Flavor.Description)
		table.AddSeparator()
		table.AddRow("TYPE", instanceType)
		table.AddSeparator()
		table.AddRow("REPLICAS", *instance.Replicas)
		table.AddSeparator()
		table.AddRow("CPU", *instance.Flavor.Cpu)
		table.AddSeparator()
		table.AddRow("RAM", *instance.Flavor.Memory)
		table.AddSeparator()
		table.AddRow("BACKUP SCHEDULE (UTC)", *instance.BackupSchedule)
		table.AddSeparator()
		err = table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
