package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	ufw "github.com/stackitcloud/stackit-sdk-go/services/ufw/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/ufw/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ufw/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const (
	instanceIdArg = "INSTANCE_ID"

	sourceIpFlag = "sourceIp"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string

	SourceIp    *string
	Direction   *string
	Description *string
	EtherType   *string
	PortRange   *string
	Protocol    *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a UFW rule instance",
		Long:  "Updates a STACKIT Unified Firewall (UFW) rule instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update a UFW rule instance with "1.1.1.1/32" as sourceIp for instance with ID "xxx"`,
				"$ stackit ufw rules update xxx --sourceIp 1.1.1.1/32"),
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

			prompt := fmt.Sprintf("Are you sure you want to update a UFW rule instance for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update UFW rule instance: %w", err)
			}
			instanceId := resp.RefId

			if !model.Async {
				err := spinner.Run(params.Printer, "Updating ufw rule instance", func() error {
					_, err = wait.UpdateRuleWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, *instanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for UFW rule instance updating process: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(sourceIpFlag, "s", "", "The IP (CIDR) to which the rule applies (e.g. 192.168.0.1/32)")

	err := flags.MarkFlagsRequired(cmd, sourceIpFlag)
	cobra.CheckErr(err)
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

		SourceIp: flags.FlagToStringPointer(p, cmd, sourceIpFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ufw.APIClient) ufw.ApiUpdateRuleRequest {
	req := apiClient.DefaultAPI.UpdateRule(ctx, model.ProjectId, model.Region, model.InstanceId)

	req = req.UpdateRulePayload(ufw.UpdateRulePayload{
		SourceIP: *model.SourceIp,
	})

	return req
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, rule *ufw.UpdateRuleResponse) error {
	if rule == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(outputFormat, rule, func() error {
		operationState := "Updated"
		if async {
			operationState = "Triggered updating process of"
		}
		p.Outputf("%s rule for project %q. New rule refID: %s\n", operationState, projectLabel, utils.PtrString(rule.RefId))
		return nil
	})
}
