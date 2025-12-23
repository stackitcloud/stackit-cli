package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
)

const (
	resourcePoolIdArg = "RESOURCE_POOL_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ResourcePoolId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of a SFS resource pool",
		Long:  "Shows details of a SFS resource pool.",
		Args:  args.SingleArg(resourcePoolIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Describe the SFS resource pool with ID "xxx"`,
				"$ stackit beta sfs resource-pool describe xxx"),
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
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("describe SFS resource pool: %w", err)
			}

			// Get projectLabel
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			return outputResult(params.Printer, model.OutputFormat, model.ResourcePoolId, projectLabel, resp.ResourcePool)
		},
	}
	return cmd
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiGetResourcePoolRequest {
	req := apiClient.GetResourcePool(ctx, model.ProjectId, model.Region, model.ResourcePoolId)
	return req
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	resourcePoolId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ResourcePoolId:  resourcePoolId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, resourcePoolId, projectLabel string, resourcePool *sfs.GetResourcePoolResponseResourcePool) error {
	return p.OutputResult(outputFormat, resourcePool, func() error {
		if resourcePool == nil {
			p.Outputf("Resource pool %q not found in project %q\n", resourcePoolId, projectLabel)
			return nil
		}
		table := tables.NewTable()

		// convert the string slice to a comma separated list
		var ipAclStr string
		if resourcePool.IpAcl != nil {
			ipAclStr = strings.Join(*resourcePool.IpAcl, ", ")
		}

		table.AddRow("ID", utils.PtrString(resourcePool.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(resourcePool.Name))
		table.AddSeparator()
		table.AddRow("AVAILABILITY ZONE", utils.PtrString(resourcePool.AvailabilityZone))
		table.AddSeparator()
		table.AddRow("NUMBER OF SHARES", utils.PtrString(resourcePool.CountShares))
		table.AddSeparator()
		table.AddRow("IP ACL", ipAclStr)
		table.AddSeparator()
		table.AddRow("MOUNT PATH", utils.PtrString(resourcePool.MountPath))
		table.AddSeparator()
		if resourcePool.PerformanceClass != nil {
			table.AddRow("PERFORMANCE CLASS", utils.PtrString(resourcePool.PerformanceClass.Name))
			table.AddSeparator()
		}
		table.AddRow("SNAPSHOTS ARE VISIBLE", utils.PtrString(resourcePool.SnapshotsAreVisible))
		table.AddSeparator()
		table.AddRow("NEXT PERFORMANCE CLASS DOWNGRADE TIME", utils.PtrString(resourcePool.PerformanceClassDowngradableAt))
		table.AddSeparator()
		table.AddRow("NEXT SIZE REDUCTION TIME", utils.PtrString(resourcePool.SizeReducibleAt))
		table.AddSeparator()
		if resourcePool.HasSpace() {
			table.AddRow("TOTAL SIZE (GB)", utils.PtrString(resourcePool.Space.SizeGigabytes))
			table.AddSeparator()
			table.AddRow("AVAILABLE SIZE (GB)", utils.PtrString(resourcePool.Space.AvailableGigabytes))
			table.AddSeparator()
			table.AddRow("USED SIZE (GB)", utils.PtrString(resourcePool.Space.UsedGigabytes))
			table.AddSeparator()
		}
		table.AddRow("STATE", utils.PtrString(resourcePool.State))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
