package update

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	instanceNameFlag   = "name"
	aclFlag            = "acl"
	backupScheduleFlag = "backup-schedule"
	flavorIdFlag       = "flavor-id"
	storageSizeFlag    = "storage-size"
	versionFlag        = "version"
	retentionDaysFlag  = "retention-days"

	// Deprecated: Will be removed after 2027-01-31. Storage class can not be updated.
	storageClassFlag = "storage-class"

	cpuFlag = "cpu" // Deprecated: Will be removed after 2027-01-31. Flavor id should be used instead.
	ramFlag = "ram" // Deprecated: Will be removed after 2027-01-31. Flavor id should be used instead.
)

// Deprecated: Will be removed after 2027-01-31. Replicas are managed via the flavor id on API side now.
var typeFlag = flags.StringEnumFlag(
	"type",
	postgresflexUtils.AvailableInstanceTypes(),
	"Instance type,",
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId     string
	InstanceName   *string
	ACL            []string
	BackupSchedule *string
	FlavorId       *string
	StorageSize    *int64
	Version        *string
	RetentionDays  *int32

	CPU  *int64  // Deprecated: Will be removed after 2027-01-31
	RAM  *int64  // Deprecated: Will be removed after 2027-01-31
	Type *string // Deprecated: Will be removed after 2027-01-31
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a PostgreSQL Flex instance",
		Long:  "Updates a PostgreSQL Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`Update the name of a PostgreSQL Flex instance`,
				"$ stackit postgresflex instance update xxx --name my-new-name"),
			examples.NewExample(
				`Update the version of a PostgreSQL Flex instance`,
				"$ stackit postgresflex instance update xxx --version 6.0"),
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
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

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to update instance %q? (This may cause downtime)", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				return err
			}
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update PostgreSQL Flex instance: %w", err)
			}

			// update endpoint doesn't return the updated instance, so we have to call the GET endpoint to fetch it
			var getResp *postgresflex.GetInstanceResponse

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating instance", func() error {
					getResp, err = wait.PartialUpdateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for PostgreSQL Flex instance update: %w", err)
				}
			} else {
				getResp, err = apiClient.DefaultAPI.GetInstance(ctx, model.ProjectId, model.Region, model.InstanceId).Execute()
				if err != nil {
					return fmt.Errorf("fetching PostgreSQL Flex instance after async update: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, instanceLabel, getResp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")
	cmd.Flags().String(backupScheduleFlag, "", "Backup schedule")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor")
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs")            // remove after 2027-01-31
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB)")     // remove after 2027-01-31
	cmd.Flags().String(storageClassFlag, "", "Storage class.") // remove after 2027-01-31
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB)")
	cmd.Flags().String(versionFlag, "", "Version")
	cmd.Flags().String(retentionDaysFlag, "", "The days for how long the backup files should be stored before cleaned up (32 to 90).")
	typeFlag.Register(cmd.Flags())

	// remove after 2027-01-31
	err := cmd.Flags().MarkDeprecated(storageClassFlag, "This flag has no effect and will be removed after 2027-01-31.")
	cobra.CheckErr(err)

	// remove after 2027-01-31
	err = cmd.Flags().MarkDeprecated(typeFlag.Name(), fmt.Sprintf("Will be removed after 2027-01-31. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	err = cmd.Flags().MarkDeprecated(cpuFlag, fmt.Sprintf("Will be removed after 2027-01-31. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	err = cmd.Flags().MarkDeprecated(ramFlag, fmt.Sprintf("Will be removed after 2027-01-31. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	cmd.MarkFlagsMutuallyExclusive(flavorIdFlag, cpuFlag)
	cmd.MarkFlagsMutuallyExclusive(flavorIdFlag, ramFlag)
	cmd.MarkFlagsMutuallyExclusive(flavorIdFlag, "type")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	instanceName := flags.FlagToStringPointer(p, cmd, instanceNameFlag)
	flavorId := flags.FlagToStringPointer(p, cmd, flavorIdFlag)
	cpu := flags.FlagToInt64Pointer(p, cmd, cpuFlag)
	ram := flags.FlagToInt64Pointer(p, cmd, ramFlag)
	acl := flags.FlagToStringSliceValue(p, cmd, aclFlag)
	backupSchedule := flags.FlagToStringPointer(p, cmd, backupScheduleFlag)
	storageSize := flags.FlagToInt64Pointer(p, cmd, storageSizeFlag)
	version := flags.FlagToStringPointer(p, cmd, versionFlag)
	instanceType := typeFlag.Ptr()
	retentionDays := flags.FlagToInt32Pointer(p, cmd, retentionDaysFlag)

	if instanceName == nil && flavorId == nil && cpu == nil && ram == nil && acl == nil &&
		backupSchedule == nil && storageSize == nil && version == nil && instanceType == nil && retentionDays == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	if flavorId != nil && (cpu != nil || ram != nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd:  cmd,
			Args: inputArgs,
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		InstanceName:    instanceName,
		ACL:             acl,
		BackupSchedule:  backupSchedule,
		FlavorId:        flavorId,
		StorageSize:     storageSize,
		Version:         version,
		RetentionDays:   retentionDays,

		// deprecated fields
		Type: instanceType,
		CPU:  cpu,
		RAM:  ram,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient postgresflex.DefaultAPI) (postgresflex.ApiPartialUpdateInstanceRequest, error) {
	var flavorId *string
	var err error

	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.Region, model.InstanceId)

	currentInstance, err := apiClient.GetInstance(ctx, model.ProjectId, model.Region, model.InstanceId).Execute()
	if err != nil {
		return req, fmt.Errorf("failed to get instance %s: %w", model.InstanceId, err)
	}

	flavors, err := apiClient.ListFlavors(ctx, model.ProjectId, model.Region).Execute()
	if err != nil {
		return req, fmt.Errorf("get PostgreSQL Flex flavors: %w", err)
	}

	// if cpu/ram flags are used instead of the flavor id flag
	if model.FlavorId == nil && (model.RAM != nil || model.CPU != nil) {
		ram := model.RAM
		cpu := model.CPU

		// if only one of the cpu/ram flags is set
		if model.RAM == nil || model.CPU == nil {
			var currentFlavor *postgresflex.ListFlavors
			for _, f := range flavors.Flavors {
				if f.Id == currentInstance.FlavorId {
					currentFlavor = &f
				}
			}

			if currentFlavor == nil {
				return req, fmt.Errorf("flavor %s not found", currentInstance.FlavorId)
			}

			if model.RAM == nil {
				ram = utils.Ptr(currentFlavor.Memory)
			}
			if model.CPU == nil {
				cpu = utils.Ptr(currentFlavor.Cpu)
			}
		}

		flavorId, err = postgresflexUtils.LoadFlavorId(*cpu, *ram, flavors.Flavors)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load flavor ID: %w", err)
			}
			return req, err
		}
	} else if model.FlavorId != nil {
		flavorId = model.FlavorId
	}

	var payloadNetwork *postgresflex.InstanceNetworkOpt
	if model.ACL != nil {
		payloadNetwork = &postgresflex.InstanceNetworkOpt{
			Acl: model.ACL,
		}
	}

	var payloadStorage *postgresflex.StorageUpdate
	if model.StorageSize != nil {
		payloadStorage = &postgresflex.StorageUpdate{
			Size: model.StorageSize,
		}
	}

	payload := postgresflex.PartialUpdateInstancePayload{
		BackupSchedule: model.BackupSchedule,
		FlavorId:       flavorId,
		Name:           model.InstanceName,
		Network:        payloadNetwork,
		RetentionDays:  model.RetentionDays,
		Storage:        payloadStorage,
		Version:        model.Version,
	}

	return req.PartialUpdateInstancePayload(payload), nil
}

func outputResult(p *print.Printer, outputFormat string, async bool, instanceLabel string, resp *postgresflex.GetInstanceResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("no response passed")
		}

		operationState := "Updated"
		if async {
			operationState = "Triggered update of"
		}

		p.Outputf("%s instance %q\n", operationState, instanceLabel)
		return nil
	})
}
