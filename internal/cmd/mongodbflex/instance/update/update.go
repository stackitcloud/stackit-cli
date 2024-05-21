package update

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
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
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex"
	"github.com/stackitcloud/stackit-sdk-go/services/mongodbflex/wait"
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
		Short: "Updates a MongoDB Flex instance",
		Long:  "Updates a MongoDB Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`Update the name of a MongoDB Flex instance`,
				"$ stackit mongodbflex instance update xxx --name my-new-name"),
			examples.NewExample(
				`Update the version of a MongoDB Flex instance`,
				"$ stackit mongodbflex instance update xxx --version 6.0"),
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			instanceLabel, err := mongodbflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update instance %q?", instanceLabel)
				err = p.PromptForConfirmation(prompt)
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
				return fmt.Errorf("update MongoDB Flex instance: %w", err)
			}
			instanceId := *resp.Item.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Updating instance")
				_, err = wait.PartialUpdateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for MongoDB Flex instance update: %w", err)
				}
				s.Stop()
			}

			return outputResult(p, model, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	typeFlagOptions := mongodbflexUtils.AvailableInstanceTypes()

	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "Lists of IP networks in CIDR notation which are allowed to access this instance")
	cmd.Flags().String(backupScheduleFlag, "", "Backup schedule")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor")
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs")
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB)")
	cmd.Flags().String(storageClassFlag, "", "Storage class")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB)")
	cmd.Flags().String(versionFlag, "", "Version")
	cmd.Flags().Var(flags.EnumFlag(false, "", typeFlagOptions...), typeFlag, fmt.Sprintf("Instance type, one of %q", typeFlagOptions))
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
	acl := flags.FlagToStringSlicePointer(p, cmd, aclFlag)
	backupSchedule := flags.FlagToStringPointer(p, cmd, backupScheduleFlag)
	storageClass := flags.FlagToStringPointer(p, cmd, storageClassFlag)
	storageSize := flags.FlagToInt64Pointer(p, cmd, storageSizeFlag)
	version := flags.FlagToStringPointer(p, cmd, versionFlag)
	instanceType := flags.FlagToStringPointer(p, cmd, typeFlag)

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
		CPU:             cpu,
		RAM:             ram,
		StorageClass:    storageClass,
		StorageSize:     storageSize,
		Version:         version,
		Type:            instanceType,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

type MongoDBFlexClient interface {
	PartialUpdateInstance(ctx context.Context, projectId, instanceId string) mongodbflex.ApiPartialUpdateInstanceRequest
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*mongodbflex.GetInstanceResponse, error)
	ListFlavorsExecute(ctx context.Context, projectId string) (*mongodbflex.ListFlavorsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, flavorId string) (*mongodbflex.ListStoragesResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient MongoDBFlexClient) (mongodbflex.ApiPartialUpdateInstanceRequest, error) {
	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.InstanceId)

	var flavorId *string
	var err error

	flavors, err := apiClient.ListFlavorsExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get MongoDB Flex flavors: %w", err)
	}

	if model.FlavorId == nil && (model.RAM != nil || model.CPU != nil) {
		ram := model.RAM
		cpu := model.CPU
		if model.RAM == nil || model.CPU == nil {
			currentInstance, err := apiClient.GetInstanceExecute(ctx, model.ProjectId, model.InstanceId)
			if err != nil {
				return req, fmt.Errorf("get MongoDB Flex instance: %w", err)
			}
			if model.RAM == nil {
				ram = currentInstance.Item.Flavor.Memory
			}
			if model.CPU == nil {
				cpu = currentInstance.Item.Flavor.Cpu
			}
		}
		flavorId, err = mongodbflexUtils.LoadFlavorId(*cpu, *ram, flavors.Flavors)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load flavor ID: %w", err)
			}
			return req, err
		}
	} else if model.FlavorId != nil {
		err := mongodbflexUtils.ValidateFlavorId(*model.FlavorId, flavors.Flavors)
		if err != nil {
			return req, err
		}
		flavorId = model.FlavorId
	}

	var storages *mongodbflex.ListStoragesResponse
	if model.StorageClass != nil || model.StorageSize != nil {
		validationFlavorId := flavorId
		if validationFlavorId == nil {
			currentInstance, err := apiClient.GetInstanceExecute(ctx, model.ProjectId, model.InstanceId)
			if err != nil {
				return req, fmt.Errorf("get MongoDB Flex instance: %w", err)
			}
			validationFlavorId = currentInstance.Item.Flavor.Id
		}
		storages, err = apiClient.ListStoragesExecute(ctx, model.ProjectId, *validationFlavorId)
		if err != nil {
			return req, fmt.Errorf("get MongoDB Flex storages: %w", err)
		}
		err = mongodbflexUtils.ValidateStorage(model.StorageClass, model.StorageSize, storages, *validationFlavorId)
		if err != nil {
			return req, err
		}
	}

	var payloadAcl *mongodbflex.ACL
	if model.ACL != nil {
		payloadAcl = &mongodbflex.ACL{Items: model.ACL}
	}

	var payloadStorage *mongodbflex.Storage
	if model.StorageClass != nil || model.StorageSize != nil {
		payloadStorage = &mongodbflex.Storage{
			Class: model.StorageClass,
			Size:  model.StorageSize,
		}
	}

	var replicas *int64
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

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *mongodbflex.UpdateInstanceResponse) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal update MongoDBFlex instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal update MongoDBFlex instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Info("%s instance %q\n", operationState, instanceLabel)
		return nil
	}
}
