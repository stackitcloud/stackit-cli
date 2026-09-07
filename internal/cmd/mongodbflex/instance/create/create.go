package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	mongodbflex "github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/v2api"
	wait "github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/v2api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	instanceNameFlag   = "name"
	aclFlag            = "acl"
	backupScheduleFlag = "backup-schedule"
	flavorIdFlag       = "flavor-id"
	storageClassFlag   = "storage-class"
	storageSizeFlag    = "storage-size"
	versionFlag        = "version"
	defaultType        = "Replica"

	cpuFlag = "cpu" // Deprecated: Will be removed after 2027-03-07.
	ramFlag = "ram" // Deprecated: Will be removed after 2027-03-07.

	defaultBackupSchedule = "0 0/6 * * *"           // Deprecated: Will be removed after 2027-03-07.
	defaultStorageClass   = "premium-perf2-mongodb" // Deprecated: Will be removed after 2027-03-07.
	defaultStorageSize    = 10                      // Deprecated: Will be removed after 2027-03-07.

)

var typeFlag = flags.StringEnumFlag(
	"type",
	mongodbflexUtils.AvailableInstanceTypes(),
	"Instance type,",
	flags.StringEnumDefaultValue(defaultType),
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceName   string
	ACL            []string
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
		Use:   "create",
		Short: "Creates a MongoDB Flex instance",
		Long:  "Creates a MongoDB Flex instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a MongoDB Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access).`,
				`$ stackit mongodbflex instance create --name my-instance --flavor-id xxx --acl 0.0.0.0/0 --type Replica --storage-size 20 --version 8.0 --backup-schedule "6 6 * * *" --storage-size 10 --storage-class premium-perf2-mongodb`),
			examples.NewExample(
				`Create a MongoDB Flex instance with name "my-instance", allow access to a specific range of IP addresses.`,
				`$ stackit mongodbflex instance create --name my-instance --flavor-id xxx --acl 1.2.3.0/24 --type Replica --storage-size 20 --version 8.0 --backup-schedule "6 6 * * *" --storage-size 10 --storage-class premium-perf2-mongodb`),
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

			// load flavor id - remove after 2027-03-07
			if model.FlavorId == nil {
				// transform the model.FlavorId field from "*string" to "string" once this is removed
				params.Printer.Warn("The --%s flag is not set, determining flavor ID by CPU und RAM. This behavior is deprecated, the --%s flag will be required after 2027-03-07.\n", flavorIdFlag, flavorIdFlag)
			}
			model.FlavorId, err = getFlavorId(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "determining flavor id: %v", err)
			}

			// remove after 2027-03-07
			if model.BackupSchedule == nil {
				// transform the model.BackupSchedule field from "*string" to "string" once this is removed
				params.Printer.Warn("The --%s flag is not set. Using the default value \"%s\". This behavior is deprecated, the --%s flag will be required after 2027-03-07.\n", backupScheduleFlag, defaultBackupSchedule, backupScheduleFlag)
				model.BackupSchedule = utils.Ptr(defaultBackupSchedule)
			}

			// Fill in version, if needed - remove after 2027-03-07
			if model.Version == nil {
				params.Printer.Warn("The --%s flag is not set. Using the latest version as a default. This behavior is deprecated, the --%s flag will be required after 2027-03-07.\n", versionFlag, versionFlag)
				// transform the model.Version field from "*string" to "string" once this is removed

				version, err := mongodbflexUtils.GetLatestMongoDBVersion(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region) //nolint:staticcheck // SA1019 - deprecated but still supported until 2027-03-07
				if err != nil {
					return fmt.Errorf("get latest MongoDB version: %w", err)
				}
				model.Version = utils.Ptr(version)
			}

			// remove after 2027-03-07
			if model.StorageSize == nil {
				params.Printer.Warn("The --%s flag is not set. Using the default value (%d). This behavior is deprecated, the --%s flag will be required after 2027-03-07.\n", storageSizeFlag, defaultStorageSize, storageSizeFlag)
				model.StorageSize = utils.Ptr(int64(defaultStorageSize))
			}

			// remove after 2027-03-07
			if model.StorageClass == nil {
				params.Printer.Warn("The --%s flag is not set. Using the default value (%s). This behavior is deprecated, the --%s flag will be required after 2027-03-07.\n", storageClassFlag, defaultStorageClass, storageClassFlag)
				model.StorageClass = utils.Ptr(defaultStorageClass)
			}

			prompt := fmt.Sprintf("Are you sure you want to create a MongoDB Flex instance for project %q?", projectLabel)
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
				return fmt.Errorf("create MongoDB Flex instance: %w", err)
			}
			instanceId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating instance", func() error {
					_, err = wait.CreateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, instanceId, model.Region).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for MongoDB Flex instance creation: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc.")
	cmd.Flags().String(backupScheduleFlag, defaultBackupSchedule, "Backup schedule. This flag will be required after 2027-03-07.")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor. This flag will be required after 2027-03-07.")
	cmd.Flags().String(storageClassFlag, defaultStorageClass, "Storage class. This flag will be required after 2027-03-07.")
	cmd.Flags().Int64(storageSizeFlag, defaultStorageSize, "Storage size (in GB). This flag will be required after 2027-03-07.")
	cmd.Flags().String(versionFlag, "", "MongoDB version. Defaults to the latest version available. This flag will be required after 2027-03-07.")
	typeFlag.Register(cmd.Flags())

	// remove after 2027-03-07
	cmd.Flags().Int32(cpuFlag, 0, "Number of CPUs")
	cmd.Flags().Int32(ramFlag, 0, "Amount of RAM (in GB)")

	// after 2027-03-07: add backupScheduleFlag, storageClassFlag, storageSizeFlag, versionFlag, flavorIdFlag, replicasFlag
	err := flags.MarkFlagsRequired(cmd, instanceNameFlag, aclFlag)
	cobra.CheckErr(err)

	// remove after 2027-03-07
	err = cmd.Flags().MarkDeprecated(cpuFlag, fmt.Sprintf("Will be removed after 2027-03-07. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	err = cmd.Flags().MarkDeprecated(ramFlag, fmt.Sprintf("Will be removed after 2027-03-07. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	cmd.MarkFlagsRequiredTogether(cpuFlag, ramFlag)
	cmd.MarkFlagsMutuallyExclusive(flavorIdFlag, cpuFlag)
	cmd.MarkFlagsMutuallyExclusive(flavorIdFlag, ramFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	flavorId := flags.FlagToStringPointer(p, cmd, flavorIdFlag)
	cpu := flags.FlagToInt32Pointer(p, cmd, cpuFlag)
	ram := flags.FlagToInt32Pointer(p, cmd, ramFlag)

	// remove after 2027-03-07: flavor id flag will be required then
	if flavorId == nil && (cpu == nil || ram == nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd: cmd,
		}
	}
	// remove after 2027-03-07: flavor id flag will be required then
	if flavorId != nil && (cpu != nil || ram != nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd: cmd,
		}
	}

	// remove after 2027-03-07: storage size flag will be required then (no pointer anymore)
	var storageSize *int64
	if cmd.Flags().Changed(storageSizeFlag) {
		storageSize = flags.FlagToInt64Pointer(p, cmd, storageSizeFlag)
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceName:    flags.FlagToStringValue(p, cmd, instanceNameFlag),
		ACL:             flags.FlagToStringSliceValue(p, cmd, aclFlag),
		BackupSchedule:  flags.FlagToStringPointer(p, cmd, backupScheduleFlag),
		FlavorId:        flavorId,
		StorageClass:    flags.FlagToStringPointer(p, cmd, storageClassFlag),
		StorageSize:     storageSize,
		Version:         flags.FlagToStringPointer(p, cmd, versionFlag),
		Type:            typeFlag.Ptr(),

		// remove after 2027-03-07: deprecated fields
		CPU: cpu,
		RAM: ram,
	}

	p.DebugInputModel(model)
	return &model, nil
}

// Deprecated: remove after 2027-03-07
func getFlavorId(ctx context.Context, model *inputModel, apiClient mongodbflex.DefaultAPI) (*string, error) {
	if model == nil {
		return nil, fmt.Errorf("model is nil")
	}

	if model.FlavorId != nil {
		return model.FlavorId, nil
	}

	// Load all flavors
	flavors, err := apiClient.ListFlavors(ctx, model.ProjectId, model.Region).Execute()
	if err != nil {
		return nil, fmt.Errorf("loading flavors: %w", err)
	}

	for _, flavor := range flavors.Flavors {
		if *flavor.Cpu == *model.CPU && *flavor.Memory == *model.RAM {
			return flavor.Id, nil
		}
	}

	return nil, fmt.Errorf("no matching flavor found")
}

func buildRequest(ctx context.Context, model *inputModel, apiClient mongodbflex.DefaultAPI) (mongodbflex.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId, model.Region)

	// remove after 2027-03-07
	if model.BackupSchedule == nil {
		return mongodbflex.ApiCreateInstanceRequest{}, fmt.Errorf("backup schedule is nil")
	} else if model.StorageSize == nil {
		return mongodbflex.ApiCreateInstanceRequest{}, fmt.Errorf("storage size is nil")
	} else if model.Version == nil {
		return mongodbflex.ApiCreateInstanceRequest{}, fmt.Errorf("version is nil")
	} else if model.StorageClass == nil {
		return mongodbflex.ApiCreateInstanceRequest{}, fmt.Errorf("storage class is nil")
	}

	replicas, err := mongodbflexUtils.GetInstanceReplicas(*model.Type)
	if err != nil {
		return req, fmt.Errorf("get MongoDB Flex instance type: %w", err)
	}

	req = req.CreateInstancePayload(mongodbflex.CreateInstancePayload{
		Name:           model.InstanceName,
		Acl:            mongodbflex.ACL{Items: model.ACL},
		BackupSchedule: *model.BackupSchedule,
		FlavorId:       utils.PtrString(model.FlavorId),
		Replicas:       replicas,
		Storage: mongodbflex.Storage{
			Class: model.StorageClass,
			Size:  model.StorageSize,
		},
		Version: *model.Version,
		Options: map[string]string{
			"type": *model.Type,
		},
	})

	return req, nil
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, resp *mongodbflex.CreateInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("create instance response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, utils.PtrString(resp.Id))
		return nil
	})
}
