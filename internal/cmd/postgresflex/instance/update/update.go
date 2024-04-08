package update

import (
	"context"
	"errors"
	"fmt"

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
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

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
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId     string
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

func NewCmd(p *print.Printer) *cobra.Command {
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

			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			instanceLabel, err := postgresflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update instance %q? (This may cause downtime)", instanceLabel)
				err = p.PromptForConfirmation(cmd, prompt)
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
				return fmt.Errorf("update PostgreSQL Flex instance: %w", err)
			}
			instanceId := *resp.Item.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Updating instance")
				_, err = wait.PartialUpdateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for PostgreSQL Flex instance update: %w", err)
				}
				s.Stop()
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			p.Info("%s instance %q\n", operationState, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	typeFlagOptions := postgresflexUtils.AvailableInstanceTypes()

	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")
	cmd.Flags().String(backupScheduleFlag, "", "Backup schedule")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor")
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs")
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB)")
	cmd.Flags().String(storageClassFlag, "", "Storage class")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB)")
	cmd.Flags().String(versionFlag, "", "Version")
	cmd.Flags().Var(flags.EnumFlag(false, "", typeFlagOptions...), typeFlag, fmt.Sprintf("Instance type, one of %q", typeFlagOptions))
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	instanceName := flags.FlagToStringPointer(cmd, instanceNameFlag)
	flavorId := flags.FlagToStringPointer(cmd, flavorIdFlag)
	cpu := flags.FlagToInt64Pointer(cmd, cpuFlag)
	ram := flags.FlagToInt64Pointer(cmd, ramFlag)
	acl := flags.FlagToStringSlicePointer(cmd, aclFlag)
	backupSchedule := flags.FlagToStringPointer(cmd, backupScheduleFlag)
	storageClass := flags.FlagToStringPointer(cmd, storageClassFlag)
	storageSize := flags.FlagToInt64Pointer(cmd, storageSizeFlag)
	version := flags.FlagToStringPointer(cmd, versionFlag)
	instanceType := flags.FlagToStringPointer(cmd, typeFlag)

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

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		InstanceName:    instanceName,
		ACL:             acl,
		BackupSchedule:  backupSchedule,
		FlavorId:        flavorId,
		CPU:             cpu,
		RAM:             ram,
		StorageClass:    storageClass,
		StorageSize:     storageSize,
		Version:         version,
		Type:            instanceType,
	}, nil
}

type PostgreSQLFlexClient interface {
	PartialUpdateInstance(ctx context.Context, projectId, instanceId string) postgresflex.ApiPartialUpdateInstanceRequest
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*postgresflex.InstanceResponse, error)
	ListFlavorsExecute(ctx context.Context, projectId string) (*postgresflex.ListFlavorsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, flavorId string) (*postgresflex.ListStoragesResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient PostgreSQLFlexClient) (postgresflex.ApiPartialUpdateInstanceRequest, error) {
	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.InstanceId)

	var flavorId *string
	var err error

	flavors, err := apiClient.ListFlavorsExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get PostgreSQL Flex flavors: %w", err)
	}

	if model.FlavorId == nil && (model.RAM != nil || model.CPU != nil) {
		ram := model.RAM
		cpu := model.CPU
		if model.RAM == nil || model.CPU == nil {
			currentInstance, err := apiClient.GetInstanceExecute(ctx, model.ProjectId, model.InstanceId)
			if err != nil {
				return req, fmt.Errorf("get PostgreSQL Flex instance: %w", err)
			}
			if model.RAM == nil {
				ram = currentInstance.Item.Flavor.Memory
			}
			if model.CPU == nil {
				cpu = currentInstance.Item.Flavor.Cpu
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
		err := postgresflexUtils.ValidateFlavorId(*model.FlavorId, flavors.Flavors)
		if err != nil {
			return req, err
		}
		flavorId = model.FlavorId
	}

	var storages *postgresflex.ListStoragesResponse
	if model.StorageClass != nil || model.StorageSize != nil {
		validationFlavorId := flavorId
		if validationFlavorId == nil {
			currentInstance, err := apiClient.GetInstanceExecute(ctx, model.ProjectId, model.InstanceId)
			if err != nil {
				return req, fmt.Errorf("get PostgreSQL Flex instance: %w", err)
			}
			validationFlavorId = currentInstance.Item.Flavor.Id
		}
		storages, err = apiClient.ListStoragesExecute(ctx, model.ProjectId, *validationFlavorId)
		if err != nil {
			return req, fmt.Errorf("get PostgreSQL Flex storages: %w", err)
		}
		err = postgresflexUtils.ValidateStorage(model.StorageClass, model.StorageSize, storages, *validationFlavorId)
		if err != nil {
			return req, err
		}
	}

	var payloadAcl *postgresflex.ACL
	if model.ACL != nil {
		payloadAcl = &postgresflex.ACL{Items: model.ACL}
	}

	var payloadStorage *postgresflex.Storage
	if model.StorageClass != nil || model.StorageSize != nil {
		payloadStorage = &postgresflex.Storage{
			Class: model.StorageClass,
			Size:  model.StorageSize,
		}
	}

	var replicas *int64
	var payloadOptions *map[string]string
	if model.Type != nil {
		replicasInt, err := postgresflexUtils.GetInstanceReplicas(*model.Type)
		if err != nil {
			return req, fmt.Errorf("get PostgreSQL Flex instance type: %w", err)
		}

		replicas = &replicasInt
		payloadOptions = utils.Ptr(map[string]string{
			"type": *model.Type,
		})
	}

	req = req.PartialUpdateInstancePayload(postgresflex.PartialUpdateInstancePayload{
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
