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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SecurityGroupId string
}

const groupIdArg = "GROUP_ID"

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", groupIdArg),
		Short: "Describes security groups",
		Long:  "Describes security groups by its internal ID.",
		Args:  args.SingleArg(groupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Describe group "xxx"`, `$ stackit beta security-group describe xxx`),
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
			request := buildRequest(ctx, model, apiClient)

			group, err := request.Execute()
			if err != nil {
				return fmt.Errorf("get security group: %w", err)
			}

			if err := outputResult(p, model.OutputFormat, group); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetSecurityGroupRequest {
	request := apiClient.GetSecurityGroup(ctx, model.ProjectId, model.SecurityGroupId)
	return request
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SecurityGroupId: cliArgs[0],
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

func outputResult(p *print.Printer, outputFormat string, resp *iaas.SecurityGroup) error {
	if resp == nil {
		return fmt.Errorf("security group response is empty")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal security group: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal security group: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		var content []tables.Table

		table := tables.NewTable()
		table.SetTitle("SECURITY GROUP")

		if id := resp.Id; id != nil {
			table.AddRow("ID", *id)
		}
		table.AddSeparator()

		if name := resp.Name; name != nil {
			table.AddRow("NAME", *name)
			table.AddSeparator()
		}

		if description := resp.Description; description != nil {
			table.AddRow("DESCRIPTION", *description)
			table.AddSeparator()
		}

		if stateful := resp.Stateful; stateful != nil {
			table.AddRow("STATEFUL", *stateful)
			table.AddSeparator()
		}

		if resp.Labels != nil && len(*resp.Labels) > 0 {
			var labels []string
			for key, value := range *resp.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		if resp.CreatedAt != nil {
			table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(resp.CreatedAt))
			table.AddSeparator()
		}

		if resp.UpdatedAt != nil {
			table.AddRow("UPDATED AT", utils.ConvertTimePToDateTimeString(resp.UpdatedAt))
			table.AddSeparator()
		}

		content = append(content, table)

		if resp.Rules != nil && len(*resp.Rules) > 0 {
			rulesTable := tables.NewTable()
			rulesTable.SetTitle("RULES")
			rulesTable.SetHeader(
				"ID",
				"DESCRIPTION",
				"PROTOCOL",
				"DIRECTION",
				"ETHER TYPE",
				"PORT RANGE",
				"IP RANGE",
				"ICMP PARAMETERS",
				"REMOTE SECURITY GROUP ID",
			)

			for _, rule := range *resp.Rules {
				var portRange string
				if rule.PortRange != nil {
					portRange = fmt.Sprintf("%s-%s", utils.PtrString(rule.PortRange.Min), utils.PtrString(rule.PortRange.Max))
				}

				var protocol string
				if rule.Protocol != nil {
					protocol = utils.PtrString(rule.Protocol.Name)
				}

				var icmpParameter string
				if rule.IcmpParameters != nil {
					icmpParameter = fmt.Sprintf("type: %s, code: %s", utils.PtrString(rule.IcmpParameters.Type), utils.PtrString(rule.IcmpParameters.Code))
				}

				rulesTable.AddRow(
					utils.PtrString(rule.Id),
					utils.PtrString(rule.Description),
					protocol,
					utils.PtrString(rule.Direction),
					utils.PtrString(rule.Ethertype),
					portRange,
					utils.PtrString(rule.IpRange),
					icmpParameter,
					utils.PtrString(rule.RemoteSecurityGroupId),
				)
			}

			content = append(content, rulesTable)
		}

		if err := tables.DisplayTables(p, content); err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
