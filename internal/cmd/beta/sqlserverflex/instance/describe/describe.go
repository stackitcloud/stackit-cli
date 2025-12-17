package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details  of a SQLServer Flex instance",
		Long:  "Shows details  of a SQLServer Flex instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a SQLServer Flex instance with ID "xxx"`,
				"$ stackit beta sqlserverflex instance describe xxx"),
			examples.NewExample(
				`Get details of a SQLServer Flex instance with ID "xxx" in JSON format`,
				"$ stackit beta sqlserverflex instance describe xxx --output-format json"),
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
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read SQLServer Flex instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp.Item)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sqlserverflex.APIClient) sqlserverflex.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *sqlserverflex.Instance) error {
	if instance == nil {
		return fmt.Errorf("instance response is empty")
	}

	return p.OutputResult(outputFormat, instance, func() error {
		var acls string
		if instance.Acl != nil && instance.Acl.HasItems() {
			aclsArray := *instance.Acl.Items
			acls = strings.Join(aclsArray, ",")
		}

		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(instance.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(instance.Name))
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(instance.Status))
		table.AddSeparator()
		if instance.Storage != nil {
			table.AddRow("STORAGE SIZE (GB)", utils.PtrString(instance.Storage.Size))
			table.AddSeparator()
		}
		table.AddRow("VERSION", utils.PtrString(instance.Version))
		table.AddSeparator()
		table.AddRow("BACKUP SCHEDULE (UTC)", utils.PtrString(instance.BackupSchedule))
		table.AddSeparator()
		table.AddRow("ACL", acls)
		table.AddSeparator()
		if instance.Flavor != nil {
			table.AddRow("FLAVOR DESCRIPTION", utils.PtrString(instance.Flavor.Description))
			table.AddSeparator()
			table.AddRow("CPU", utils.PtrString(instance.Flavor.Cpu))
			table.AddSeparator()
			table.AddRow("RAM (GB)", utils.PtrString(instance.Flavor.Memory))
			table.AddSeparator()
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
