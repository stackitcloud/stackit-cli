package create

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/client"
	mongodbflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mongodbflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/wait"
)

const (
	instanceNameFlag   = "name"
	aclFlag            = "acl"
	backupScheduleFlag = "backup-schedule"
	flavorIdFlag       = "flavor-id"
	cpuFlag            = "cpu"
	ramFlag            = "ram"
	storageClassFlag   = "storage-class"
	storageSizeFlag    = "storage-size"
	versionFlag        = "version"
	typeFlag           = "type"

	defaultBackupSchedule = "0 0/6 * * *"
	defaultStorageClass   = "premium-perf2-mongodb"
	defaultStorageSize    = 10
	defaultType           = "Replica"
	defaultVersion        = "6.0"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceName   *string
	ACL            *[]string
	BackupSchedule *string
	FlavorId       *string
	CPU            *int64
	RAM            *int64
	StorageClass   *string
	StorageSize    *int64
	Version        *string
	Type           *string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a MongoDB Flex instance",
		Long:  "Create a MongoDB Flex instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a MongoDB Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access) and specify flavor by CPU and RAM. Other parameters are set to default values`,
				`$ stackit mongodbflex instance create --name my-instance --cpu 1 --ram 4 --acl 0.0.0.0/0`),
			examples.NewExample(
				`Create a MongoDB Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access) and specify flavor by ID. Other parameters are set to default values`,
				`$ stackit mongodbflex instance create --name my-instance --flavor-id xxx --acl 0.0.0.0/0`),
			examples.NewExample(
				`Create a MongoDB Flex instance with name "my-instance", allow access to a specific range of IP addresses, specify flavor by CPU and RAM and set storage size to 20 GB. Other parameters are set to default values`,
				`$ stackit mongodbflex instance create --name my-instance --cpu 1 --ram 4 --acl 1.2.3.0/24 --storage-size 20`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a MongoDB Flex instance for project %s?", projectLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
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
				s := spinner.New(cmd)
				s.Start("Creating instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for MongoDB Flex instance creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			cmd.Printf("%s instance for project %s. Instance ID: %s\n", operationState, projectLabel, instanceId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	typeFlagOptions := mongodbflexUtils.AvailableInstanceTypes()

	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc.")
	cmd.Flags().String(backupScheduleFlag, defaultBackupSchedule, "Backup schedule")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor")
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs")
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB)")
	cmd.Flags().String(storageClassFlag, defaultStorageClass, "Storage class")
	cmd.Flags().Int64(storageSizeFlag, defaultStorageSize, "Storage size (in GB)")
	cmd.Flags().String(versionFlag, defaultVersion, "Version")
	cmd.Flags().Var(flags.EnumFlag(false, defaultType, typeFlagOptions...), typeFlag, fmt.Sprintf("Instance type, one of %q", typeFlagOptions))

	err := flags.MarkFlagsRequired(cmd, instanceNameFlag, aclFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	storageSize := flags.FlagWithDefaultToInt64Value(cmd, storageSizeFlag)

	flavorId := flags.FlagToStringPointer(cmd, flavorIdFlag)
	cpu := flags.FlagToInt64Pointer(cmd, cpuFlag)
	ram := flags.FlagToInt64Pointer(cmd, ramFlag)

	if flavorId == nil && (cpu == nil || ram == nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Service:   "mongodbflex",
			Operation: cmd.Use,
		}
	}
	if flavorId != nil && (cpu != nil || ram != nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Service:   "mongodbflex",
			Operation: cmd.Use,
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceName:    flags.FlagToStringPointer(cmd, instanceNameFlag),
		ACL:             flags.FlagToStringSlicePointer(cmd, aclFlag),
		BackupSchedule:  utils.Ptr(flags.FlagWithDefaultToStringValue(cmd, backupScheduleFlag)),
		FlavorId:        flavorId,
		CPU:             cpu,
		RAM:             ram,
		StorageClass:    utils.Ptr(flags.FlagWithDefaultToStringValue(cmd, storageClassFlag)),
		StorageSize:     &storageSize,
		Version:         utils.Ptr(flags.FlagWithDefaultToStringValue(cmd, versionFlag)),
		Type:            utils.Ptr(flags.FlagWithDefaultToStringValue(cmd, typeFlag)),
	}, nil
}

type MongoDBFlexClient interface {
	CreateInstance(ctx context.Context, projectId string) mongodbflex.ApiCreateInstanceRequest
	ListFlavorsExecute(ctx context.Context, projectId string) (*mongodbflex.ListFlavorsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, flavorId string) (*mongodbflex.ListStoragesResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient MongoDBFlexClient) (mongodbflex.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId)

	var flavorId *string
	var err error

	flavors, err := apiClient.ListFlavorsExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get MongoDB Flex flavors: %w", err)
	}

	if model.FlavorId == nil {
		flavorId, err = mongodbflexUtils.LoadFlavorId(*model.CPU, *model.RAM, flavors.Flavors)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load flavor ID: %w", err)
			}
			return req, err
		}
	} else {
		err := mongodbflexUtils.ValidateFlavorId(*model.FlavorId, flavors.Flavors)
		if err != nil {
			return req, err
		}
		flavorId = model.FlavorId
	}

	storages, err := apiClient.ListStoragesExecute(ctx, model.ProjectId, *flavorId)
	if err != nil {
		return req, fmt.Errorf("get MongoDB Flex storages: %w", err)
	}
	err = mongodbflexUtils.ValidateStorage(model.StorageClass, model.StorageSize, storages, *flavorId)
	if err != nil {
		return req, err
	}

	replicas, err := mongodbflexUtils.GetInstanceReplicas(*model.Type)
	if err != nil {
		return req, fmt.Errorf("get MongoDB Flex intance type: %w", err)
	}

	req = req.CreateInstancePayload(mongodbflex.CreateInstancePayload{
		Name:           model.InstanceName,
		Acl:            &mongodbflex.ACL{Items: model.ACL},
		BackupSchedule: model.BackupSchedule,
		FlavorId:       flavorId,
		Replicas:       &replicas,
		Storage: &mongodbflex.Storage{
			Class: model.StorageClass,
			Size:  model.StorageSize,
		},
		Version: model.Version,
		Options: utils.Ptr(map[string]string{
			"type": *model.Type,
		}),
	})
	return req, nil
}
