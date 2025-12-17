package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mariadb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mariadb"
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
		Short: "Shows details  of a MariaDB instance",
		Long:  "Shows details  of a MariaDB instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a MariaDB instance with ID "xxx"`,
				"$ stackit mariadb instance describe xxx"),
			examples.NewExample(
				`Get details of a MariaDB instance with ID "xxx" in JSON format`,
				"$ stackit mariadb instance describe xxx --output-format json"),
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
				return fmt.Errorf("read MariaDB instance: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *mariadb.APIClient) mariadb.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *mariadb.Instance) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}

	return p.OutputResult(outputFormat, instance, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(instance.InstanceId))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(instance.Name))
		table.AddSeparator()
		if instance.LastOperation != nil {
			table.AddRow("LAST OPERATION TYPE", utils.PtrString(instance.LastOperation.Type))
			table.AddSeparator()
			table.AddRow("LAST OPERATION STATE", utils.PtrString(instance.LastOperation.State))
			table.AddSeparator()
		}
		table.AddRow("PLAN ID", utils.PtrString(instance.PlanId))
		// Only show ACL if it's present and not empty
		if instance.Parameters != nil {
			acl := (*instance.Parameters)[aclParameterKey]
			aclStr, ok := acl.(string)
			if ok {
				if aclStr != "" {
					table.AddSeparator()
					table.AddRow("ACL", aclStr)
				}
			}
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
