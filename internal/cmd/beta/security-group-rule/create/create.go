package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	securityGroupIdFlag       = "security-group-id"
	directionFlag             = "direction"
	descriptionFlag           = "description"
	etherTypeFlag             = "ether-type"
	icmpParameterCodeFlag     = "icmp-parameter-code"
	icmpParameterTypeFlag     = "icmp-parameter-type"
	ipRangeFlag               = "ip-range"
	portRangeMaxFlag          = "port-range-max"
	portRangeMinFlag          = "port-range-min"
	remoteSecurityGroupIdFlag = "remote-security-group-id"
	protocolNumberFlag        = "protocol-number"
	protocolNameFlag          = "protocol-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SecurityGroupId       string
	Direction             *string
	Description           *string
	EtherType             *string
	IcmpParameterCode     *int64
	IcmpParameterType     *int64
	IpRange               *string
	PortRangeMax          *int64
	PortRangeMin          *int64
	RemoteSecurityGroupId *string
	ProtocolNumber        *int64
	ProtocolName          *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a security group rule",
		Long:  "Creates a security group rule.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a security group rule for security group with ID "xxx" with direction "ingress"`,
				`$ stackit beta security-group-rule create --security-group-id xxx --direction ingress`,
			),
			examples.NewExample(
				`Create a security group rule for security group with ID "xxx" with direction "egress", protocol "icmp" and icmp parameters`,
				`$ stackit beta security-group-rule create --security-group-id xxx --direction egress --protocol icmp --icmp-parameter-code 0 --icmp-parameter-type 8`,
			),
			examples.NewExample(
				`Create a security group rule for security group with ID "xxx" with direction "ingress" and port range values`,
				`$ stackit beta security-group-rule create --security-group-id xxx --direction ingress --port-range-max 24 --port-range-min 22`,
			),
			examples.NewExample(
				`Create a security group rule for security group with ID "xxx" with direction "ingress" and protocol number 1 `,
				`$ stackit beta security-group-rule create --security-group-id xxx --direction ingress --protocol-number 1`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			securityGroupLabel, err := iaasUtils.GetSecurityGroupName(ctx, apiClient, model.ProjectId, model.SecurityGroupId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get security group name: %v", err)
				securityGroupLabel = model.SecurityGroupId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a security group rule for security group %q for project %q?", securityGroupLabel, projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create security group rule : %w", err)
			}

			return outputResult(p, model, projectLabel, securityGroupLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(securityGroupIdFlag, "", "The security group ID")
	cmd.Flags().String(directionFlag, "", "The direction of the traffic which the rule should match. The possible values are: `ingress`, `egress`")
	cmd.Flags().String(descriptionFlag, "", "The rule description")
	cmd.Flags().String(etherTypeFlag, "", "The ethertype which the rule should match")
	cmd.Flags().Int64(icmpParameterCodeFlag, 0, "ICMP code. Can be set if the protocol is ICMP")
	cmd.Flags().Int64(icmpParameterTypeFlag, 0, "ICMP type. Can be set if the protocol is ICMP")
	cmd.Flags().String(ipRangeFlag, "", "The remote IP range which the rule should match")
	cmd.Flags().Int64(portRangeMaxFlag, 0, "The maximum port number. Should be greater or equal to the minimum. This should only be provided if the protocol is not ICMP")
	cmd.Flags().Int64(portRangeMinFlag, 0, "The minimum port number. Should be less or equal to the maximum. This should only be provided if the protocol is not ICMP")
	cmd.Flags().Var(flags.UUIDFlag(), remoteSecurityGroupIdFlag, "The remote security group which the rule should match")
	cmd.Flags().Int64(protocolNumberFlag, 0, "The protocol number which the rule should match. Either `name` or `number` must be provided")
	cmd.Flags().String(protocolNameFlag, "", "The protocol name which the rule should match. Either `name` or `number` must be provided")

	err := flags.MarkFlagsRequired(cmd, securityGroupIdFlag, directionFlag)
	cmd.MarkFlagsMutuallyExclusive(protocolNumberFlag, protocolNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:       globalFlags,
		SecurityGroupId:       flags.FlagToStringValue(p, cmd, securityGroupIdFlag),
		Direction:             flags.FlagToStringPointer(p, cmd, directionFlag),
		Description:           flags.FlagToStringPointer(p, cmd, descriptionFlag),
		EtherType:             flags.FlagToStringPointer(p, cmd, etherTypeFlag),
		IcmpParameterCode:     flags.FlagToInt64Pointer(p, cmd, icmpParameterCodeFlag),
		IcmpParameterType:     flags.FlagToInt64Pointer(p, cmd, icmpParameterTypeFlag),
		IpRange:               flags.FlagToStringPointer(p, cmd, ipRangeFlag),
		PortRangeMax:          flags.FlagToInt64Pointer(p, cmd, portRangeMaxFlag),
		PortRangeMin:          flags.FlagToInt64Pointer(p, cmd, portRangeMinFlag),
		RemoteSecurityGroupId: flags.FlagToStringPointer(p, cmd, remoteSecurityGroupIdFlag),
		ProtocolNumber:        flags.FlagToInt64Pointer(p, cmd, protocolNumberFlag),
		ProtocolName:          flags.FlagToStringPointer(p, cmd, protocolNameFlag),
	}

	if model.ProtocolName != nil {
		if *model.ProtocolName == "icmp" || *model.ProtocolName == "ipv6-icmp" {
			if model.PortRangeMin != nil || model.PortRangeMax != nil {
				return nil, &cliErr.SecurityGroupRuleProtocolPortRangeConflictError{
					Cmd: cmd,
				}
			}
		} else {
			if model.IcmpParameterCode != nil || model.IcmpParameterType != nil {
				return nil, &cliErr.SecurityGroupRuleProtocolParametersConflictError{
					Cmd: cmd,
				}
			}
		}
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateSecurityGroupRuleRequest {
	req := apiClient.CreateSecurityGroupRule(ctx, model.ProjectId, model.SecurityGroupId)
	icmpParameters := &iaas.ICMPParameters{}
	portRange := &iaas.PortRange{}
	protocol := &iaas.CreateProtocol{}

	payload := iaas.CreateSecurityGroupRulePayload{
		Direction:             model.Direction,
		Description:           model.Description,
		Ethertype:             model.EtherType,
		IpRange:               model.IpRange,
		RemoteSecurityGroupId: model.RemoteSecurityGroupId,
	}

	if model.IcmpParameterCode != nil || model.IcmpParameterType != nil {
		icmpParameters.Code = model.IcmpParameterCode
		icmpParameters.Type = model.IcmpParameterType

		payload.IcmpParameters = icmpParameters
	}

	if model.PortRangeMax != nil || model.PortRangeMin != nil {
		portRange.Max = model.PortRangeMax
		portRange.Min = model.PortRangeMin

		payload.PortRange = portRange
	}

	if model.ProtocolNumber != nil || model.ProtocolName != nil {
		protocol.Int64 = model.ProtocolNumber
		protocol.String = model.ProtocolName

		payload.Protocol = protocol
	}

	if model.RemoteSecurityGroupId == nil {
		payload.RemoteSecurityGroupId = nil
	}

	return req.CreateSecurityGroupRulePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel, securityGroupName string, securityGroupRule *iaas.SecurityGroupRule) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(securityGroupRule, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal security group rule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(securityGroupRule, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal security group rule: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s security group rule for security group %q in project %q.\nSecurity group rule ID: %s\n", operationState, projectLabel, securityGroupName, *securityGroupRule.SecurityGroupId)
		return nil
	}
}
