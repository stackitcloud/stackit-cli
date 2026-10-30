package describe

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	platformIdArg = "PLATFORM_ID"

	quotaIdFlag = "quota-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PlatformId string
	QuotaId    *string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", platformIdArg),
		Short: "Shows details of a STACKIT Cloud Foundry platform quota",
		Long:  "Shows details of a STACKIT Cloud Foundry platform quota.",
		Args:  args.SingleArg(platformIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of the quota with ID "yyy" of a STACKIT Cloud Foundry platform with ID "xxx"`,
				`$ stackit beta scf platform quota describe xxx --quota-id yyy`),
			examples.NewExample(
				`Get details of the quota with ID "yyy" of a STACKIT Cloud Foundry platform with ID "xxx" in JSON format`,
				`$ stackit beta scf platform quota describe xxx --quota-id yyy --output-format json`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get STACKIT Cloud Foundry platform quota: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), quotaIdFlag, "The UUID of the quota")

	err := flags.MarkFlagsRequired(cmd, quotaIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	platformId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		PlatformId:      platformId,
		QuotaId:         flags.FlagToStringPointer(p, cmd, quotaIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiGetPlatformQuotaRequest {
	return apiClient.GetPlatformQuota(ctx, model.ProjectId, model.Region, model.PlatformId, *model.QuotaId)
}

// nullableInt64ToString renders a nullable quota limit. A null value means "unlimited".
func nullableInt64ToString(v scf.NullableInt64) string {
	if value := v.Get(); value != nil {
		return fmt.Sprintf("%d", *value)
	}
	return "unlimited"
}

func outputResult(p *print.Printer, outputFormat string, quota *scf.Quota) error {
	return p.OutputResult(outputFormat, quota, func() error {
		if quota == nil {
			return fmt.Errorf("received nil response, could not display quota details")
		}

		apps := quota.GetApps()
		domains := quota.GetDomains()
		routes := quota.GetRoutes()
		services := quota.GetServices()

		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")

		table.AddRow("ID", quota.GetGuid())
		table.AddRow("Name", quota.GetName())
		table.AddRow("Organization ID", utils.PtrString(quota.OrgId))
		table.AddRow("Platform ID", quota.GetPlatformId())
		table.AddRow("Project ID", quota.GetProjectId())
		table.AddRow("Region", quota.GetRegion())
		table.AddRow("Total App Instances", nullableInt64ToString(apps.TotalInstances))
		table.AddRow("Total App Memory (MB)", nullableInt64ToString(apps.TotalMemoryInMb))
		table.AddRow("Per-Process Memory (MB)", nullableInt64ToString(apps.PerProcessMemoryInMb))
		table.AddRow("Per-App Tasks", nullableInt64ToString(apps.PerAppTasks))
		table.AddRow("Log Rate Limit (bytes/s)", nullableInt64ToString(apps.LogRateLimitInBytesPerSecond))
		table.AddRow("Total Domains", nullableInt64ToString(domains.TotalDomains))
		table.AddRow("Total Routes", nullableInt64ToString(routes.TotalRoutes))
		table.AddRow("Total Reserved Ports", nullableInt64ToString(routes.TotalReservedPorts))
		table.AddRow("Paid Services Allowed", services.GetPaidServicesAllowed())
		table.AddRow("Total Service Instances", nullableInt64ToString(services.TotalServiceInstances))
		table.AddRow("Total Service Keys", nullableInt64ToString(services.TotalServiceKeys))
		table.AddRow("Created", quota.CreatedAt.Format(time.DateTime))
		table.AddRow("Updated", quota.UpdatedAt.Format(time.DateTime))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
