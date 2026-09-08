package create

import (
	"context"
	"fmt"

	ufw "github.com/stackitcloud/stackit-sdk-go/services/ufw/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/ufw/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ufw/client"

	"github.com/spf13/cobra"
)

const (
	productFlag         = "product"
	typeFlag            = "type"
	sourceIpFlag        = "sourceIp"
	instanceIdFlag      = "instanceId"
	directionFlag       = "direction"
	descriptionFlag     = "description"
	etherTypeFlag       = "etherType"
	portRangeFlag       = "portRange"
	protocolFlag        = "protocol"
	offsetFlag          = "offset"
	securityGroupIdFlag = "securityGroupId"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Product         *string
	Type            *string
	SourceIp        *string
	InstanceId      *string
	Direction       *string
	Description     *string
	EtherType       *string
	PortRange       *string
	Protocol        *string
	Offset          *int32
	SecurityGroupId *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an UFW rule instance",
		Long:  "Creates a STACKIT Unified Firewall (UFW) rule instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a UFW rule instance of type ACL with sourceIp "1.1.1.1/32" of product "redis" for instance with id=ID`,
				"$ stackit ufw instance create --product redis --sourceIp 1.1.1.1/32 --type ACL --instanceId ID"),
			// TODO add more examples for creating Security Rule and Group types
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a UFW rule instance for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create UFW rule instance: %w", err)
			}
			instanceId := resp.RefId

			if !model.Async {
				err := spinner.Run(params.Printer, "Creating ufw rule instance", func() error {
					_, err = wait.CreateRuleWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, *instanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for UFW rule instance creation: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(productFlag, "", "The source service (e.g., Load Balancer, Redis) where you want to attach a rule")
	cmd.Flags().StringP(typeFlag, "t", "", "Type (ACL/SecurityRule/SecurityGroup/PublicIP) You can check /provider-options route for them")
	cmd.Flags().StringP(sourceIpFlag, "s", "", "The IP (CIDR) to which the rule applies (e.g. 192.168.0.1/32)")
	cmd.Flags().StringP(instanceIdFlag, "i", "", "Instance ID that will have attached your rule")
	cmd.Flags().StringP(directionFlag, "d", "", "Direction (the direction of the traffic, typically ingress or egress, for security rules type)")
	cmd.Flags().StringP(descriptionFlag, "D", "", "Description")
	cmd.Flags().StringP(etherTypeFlag, "e", "", "Specifies the bound of the rule (for security rules type)")
	cmd.Flags().StringP(portRangeFlag, "r", "", "Port range (the Port range to which the rule applies, for security rules type)")
	cmd.Flags().String(protocolFlag, "", "The network protocol (e.g. TCP, UDP, ICMP, for security rules type)")
	cmd.Flags().Int32P(offsetFlag, "f", -1, "Offset - Position in the ACL list of an instance, will be ignored at creation")
	cmd.Flags().StringP(securityGroupIdFlag, "g", "", "Security group ID - The ID of the Security Group")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)

	err = flags.MarkFlagsRequired(cmd, productFlag)
	cobra.CheckErr(err)

	err = flags.MarkFlagsRequired(cmd, sourceIpFlag)
	cobra.CheckErr(err)

	err = flags.MarkFlagsRequired(cmd, typeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	if globalFlags.Region == "" {
		return nil, &errors.RegionError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,

		Product:         flags.FlagToStringPointer(p, cmd, productFlag),
		Type:            flags.FlagToStringPointer(p, cmd, typeFlag),
		SourceIp:        flags.FlagToStringPointer(p, cmd, sourceIpFlag),
		InstanceId:      flags.FlagToStringPointer(p, cmd, instanceIdFlag),
		Direction:       flags.FlagToStringPointer(p, cmd, directionFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		EtherType:       flags.FlagToStringPointer(p, cmd, etherTypeFlag),
		PortRange:       flags.FlagToStringPointer(p, cmd, portRangeFlag),
		Protocol:        flags.FlagToStringPointer(p, cmd, protocolFlag),
		Offset:          flags.FlagToInt32Pointer(p, cmd, offsetFlag),
		SecurityGroupId: flags.FlagToStringPointer(p, cmd, securityGroupIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ufw.APIClient) ufw.ApiCreateRuleRequest {
	req := apiClient.DefaultAPI.CreateRule(ctx, model.ProjectId, model.Region)

	// TODO - add logic for field checking: existing ACLs, correct product, type, instanceID maybe

	req = req.CreateRulePayload(ufw.CreateRulePayload{
		Product:         *model.Product,
		Type:            *model.Type,
		SourceIP:        *model.SourceIp,
		InstanceId:      *model.InstanceId,
		Direction:       model.Direction,
		Description:     model.Description,
		EtherType:       model.EtherType,
		PortRange:       model.PortRange,
		Protocol:        model.Protocol,
		Offset:          model.Offset,
		SecurityGroupId: model.SecurityGroupId,
	})

	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, rule *ufw.CreateRuleResponse) error {
	if rule == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(outputFormat, rule, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s rule for project %q. Rule refID: %s\n", operationState, projectLabel, utils.PtrString(rule.RefId))
		return nil
	})
}
