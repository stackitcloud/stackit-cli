package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	sfsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

const (
	shareIdArg = "SHARE_ID"

	resourcePoolIdFlag = "resource-pool-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ResourcePoolId string
	ShareId        string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", shareIdArg),
		Short: "Shows details of a shares",
		Long:  "Shows details of a shares.",
		Args:  args.SingleArg(shareIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe a shares with ID "xxx" from resource pool with ID "yyy"`,
				"$ stackit beta sfs export-policy describe xxx --resource-pool-id yyy",
			),
		),
		RunE: func(cmd *cobra.Command, inputArgs []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, inputArgs)
			if err != nil {
				return fmt.Errorf("unable to parse input: %w", err)
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
				return fmt.Errorf("describe SFS share: %w", err)
			}

			resourcePoolLabel, err := sfsUtils.GetResourcePoolName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolLabel = model.ResourcePoolId
			} else if resourcePoolLabel == "" {
				resourcePoolLabel = model.ResourcePoolId
			}

			return outputResult(params.Printer, model.OutputFormat, resourcePoolLabel, model.ShareId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), resourcePoolIdFlag, "The resource pool the share is assigned to")

	err := flags.MarkFlagsRequired(cmd, resourcePoolIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	shareId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ResourcePoolId:  flags.FlagToStringValue(p, cmd, resourcePoolIdFlag),
		ShareId:         shareId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiGetShareRequest {
	return apiClient.GetShare(ctx, model.ProjectId, model.Region, model.ResourcePoolId, model.ShareId)
}

func outputResult(p *print.Printer, outputFormat, resourcePoolLabel, shareId string, share *sfs.GetShareResponse) error {
	return p.OutputResult(outputFormat, share, func() error {
		if share == nil || share.Share == nil {
			p.Outputf("Share %q not found in resource pool %q\n", shareId, resourcePoolLabel)
			return nil
		}

		var content []tables.Table

		table := tables.NewTable()
		table.SetTitle("Share")
		item := *share.Share

		table.AddRow("ID", utils.PtrString(item.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(item.Name))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(item.State))
		table.AddSeparator()
		table.AddRow("MOUNT PATH", utils.PtrString(item.MountPath))
		table.AddSeparator()
		table.AddRow("HARD LIMIT (GB)", utils.PtrString(item.SpaceHardLimitGigabytes))
		table.AddSeparator()
		table.AddRow("CREATED AT", utils.ConvertTimePToDateTimeString(item.CreatedAt))

		content = append(content, table)

		if item.HasExportPolicy() {
			policyTable := tables.NewTable()
			policyTable.SetTitle("Export Policy")

			policyTable.SetHeader(
				"ID",
				"NAME",
				"SHARES USING EXPORT POLICY",
				"CREATED AT",
			)

			policy := item.ExportPolicy.Get()

			policyTable.AddRow(
				utils.PtrString(policy.Id),
				utils.PtrString(policy.Name),
				utils.PtrString(policy.SharesUsingExportPolicy),
				utils.ConvertTimePToDateTimeString(policy.CreatedAt),
			)

			content = append(content, policyTable)

			if policy.Rules != nil && len(*policy.Rules) > 0 {
				ruleTable := tables.NewTable()
				ruleTable.SetTitle("Export Policy - Rules")

				ruleTable.SetHeader("ID", "ORDER", "DESCRIPTION", "IP ACL", "READ ONLY", "SET UUID", "SUPER USER", "CREATED AT")

				for _, rule := range *policy.Rules {
					var description string
					if rule.Description != nil {
						description = utils.PtrString(rule.Description.Get())
					}
					ruleTable.AddRow(
						utils.PtrString(rule.Id),
						utils.PtrString(rule.Order),
						description,
						utils.JoinStringPtr(rule.IpAcl, ", "),
						utils.PtrString(rule.ReadOnly),
						utils.PtrString(rule.SetUuid),
						utils.PtrString(rule.SuperUser),
						utils.ConvertTimePToDateTimeString(rule.CreatedAt),
					)
					ruleTable.AddSeparator()
				}

				content = append(content, ruleTable)
			}
		}

		if err := tables.DisplayTables(p, content); err != nil {
			return fmt.Errorf("render tables: %w", err)
		}
		return nil
	})
}
