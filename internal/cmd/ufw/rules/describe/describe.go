package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	ufw "github.com/stackitcloud/stackit-sdk-go/services/ufw/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ufw/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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
		Short: "Shows details of a UFW rule instance",
		Long:  "Shows details of a STACKIT Unified Firewall (UFW) rule instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a UFW rule instance with ID "xxx"`,
				"$ stackit ufw rules describe xxx"),
			examples.NewExample(
				`Get details of a UFW rule instance with ID "xxx" in JSON format`,
				"$ stackit ufw rules describe xxx --output-format json"),
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
				return fmt.Errorf("read UFW rule instance: %w", err)
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

	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ufw.APIClient) ufw.ApiGetRuleRequest {
	return apiClient.DefaultAPI.GetRule(ctx, model.ProjectId, model.Region, model.InstanceId)
}

func outputResult(p *print.Printer, outputFormat string, rule *ufw.RuleResponse) error {
	return p.OutputResult(outputFormat, rule, func() error {
		if rule == nil {
			return fmt.Errorf("no instance rule passed")
		}

		table := tables.NewTable()
		table.AddRow("PRODUCT", utils.PtrString(&rule.Product))
		table.AddSeparator()
		table.AddRow("SOURCE", utils.PtrString(&rule.SourceIP))
		table.AddSeparator()
		table.AddRow("DEPLOYMENT TARGET", utils.PtrString(&rule.InstanceName))
		table.AddSeparator()
		table.AddRow("PROTOCOL", utils.PtrString(&rule.Protocol))
		table.AddSeparator()
		table.AddRow("DIRECTION", utils.PtrString(&rule.Direction))
		table.AddSeparator()
		table.AddRow("PORT RANGE", fmt.Sprintf("%d", rule.PortRange))
		table.AddSeparator()
		table.AddRow("ETHER TYPE", utils.PtrString(&rule.EtherType))
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(&rule.Status))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
