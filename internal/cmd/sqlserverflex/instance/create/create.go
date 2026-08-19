package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	sqlserverflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api/wait"
)

const (
	instanceNameFlag   = "name"
	aclFlag            = "acl"
	backupScheduleFlag = "backup-schedule"
	flavorIdFlag       = "flavor-id"
	storageClassFlag   = "storage-class"
	storageSizeFlag    = "storage-size"
	versionFlag        = "version"
	editionFlag        = "edition"
	retentionDaysFlag  = "retention-days"

	encryptionKekKeyIdFlag       = "encryption-kek-key-id"
	encryptionKekKeyringIdFlag   = "encryption-kek-keyring-id"
	encryptionKekKeyVersionFlag  = "encryption-kek-key-version"
	encryptionServiceAccountFlag = "encryption-service-account"

	// Deprecated: cpuFlag is deprecated and will be removed after 2027-02-28. Use flavorIdFlag instead.
	cpuFlag = "cpu"
	// Deprecated: ramFlag is deprecated and will be removed after 2027-02-28. Use flavorIdFlag instead.
	ramFlag = "ram"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceName   string
	ACL            []string
	BackupSchedule string
	FlavorId       *string
	StorageClass   string
	StorageSize    *int64
	Version        string
	RetentionDays  int32

	EncryptionKekKeyId       *string
	EncryptionKekKeyringId   *string
	EncryptionKekKeyVersion  *string
	EncryptionServiceAccount *string

	// Deprecated: CPU is deprecated and will be removed after 2027-02-28. Use FlavorId instead.
	CPU *int64
	// Deprecated: RAM is deprecated and will be removed after 2027-02-28. Use FlavorId instead.
	RAM *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a SQLServer Flex instance",
		Long:  "Creates a SQLServer Flex instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a SQLServer Flex instance with name "my-instance" and specify flavor by ID.
  The flavor ID can be retrieved by running "$ stackit sqlserverflex flavor list"`,
				`$ stackit sqlserverflex instance create --name my-instance --flavor-id xxx --backup-schedule "0 2 * * *" --retention-days 30 --storage-class premium-perf2-stackit --storage-size 10 --version 2022 --acl 1.2.3.0/24`),
			examples.NewExample(
				`Create a SQLServer Flex instance with name "my-instance", specify flavor by ID, set storage size to 20 GB, and restrict access to a specific range of IP addresses.`,
				`$ stackit sqlserverflex instance create --name my-instance --flavor-id xxx --storage-size 20 --backup-schedule "0 2 * * *" --retention-days 30 --storage-class premium-perf2-stackit --version 2022 --acl 1.2.3.0/24`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Deprecated: remove after 2027-02-28, once flavor-id is the only supported way to select a flavor.
			if model.FlavorId == nil {
				params.Printer.Warn("The --%s flag is not set, determining flavor ID by CPU and RAM. This behavior is deprecated, the --%s flag will be required after 2027-02-28.\n", flavorIdFlag, flavorIdFlag)
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

			prompt := fmt.Sprintf("Are you sure you want to create a SQLServer Flex instance for project %q?", projectLabel)
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
				return fmt.Errorf("create SQLServer Flex instance: %w", err)
			}
			instanceId := resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating instance", func() error {
					_, err = wait.CreateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, instanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for SQLServer Flex instance creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc.")
	cmd.Flags().String(backupScheduleFlag, "", "Backup schedule")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor. This flag will be required after 2027-02-28.")
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs")
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB)")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB)")
	cmd.Flags().String(storageClassFlag, "", "Storage class")
	cmd.Flags().String(versionFlag, "", "SQLServer version")
	cmd.Flags().String(editionFlag, "", "Edition of the SQLServer instance")
	cmd.Flags().Int32(retentionDaysFlag, 0, "The days for how long the backup files should be stored before being cleaned up")
	cmd.Flags().String(encryptionKekKeyIdFlag, "", "The key identifier")
	cmd.Flags().String(encryptionKekKeyringIdFlag, "", "The keyring identifier")
	cmd.Flags().String(encryptionKekKeyVersionFlag, "", "The key version")
	cmd.Flags().String(encryptionServiceAccountFlag, "", "The service account")

	err := flags.MarkFlagsRequired(cmd, instanceNameFlag, backupScheduleFlag, retentionDaysFlag, storageClassFlag, storageSizeFlag, versionFlag)
	cobra.CheckErr(err)

	cmd.MarkFlagsRequiredTogether(encryptionKekKeyIdFlag, encryptionKekKeyringIdFlag, encryptionKekKeyVersionFlag, encryptionServiceAccountFlag)

	// Deprecated: remove after 2027-02-28
	err = cmd.Flags().MarkDeprecated(cpuFlag, fmt.Sprintf("will be removed after 2027-02-28. Use the --%s flag instead.", flavorIdFlag))
	cobra.CheckErr(err)
	err = cmd.Flags().MarkDeprecated(ramFlag, fmt.Sprintf("will be removed after 2027-02-28. Use the --%s flag instead.", flavorIdFlag))
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

	if flavorId == nil && (cpu == nil || ram == nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd:     cmd,
			Service: sqlserverflexUtils.ServiceCmd,
		}
	}
	if flavorId != nil && (cpu != nil || ram != nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd:     cmd,
			Service: sqlserverflexUtils.ServiceCmd,
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceName:    flags.FlagToStringValue(p, cmd, instanceNameFlag),
		ACL:             flags.FlagToStringSliceValue(p, cmd, aclFlag),
		BackupSchedule:  flags.FlagToStringValue(p, cmd, backupScheduleFlag),
		FlavorId:        flavorId,
		CPU:             cpu,
		RAM:             ram,
		StorageClass:    flags.FlagToStringValue(p, cmd, storageClassFlag),
		StorageSize:     flags.FlagToInt64Pointer(p, cmd, storageSizeFlag),
		Version:         flags.FlagToStringValue(p, cmd, versionFlag),
		RetentionDays:   flags.FlagWithDefaultToInt32Value(p, cmd, retentionDaysFlag),

		EncryptionKekKeyId:       flags.FlagToStringPointer(p, cmd, encryptionKekKeyIdFlag),
		EncryptionKekKeyringId:   flags.FlagToStringPointer(p, cmd, encryptionKekKeyringIdFlag),
		EncryptionKekKeyVersion:  flags.FlagToStringPointer(p, cmd, encryptionKekKeyVersionFlag),
		EncryptionServiceAccount: flags.FlagToStringPointer(p, cmd, encryptionServiceAccountFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient sqlserverflex.DefaultAPI) (sqlserverflex.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId, model.Region)

	var flavorId string
	var err error

	flavors, err := sqlserverflexUtils.ListAllFlavors(ctx, apiClient, model.ProjectId, model.Region)
	if err != nil {
		return req, fmt.Errorf("get SQLServer Flex flavors: %w", err)
	}

	if model.FlavorId == nil {
		flavorId, err = sqlserverflexUtils.LoadFlavorId(*model.CPU, *model.RAM, flavors)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load flavor ID: %w", err)
			}
			return req, err
		}
	} else {
		err := sqlserverflexUtils.ValidateFlavorId(*model.FlavorId, flavors)
		if err != nil {
			return req, err
		}
		flavorId = *model.FlavorId
	}

	storages, err := apiClient.ListStorages(ctx, model.ProjectId, model.Region, flavorId).Execute()
	if err != nil {
		return req, fmt.Errorf("get SQLServer Flex storages: %w", err)
	}
	err = sqlserverflexUtils.ValidateStorage(model.StorageClass, model.StorageSize, storages, flavorId)
	if err != nil {
		return req, err
	}

	var encryption *sqlserverflex.InstanceEncryption
	if model.EncryptionKekKeyId != nil && model.EncryptionKekKeyringId != nil && model.EncryptionKekKeyVersion != nil && model.EncryptionServiceAccount != nil {
		encryption = &sqlserverflex.InstanceEncryption{
			KekKeyId:       *model.EncryptionKekKeyId,
			KekKeyRingId:   *model.EncryptionKekKeyringId,
			KekKeyVersion:  *model.EncryptionKekKeyVersion,
			ServiceAccount: *model.EncryptionServiceAccount,
		}
		if model.ACL == nil {
			model.ACL = make([]string, 0)
		}
	}

	req = req.CreateInstancePayload(sqlserverflex.CreateInstancePayload{
		Name:       model.InstanceName,
		Encryption: encryption,
		Network: sqlserverflex.CreateInstancePayloadNetwork{
			Acl: model.ACL,
		},
		BackupSchedule: model.BackupSchedule,
		FlavorId:       flavorId,
		Storage: sqlserverflex.StorageCreate{
			Class: model.StorageClass,
			Size:  utils.PtrValue(model.StorageSize),
		},
		Version:       sqlserverflex.InstanceVersion(model.Version),
		RetentionDays: model.RetentionDays,
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *sqlserverflex.CreateInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("sqlserverflex response is empty")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, resp.Id)
		return nil
	})
}
