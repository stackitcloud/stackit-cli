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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	mongodbflex "github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/v2api"
	wait "github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/v2api/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	instanceNameFlag   = "name"
	aclFlag            = "acl"
	backupScheduleFlag = "backup-schedule"
	flavorIdFlag       = "flavor-id"
	storageClassFlag   = "storage-class"
	storageSizeFlag    = "storage-size"
	versionFlag        = "version"

	cpuFlag = "cpu" // Deprecated: Will be removed after 2027-03-07. Flavor id should be used instead.
	ramFlag = "ram" // Deprecated: Will be removed after 2027-03-07. Flavor id should be used instead.
)

var typeFlag = flags.StringEnumFlag(
	"type",
	mongodbflexUtils.AvailableInstanceTypes(),
	"Instance type,",
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId     string
	InstanceName   *string
	ACL            *[]string
	BackupSchedule *string
	FlavorId       *string
	StorageClass   *string
	StorageSize    *int64
	Version        *string
	Type           *string
	CPU            *int32 // Deprecated: Will be removed after 2027-03-07.
	RAM            *int32 // Deprecated: Will be removed after 2027-03-07.
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a MongoDB Flex instance",
		Long:  "Updates a MongoDB Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`Update the name of a MongoDB Flex instance`,
				"$ stackit mongodbflex instance update xxx --name my-new-name"),
			examples.NewExample(
				`Update the version of a MongoDB Flex instance`,
				"$ stackit mongodbflex instance update xxx --version 8.0"),
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

			instanceLabel, err := mongodbflexUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.InstanceId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to update instance %q?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update MongoDB Flex instance: %w", err)
			}
			instanceId := *resp.Item.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating instance", func() error {
					_, err = wait.PartialUpdateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, instanceId, model.Region).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for MongoDB Flex instance update: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "Lists of IP networks in CIDR notation which are allowed to access this instance")
	cmd.Flags().String(backupScheduleFlag, "", "Backup schedule")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor")
	cmd.Flags().String(storageClassFlag, "", "Storage class")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB)")
	cmd.Flags().String(versionFlag, "", "Version")
	cmd.Flags().Int32(cpuFlag, 0, "Number of CPUs")        // Deprecated: Will be removed after 2027-03-07.
	cmd.Flags().Int32(ramFlag, 0, "Amount of RAM (in GB)") // Deprecated: Will be removed after 2027-03-07.
	typeFlag.Register(cmd.Flags())

	// Deprecated: Will be removed after 2027-03-07.
	err := cmd.Flags().MarkDeprecated(cpuFlag, fmt.Sprintf("Will be removed after 2027-03-07. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	err = cmd.Flags().MarkDeprecated(ramFlag, fmt.Sprintf("Will be removed after 2027-03-07. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	cmd.MarkFlagsRequiredTogether(cpuFlag, ramFlag)
	cmd.MarkFlagsMutuallyExclusive(flavorIdFlag, cpuFlag)
	cmd.MarkFlagsMutuallyExclusive(flavorIdFlag, ramFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	instanceName := flags.FlagToStringPointer(p, cmd, instanceNameFlag)
	flavorId := flags.FlagToStringPointer(p, cmd, flavorIdFlag)
	cpu := flags.FlagToInt32Pointer(p, cmd, cpuFlag)
	ram := flags.FlagToInt32Pointer(p, cmd, ramFlag)
	acl := flags.FlagToStringSlicePointer(p, cmd, aclFlag)
	backupSchedule := flags.FlagToStringPointer(p, cmd, backupScheduleFlag)
	storageClass := flags.FlagToStringPointer(p, cmd, storageClassFlag)
	storageSize := flags.FlagToInt64Pointer(p, cmd, storageSizeFlag)
	version := flags.FlagToStringPointer(p, cmd, versionFlag)
	instanceType := typeFlag.Ptr()

	if instanceName == nil && flavorId == nil && cpu == nil && ram == nil && acl == nil &&
		backupSchedule == nil && storageClass == nil && storageSize == nil && version == nil && instanceType == nil {
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
		StorageClass:    storageClass,
		StorageSize:     storageSize,
		Version:         version,
		Type:            instanceType,

		// deprecated fields
		CPU: cpu,
		RAM: ram,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient mongodbflex.DefaultAPI) (mongodbflex.ApiPartialUpdateInstanceRequest, error) {
	var flavorId *string
	var err error

	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.InstanceId, model.Region)

	currentInstance, err := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId, model.Region).Execute()
	if err != nil {
		return req, fmt.Errorf("get MongoDB Flex instance: %w", err)
	}

	flavors, err := apiClient.ListFlavors(ctx, model.ProjectId, model.Region).Execute()
	if err != nil {
		return req, fmt.Errorf("get MongoDB Flex flavors: %w", err)
	}

	// if cpu/ram flags are used instead of the flavor id flag
	if model.FlavorId == nil && (model.RAM != nil || model.CPU != nil) {
		ram := model.RAM
		cpu := model.CPU

		// if only one of the cpu/ram flags is set
		if model.RAM == nil || model.CPU == nil {
			var currentFlavor *mongodbflex.InstanceFlavor
			for _, f := range flavors.Flavors {
				if f.Id == currentInstance.Item.Flavor.Id {
					currentFlavor = &f
				}
			}

			if currentFlavor == nil {
				return req, fmt.Errorf("flavor %s not found", currentInstance.Item.Flavor.GetId())
			}

			if model.RAM == nil {
				ram = currentFlavor.Memory
			}
			if model.CPU == nil {
				cpu = currentFlavor.Cpu
			}
		}

		flavorId, err = mongodbflexUtils.LoadFlavorId(*cpu, *ram, &flavors.Flavors) //nolint:staticcheck // SA1019 - deprecated but still supported until 2027-03-07
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

	var payloadAcl *mongodbflex.ACL
	if model.ACL != nil {
		payloadAcl = &mongodbflex.ACL{Items: *model.ACL}
	}

	var payloadStorage *mongodbflex.Storage
	if model.StorageClass != nil || model.StorageSize != nil {
		payloadStorage = &mongodbflex.Storage{
			Class: model.StorageClass,
			Size:  model.StorageSize,
		}
	}

	var replicas *int32
	var payloadOptions *map[string]string
	if model.Type != nil {
		replicasInt, err := mongodbflexUtils.GetInstanceReplicas(*model.Type)
		if err != nil {
			return req, fmt.Errorf("get PostgreSQL Flex instance type: %w", err)
		}

		replicas = &replicasInt
		payloadOptions = utils.Ptr(map[string]string{
			"type": *model.Type,
		})
	}

	req = req.PartialUpdateInstancePayload(mongodbflex.PartialUpdateInstancePayload{
		Name:           model.InstanceName,
		Acl:            payloadAcl,
		BackupSchedule: model.BackupSchedule,
		FlavorId:       flavorId,
		Replicas:       replicas,
		Storage:        payloadStorage,
		Version:        model.Version,
		Options:        payloadOptions,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat string, async bool, instanceLabel string, resp *mongodbflex.UpdateInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("resp is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		operationState := "Updated"
		if async {
			operationState = "Triggered update of"
		}
		p.Info("%s instance %q\n", operationState, instanceLabel)
		return nil
	})
}
