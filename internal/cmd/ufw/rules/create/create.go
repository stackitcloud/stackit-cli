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
	productFlag    = "product"
	typeFlag       = "type"
	sourceIpFlag   = "sourceIp"
	instanceIdFlag = "instanceId"
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
		Short: "Creates a UFW rule instance",
		Long:  "Creates a STACKIT Unified Firewall (UFW) rule instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a UFW rule instance of type ACL with sourceIp "1.1.1.1/32" of product "Redis" for instance with id=ID`,
				"$ stackit ufw rules create --product redis --sourceIp 1.1.1.1/32 --type ACL --instanceId ID"),
			examples.NewExample(
				`Create a UFW rule instance of type ACL with sourceIp "2.2.2.2/32" of product "Edge Cloud" for instance with id=ID`,
				"$ stackit ufw rules create --product edge-cloud --sourceIp 2.2.2.2/32 --type ACL --instanceId ID"),
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

			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}

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
	cmd.Flags().String(productFlag, "", "The source service (e.g., Edge Cloud, Redis) where you want to attach a rule")
	cmd.Flags().StringP(typeFlag, "t", "", "Type (ACL/SecurityRule/SecurityGroup) You can check /provider-options route for them. Unfortunately, this field could be only ACL for the CLI version")
	cmd.Flags().StringP(sourceIpFlag, "s", "", "The IP (CIDR) to which the rule applies (e.g. 192.168.0.1/32)")
	cmd.Flags().StringP(instanceIdFlag, "i", "", "Instance ID that will have attached your rule")

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

		Product:    flags.FlagToStringPointer(p, cmd, productFlag),
		Type:       flags.FlagToStringPointer(p, cmd, typeFlag),
		SourceIp:   flags.FlagToStringPointer(p, cmd, sourceIpFlag),
		InstanceId: flags.FlagToStringPointer(p, cmd, instanceIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ufw.APIClient) (ufw.ApiCreateRuleRequest, error) {
	req := apiClient.DefaultAPI.CreateRule(ctx, model.ProjectId, model.Region)

	if *model.Type != "ACL" {
		return req, fmt.Errorf("invalid rule type: %s", *model.Type)
	}

	req = req.CreateRulePayload(ufw.CreateRulePayload{
		Product:    *model.Product,
		Type:       *model.Type,
		SourceIP:   *model.SourceIp,
		InstanceId: *model.InstanceId,
	})

	return req, nil
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
