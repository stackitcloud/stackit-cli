package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	instanceIdArg = "INSTANCE_ID"

	aclParameterKey = "sgw_acl"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of a Valkey instance",
		Long:  "Shows details of a Valkey instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Valkey instance with ID "xxx"`,
				"$ stackit beta valkey instance describe xxx"),
			examples.NewExample(
				`Get details of a Valkey instance with ID "xxx" in JSON format`,
				"$ stackit beta valkey instance describe xxx --output-format json"),
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
				return fmt.Errorf("read Valkey instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *valkey.APIClient) valkey.ApiGetInstanceRequest {
	return apiClient.DefaultAPI.GetInstance(ctx, model.ProjectId, model.Region, model.InstanceId)
}

func outputResult(p *print.Printer, outputFormat string, instance *valkey.Instance) error {
	return p.OutputResult(outputFormat, instance, func() error {
		if instance == nil {
			return fmt.Errorf("no instance passed")
		}

		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(instance.InstanceId))
		table.AddSeparator()
		table.AddRow("NAME", instance.Name)
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(instance.Status))
		table.AddSeparator()
		table.AddRow("PLAN ID", instance.PlanId)
		table.AddSeparator()
		table.AddRow("PLAN NAME", instance.PlanName)
		table.AddSeparator()
		table.AddRow("OFFERING NAME", instance.OfferingName)
		table.AddSeparator()
		table.AddRow("OFFERING VERSION", instance.OfferingVersion)
		table.AddSeparator()
		table.AddRow("LAST OPERATION TYPE", string(instance.LastOperation.Type))
		table.AddSeparator()
		table.AddRow("LAST OPERATION STATE", string(instance.LastOperation.State))
		if acl, ok := instance.Parameters[aclParameterKey].(string); ok && acl != "" {
			table.AddSeparator()
			table.AddRow("ACL", acl)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
