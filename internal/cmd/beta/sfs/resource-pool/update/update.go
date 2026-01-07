package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/client"
	sfsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sfs/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs/wait"
)

const (
	resourcePoolIdArg    = "RESOURCE_POOL_ID"
	performanceClassFlag = "performance-class"
	sizeFlag             = "size"
	ipAclFlag            = "ip-acl"
	snapshotsVisibleFlag = "snapshots-visible"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SizeGigabytes    *int64
	PerformanceClass *string
	IpAcl            *[]string
	ResourcePoolId   string
	SnapshotsVisible *bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a SFS resource pool",
		Long: `Updates a SFS resource pool.

The available performance class values can be obtained by running:
 $ stackit beta sfs performance-class list`,
		Args: args.SingleArg(resourcePoolIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the SFS resource pool with ID "xxx"`,
				"$ stackit beta sfs resource-pool update xxx --ip-acl 10.88.135.144/28 --performance-class Standard --size 5"),
			examples.NewExample(
				`Update the SFS resource pool with ID "xxx", allow only a single IP which can mount the resource pool`,
				"$ stackit beta sfs resource-pool update xxx --ip-acl 250.81.87.224/32 --performance-class Standard --size 5"),
			examples.NewExample(
				`Update the SFS resource pool with ID "xxx", allow multiple IP ACL which can mount the resource pool`,
				"$ stackit beta sfs resource-pool update xxx --ip-acl \"10.88.135.144/28,250.81.87.224/32\" --performance-class Standard --size 5"),
			examples.NewExample(
				`Update the SFS resource pool with ID "xxx", set snapshots visible to false`,
				"$ stackit beta sfs resource-pool update xxx --snapshots-visible=false"),
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
			}

			resourcePoolName, err := sfsUtils.GetResourcePoolName(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get resource pool name: %v", err)
				resourcePoolName = model.ResourcePoolId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update resource-pool %q for project %q?", resourcePoolName, projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("update SFS resource pool: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Update resource pool")
				_, err = wait.UpdateResourcePoolWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.ResourcePoolId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for resource pool update: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(sizeFlag, 0, "Size of the pool in Gigabytes")
	cmd.Flags().String(performanceClassFlag, "", "Performance class")
	cmd.Flags().Var(flags.CIDRSliceFlag(), ipAclFlag, "List of network addresses in the form <address/prefix>, e.g. 192.168.10.0/24 that can mount the resource pool readonly")
	cmd.Flags().Bool(snapshotsVisibleFlag, false, "Set snapshots visible and accessible to users")
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiUpdateResourcePoolRequest {
	req := apiClient.UpdateResourcePool(ctx, model.ProjectId, model.Region, model.ResourcePoolId)
	req = req.UpdateResourcePoolPayload(sfs.UpdateResourcePoolPayload{
		IpAcl:               model.IpAcl,
		PerformanceClass:    model.PerformanceClass,
		SizeGigabytes:       model.SizeGigabytes,
		SnapshotsAreVisible: model.SnapshotsVisible,
	})
	return req
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	resourcePoolId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	performanceClass := flags.FlagToStringPointer(p, cmd, performanceClassFlag)
	size := flags.FlagToInt64Pointer(p, cmd, sizeFlag)
	ipAcls := flags.FlagToStringSlicePointer(p, cmd, ipAclFlag)
	snapshotsVisible := flags.FlagToBoolPointer(p, cmd, snapshotsVisibleFlag)

	if performanceClass == nil && size == nil && ipAcls == nil && snapshotsVisible == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		SizeGigabytes:    size,
		IpAcl:            ipAcls,
		PerformanceClass: performanceClass,
		ResourcePoolId:   resourcePoolId,
		SnapshotsVisible: snapshotsVisible,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat string, resp *sfs.UpdateResourcePoolResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil || resp.ResourcePool == nil {
			p.Outputln("Resource pool response is empty")
			return nil
		}
		p.Outputf("Updated resource pool %s\n", utils.PtrString(resp.ResourcePool.Name))
		return nil
	})
}
