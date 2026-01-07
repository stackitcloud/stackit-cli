package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs"
	"github.com/stackitcloud/stackit-sdk-go/services/sfs/wait"
)

const (
	performanceClassFlag = "performance-class"
	sizeFlag             = "size"
	ipAclFlag            = "ip-acl"
	availabilityZoneFlag = "availability-zone"
	nameFlag             = "name"
	snapshotsVisibleFlag = "snapshots-visible"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SizeInGB         int64
	PerformanceClass string
	IpAcl            []string
	Name             string
	AvailabilityZone string
	SnapshotsVisible bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a SFS resource pool",
		Long: `Creates a SFS resource pool.

The available performance class values can be obtained by running:
 $ stackit beta sfs performance-class list`,
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a SFS resource pool`,
				"$ stackit beta sfs resource-pool create --availability-zone eu01-m --ip-acl 10.88.135.144/28 --performance-class Standard --size 500 --name resource-pool-01"),
			examples.NewExample(
				`Create a SFS resource pool, allow only a single IP which can mount the resource pool`,
				"$ stackit beta sfs resource-pool create --availability-zone eu01-m --ip-acl 250.81.87.224/32 --performance-class Standard --size 500 --name resource-pool-01"),
			examples.NewExample(
				`Create a SFS resource pool, allow multiple IP ACL which can mount the resource pool`,
				"$ stackit beta sfs resource-pool create --availability-zone eu01-m --ip-acl \"10.88.135.144/28,250.81.87.224/32\" --performance-class Standard --size 500 --name resource-pool-01"),
			examples.NewExample(
				`Create a SFS resource pool with visible snapshots`,
				"$ stackit beta sfs resource-pool create --availability-zone eu01-m --ip-acl 10.88.135.144/28 --performance-class Standard --size 500 --name resource-pool-01 --snapshots-visible"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a resource-pool for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			resp, err := buildRequest(ctx, model, apiClient).Execute()
			if err != nil {
				return fmt.Errorf("create SFS resource pool: %w", err)
			}
			var resourcePoolId string
			if resp != nil && resp.HasResourcePool() && resp.ResourcePool.HasId() {
				resourcePoolId = *resp.ResourcePool.Id
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Create resource pool")
				_, err = wait.CreateResourcePoolWaitHandler(ctx, apiClient, model.ProjectId, model.Region, resourcePoolId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for resource pool creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(sizeFlag, 0, "Size of the pool in Gigabytes")
	cmd.Flags().String(performanceClassFlag, "", "Performance class")
	cmd.Flags().Var(flags.CIDRSliceFlag(), ipAclFlag, "List of network addresses in the form <address/prefix>, e.g. 192.168.10.0/24 that can mount the resource pool readonly")
	cmd.Flags().String(availabilityZoneFlag, "", "Availability zone")
	cmd.Flags().String(nameFlag, "", "Name")
	cmd.Flags().Bool(snapshotsVisibleFlag, false, "Set snapshots visible and accessible to users")

	for _, flag := range []string{sizeFlag, performanceClassFlag, ipAclFlag, availabilityZoneFlag, nameFlag} {
		err := flags.MarkFlagsRequired(cmd, flag)
		cobra.CheckErr(err)
	}
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *sfs.APIClient) sfs.ApiCreateResourcePoolRequest {
	req := apiClient.CreateResourcePool(ctx, model.ProjectId, model.Region)
	req = req.CreateResourcePoolPayload(sfs.CreateResourcePoolPayload{
		AvailabilityZone:    &model.AvailabilityZone,
		IpAcl:               &model.IpAcl,
		Name:                &model.Name,
		PerformanceClass:    &model.PerformanceClass,
		SizeGigabytes:       &model.SizeInGB,
		SnapshotsAreVisible: &model.SnapshotsVisible,
	})
	return req
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	performanceClass := flags.FlagToStringValue(p, cmd, performanceClassFlag)
	size := flags.FlagWithDefaultToInt64Value(p, cmd, sizeFlag)
	availabilityZone := flags.FlagToStringValue(p, cmd, availabilityZoneFlag)
	ipAcls := flags.FlagToStringSlicePointer(p, cmd, ipAclFlag)
	name := flags.FlagToStringValue(p, cmd, nameFlag)
	snapshotsVisible := flags.FlagToBoolValue(p, cmd, snapshotsVisibleFlag)

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		SizeInGB:         size,
		IpAcl:            *ipAcls,
		PerformanceClass: performanceClass,
		AvailabilityZone: availabilityZone,
		Name:             name,
		SnapshotsVisible: snapshotsVisible,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *sfs.CreateResourcePoolResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil || resp.ResourcePool == nil {
			p.Outputln("Resource pool response is empty")
			return nil
		}
		p.Outputf("Created resource pool for project %q. Resource pool ID: %s\n", projectLabel, utils.PtrString(resp.ResourcePool.Id))
		return nil
	})
}
