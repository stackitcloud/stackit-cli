package list

import (
	"context"
	"fmt"
	"strconv"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists quotas",
		Long:  "Lists project quotas.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List available quotas`,
				`$ stackit quota list`,
			),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			response, err := request.Execute()
			if err != nil {
				return fmt.Errorf("list quotas: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, response.Quotas)
		},
	}

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListQuotasRequest {
	request := apiClient.DefaultAPI.ListQuotas(ctx, model.ProjectId, model.Region)

	return request
}

func outputResult(p *print.Printer, outputFormat string, quotas iaas.QuotaList) error {
	return p.OutputResult(outputFormat, quotas, func() error {
		table := tables.NewTable()
		table.SetHeader("NAME", "LIMIT", "CURRENT USAGE", "PERCENT")
		table.AddRow(quotaRow("Total size in GiB of backups [GiB]", quotas.BackupGigabytes)...)
		table.AddRow(quotaRow("Number of backups [Count]", quotas.Backups)...)
		table.AddRow(quotaRow("Total size in GiB of volumes and snapshots [GiB]", quotas.Gigabytes)...)
		table.AddRow(quotaRow("Number of networks [Count]", quotas.Networks)...)
		table.AddRow(quotaRow("Number of network interfaces (nics) [Count]", quotas.Nics)...)
		table.AddRow(quotaRow("Number of public IP addresses [Count]", quotas.PublicIps)...)
		table.AddRow(quotaRow("Amount of server RAM in MiB [MiB]", quotas.Ram)...)
		table.AddRow(quotaRow("Number of security group rules [Count]", quotas.SecurityGroupRules)...)
		table.AddRow(quotaRow("Number of security groups [Count]", quotas.SecurityGroups)...)
		table.AddRow(quotaRow("Number of snapshots [Count]", quotas.Snapshots)...)
		table.AddRow(quotaRow("Number of server cores (vcpu) [Count]", quotas.Vcpu)...)
		table.AddRow(quotaRow("Number of volumes [Count]", quotas.Volumes)...)
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}

func quotaRow(description string, quota iaas.Quota) []interface{} {
	result := make([]interface{}, 0, 4)
	result = append(result, description)
	result = append(result, conv(quota.Limit))
	result = append(result, conv(quota.Usage))
	result = append(result, fmt.Sprintf("%3.1f%%", 100.0/float64(quota.Limit)*float64(quota.Usage)))
	return result
}

func conv(n int64) string {
	return strconv.FormatInt(n, 10)
}

func percentage(val interface {
	GetLimitOk() (int64, bool)
	GetUsageOk() (int64, bool)
}) string {
	a, aOk := val.GetLimitOk()
	b, bOk := val.GetUsageOk()
	if aOk && bOk {
		return fmt.Sprintf("%3.1f%%", 100.0/float64(a)*float64(b))
	}
	return "n/a"
}
