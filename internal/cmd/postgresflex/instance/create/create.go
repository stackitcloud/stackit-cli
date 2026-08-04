package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	postgresflex "github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/v3api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/client"
	postgresflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
)

const (
	instanceNameFlag   = "name"
	aclFlag            = "acl"
	backupScheduleFlag = "backup-schedule"
	flavorIdFlag       = "flavor-id"
	storageClassFlag   = "storage-class"
	storageSizeFlag    = "storage-size"
	versionFlag        = "version"
	retentionDaysFlag  = "retention-days"

	defaultBackupSchedule = "0 0 * * *"             // Deprecated: Will be removed after 2027-01-31.
	defaultStorageSize    = 10                      // Deprecated: Will be removed after 2027-01-31.
	defaultStorageClass   = "premium-perf2-stackit" // Deprecated: Will be removed after 2027-01-31.

	encryptionKekKeyIdFlag       = "encryption-kek-key-id"
	encryptionKekKeyringIdFlag   = "encryption-kek-keyring-id"
	encryptionKekKeyVersionFlag  = "encryption-kek-key-version"
	encryptionServiceAccountFlag = "encryption-service-account"

	cpuFlag     = "cpu"     // Deprecated: Will be removed after 2027-01-31. Flavor id should be used instead.
	ramFlag     = "ram"     // Deprecated: Will be removed after 2027-01-31. Flavor id should be used instead.
	defaultType = "Replica" // Deprecated: Will be removed after 2027-01-31. Replicas are managed via the flavor id on API side now.
)

var (
	// Deprecated: Will be removed after 2027-01-31. Replicas are managed via the flavor id on API side now.
	typeFlag = flags.StringEnumFlag(
		"type",
		postgresflexUtils.AvailableInstanceTypes(),
		"Instance type,",
		flags.StringEnumDefaultValue(defaultType),
	)
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
	RetentionDays  *int32

	EncryptionKekKeyId       *string
	EncryptionKekKeyringId   *string
	EncryptionKekKeyVersion  *string
	EncryptionServiceAccount *string

	CPU  *int64 // Deprecated: Will be removed after 2027-01-31
	RAM  *int64 // Deprecated: Will be removed after 2027-01-31
	Type string // Deprecated: Will be removed after 2027-01-31
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a PostgreSQL Flex instance",
		Long:  "Creates a PostgreSQL Flex instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a PostgreSQL Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access).`,
				`$ stackit postgresflex instance create --name my-instance --flavor-id xxx --acl 0.0.0.0/0 --storage-size 20 --retention-days 32 --version 17 --backup-schedule "6 6 * * *" --storage-size 10 --storage-class premium-perf2-stackit`),
			examples.NewExample(
				`Create a PostgreSQL Flex instance with name "my-instance", allow access to a specific range of IP addresses.`,
				`$ stackit postgresflex instance create --name my-instance --flavor-id xxx --acl 1.2.3.0/24 --storage-size 20 --retention-days 32 --version 17 --backup-schedule "6 6 * * *" --storage-size 10 --storage-class premium-perf2-stackit`),
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

			// load flavor id - remove after 2027-01-31
			if model.FlavorId == nil {
				// transform the model.FlavorId field from "*string" to "string" once this is removed
				params.Printer.Warn("The --%s flag is not set, determining flavor ID by CPU und RAM. This behavior is deprecated, the --%s flag will be required after 2027-01-31.\n", flavorIdFlag, flavorIdFlag)
			}
			model.FlavorId, err = getFlavorId(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "determining flavor id: %v", err)
			}

			// remove after 2027-01-31
			if model.RetentionDays == nil {
				params.Printer.Warn("The --%s flag is not set. Using the default value (32). This behavior is deprecated, the --%s flag will be required after 2027-01-31.\n", retentionDaysFlag, retentionDaysFlag)
				model.RetentionDays = utils.Ptr(int32(32)) // transform the model.RetentionDays field from "*int32" to "int32" once this is removed
			}

			// remove after 2027-01-31
			if model.BackupSchedule == nil {
				params.Printer.Warn("The --%s flag is not set. Using the default value \"%s\". This behavior is deprecated, the --%s flag will be required after 2027-01-31.\n", backupScheduleFlag, defaultBackupSchedule, backupScheduleFlag)
				model.BackupSchedule = utils.Ptr(defaultBackupSchedule) // transform the model.BackupSchedule field from "*string" to "string" once this is removed
			}

			// Fill in version, if needed - remove after 2027-01-31
			if model.Version == nil {
				params.Printer.Warn("The --%s flag is not set. Using the latest version as a default. This behavior is deprecated, the --%s flag will be required after 2027-01-31.\n", versionFlag, versionFlag)

				version, err := postgresflexUtils.GetLatestPostgreSQLVersion(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region)
				if err != nil {
					return fmt.Errorf("get latest PostgreSQL version: %w", err)
				}
				model.Version = utils.Ptr(version) // transform the model.Version field from "*string" to "string" once this is removed
			}

			if model.StorageSize == nil {
				params.Printer.Warn("The --%s flag is not set. Using the default value (%d). This behavior is deprecated, the --%s flag will be required after 2027-01-31.\n", storageSizeFlag, defaultStorageSize, storageSizeFlag)
				model.StorageSize = utils.Ptr(int64(defaultStorageSize))
			}

			if model.StorageClass == nil {
				params.Printer.Warn("The --%s flag is not set. Using the default value (%s). This behavior is deprecated, the --%s flag will be required after 2027-01-31.\n", storageClassFlag, defaultStorageClass, storageClassFlag)
				model.StorageClass = utils.Ptr(defaultStorageClass)
			}

			prompt := fmt.Sprintf("Are you sure you want to create a PostgreSQL Flex instance for project %q?", projectLabel)
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
				return fmt.Errorf("create PostgreSQL Flex instance: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating instance", func() error {
					_, err = wait.CreateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, resp.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for PostgreSQL Flex instance creation: %w", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp.Id, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc.")
	cmd.Flags().String(backupScheduleFlag, "", "Backup schedule. This flag will be required after 2027-01-31.")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor. This flag will be required after 2027-01-31.")
	cmd.Flags().String(storageClassFlag, "", "Storage class. This flag will be required after 2027-01-31.")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB). This flag will be required after 2027-01-31.")
	cmd.Flags().String(versionFlag, "", "PostgreSQL version. Defaults to the latest version available. This flag will be required after 2027-01-31.")
	cmd.Flags().Int32(retentionDaysFlag, 0, "The days for how long the backup files should be stored before cleaned up (32 to 90). This flag will be required after 2027-01-31.")
	cmd.Flags().String(encryptionKekKeyIdFlag, "", "The key identifier")
	cmd.Flags().String(encryptionKekKeyringIdFlag, "", "The keyring identifier")
	cmd.Flags().String(encryptionKekKeyVersionFlag, "", "The key version")
	cmd.Flags().String(encryptionServiceAccountFlag, "", "The service account")

	// remove after 2027-01-31
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs.")
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB).")
	typeFlag.Register(cmd.Flags())

	// after 2027-01-31: add backup-schedule,storage-size,version,flavor-id,retention-days
	err := flags.MarkFlagsRequired(cmd, instanceNameFlag, aclFlag)
	cobra.CheckErr(err)

	cmd.MarkFlagsRequiredTogether(encryptionKekKeyIdFlag, encryptionKekKeyringIdFlag, encryptionKekKeyVersionFlag, encryptionServiceAccountFlag)

	// remove after 2027-01-31
	err = cmd.Flags().MarkDeprecated(typeFlag.Name(), fmt.Sprintf("Will be removed after 2027-01-31. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	err = cmd.Flags().MarkDeprecated(cpuFlag, fmt.Sprintf("Will be removed after 2027-01-31. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	err = cmd.Flags().MarkDeprecated(ramFlag, fmt.Sprintf("Will be removed after 2027-01-31. Use the --%s flag instead.", flavorIdFlag))
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
	cpu := flags.FlagToInt64Pointer(p, cmd, cpuFlag)
	ram := flags.FlagToInt64Pointer(p, cmd, ramFlag)

	// remove after 2027-01-31: flavor id flag will be required then
	if flavorId == nil && (cpu == nil || ram == nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd: cmd,
		}
	}
	// remove after 2027-01-31: flavor id flag will be required then
	if flavorId != nil && (cpu != nil || ram != nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd: cmd,
		}
	}

	// remove after 2027-01-31: storage size flag will be required then (no pointer anymore)
	var storageSize *int64
	if cmd.Flags().Changed(storageSizeFlag) {
		storageSize = flags.FlagToInt64Pointer(p, cmd, storageSizeFlag)
	}

	// remove after 2027-01-31: retention days flag will be required then (no pointer anymore)
	var retentionDays *int32
	if cmd.Flags().Changed(retentionDaysFlag) {
		retentionDays = flags.FlagToInt32Pointer(p, cmd, retentionDaysFlag)
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
		RetentionDays:   retentionDays,

		EncryptionKekKeyId:       flags.FlagToStringPointer(p, cmd, encryptionKekKeyIdFlag),
		EncryptionKekKeyringId:   flags.FlagToStringPointer(p, cmd, encryptionKekKeyringIdFlag),
		EncryptionKekKeyVersion:  flags.FlagToStringPointer(p, cmd, encryptionKekKeyVersionFlag),
		EncryptionServiceAccount: flags.FlagToStringPointer(p, cmd, encryptionServiceAccountFlag),

		// deprecated fields
		CPU:  cpu,
		RAM:  ram,
		Type: typeFlag.Get(),
	}

	p.DebugInputModel(model)
	return &model, nil
}

// Deprecated: remove after 2027-01-31
func getFlavorId(ctx context.Context, model *inputModel, apiClient postgresflex.DefaultAPI) (*string, error) {
	if model == nil {
		return nil, fmt.Errorf("model is nil")
	}

	if model.FlavorId != nil {
		return model.FlavorId, nil
	}

	// Load all flavors
	flavors, err := apiClient.ListFlavors(ctx, model.ProjectId, model.Region).Size(100).Execute()
	if err != nil {
		return nil, fmt.Errorf("loading flavors: %w", err)
	}

	for _, flavor := range flavors.Flavors {
		if flavor.Cpu == *model.CPU && flavor.Memory == *model.RAM && flavor.NodeType == model.Type {
			return &flavor.Id, nil
		}
	}

	return nil, fmt.Errorf("no matching flavor found")
}

func buildRequest(ctx context.Context, model *inputModel, apiClient postgresflex.DefaultAPI) (postgresflex.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId, model.Region)

	var encryption *postgresflex.InstanceEncryption
	if model.EncryptionKekKeyId != nil && model.EncryptionKekKeyringId != nil && model.EncryptionKekKeyVersion != nil && model.EncryptionServiceAccount != nil {
		encryption = &postgresflex.InstanceEncryption{
			KekKeyId:       *model.EncryptionKekKeyId,
			KekKeyRingId:   *model.EncryptionKekKeyringId,
			KekKeyVersion:  *model.EncryptionKekKeyVersion,
			ServiceAccount: *model.EncryptionServiceAccount,
		}
	}

	// remove after 2027-01-31
	if model.BackupSchedule == nil {
		return postgresflex.ApiCreateInstanceRequest{}, fmt.Errorf("backup schedule is nil")
	} else if model.StorageSize == nil {
		return postgresflex.ApiCreateInstanceRequest{}, fmt.Errorf("storage size is nil")
	} else if model.Version == nil {
		return postgresflex.ApiCreateInstanceRequest{}, fmt.Errorf("version is nil")
	}

	req = req.CreateInstancePayload(postgresflex.CreateInstancePayload{
		BackupSchedule: *model.BackupSchedule,
		Encryption:     encryption,
		FlavorId:       utils.PtrString(model.FlavorId),
		Name:           model.InstanceName,
		Network: postgresflex.InstanceNetworkCreate{
			Acl: model.ACL,
		},
		RetentionDays: *postgresflex.NewNullableInt32(model.RetentionDays),
		Storage: postgresflex.StorageCreate{
			Class: model.StorageClass,
			Size:  *model.StorageSize,
		},
		Version: *model.Version,
	})
	return req, nil
}

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel, instanceId string, resp *postgresflex.CreateInstanceResponse) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("no response passed")
		}
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, instanceId)
		return nil
	})
}
