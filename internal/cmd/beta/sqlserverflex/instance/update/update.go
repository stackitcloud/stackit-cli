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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	sqlserverflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	sqlserverflex "github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/v3api/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	instanceNameFlag   = "name"
	aclFlag            = "acl"
	backupScheduleFlag = "backup-schedule"
	flavorIdFlag       = "flavor-id"
	cpuFlag            = "cpu"
	ramFlag            = "ram"
	versionFlag        = "version"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId     string
	InstanceName   *string
	ACL            []string
	BackupSchedule *string
	FlavorId       *string
	CPU            *int64
	RAM            *int64
	Version        *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a SQLServer Flex instance",
		Long:  "Updates a SQLServer Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`Update the name of a SQLServer Flex instance with ID "xxx"`,
				"$ stackit beta sqlserverflex instance update xxx --name my-new-name"),
			examples.NewExample(
				`Update the backup schedule of a SQLServer Flex instance with ID "xxx"`,
				`$ stackit beta sqlserverflex instance update xxx --backup-schedule "30 0 * * *"`),
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

			instanceLabel, err := sqlserverflexUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.InstanceId, model.Region)
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
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update SQLServer Flex instance: %w", err)
			}

			var instance *sqlserverflex.GetInstanceResponse
			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating instance", func() error {
					instance, err = wait.UpdateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for SQLServer Flex instance update: %w", err)
				}
			}
			if instance == nil {
				instance, err = apiClient.DefaultAPI.GetInstance(ctx, model.ProjectId, model.Region, model.InstanceId).Execute()
				if err != nil {
					return fmt.Errorf("get SQLServer Flex instance: %w", err)
				}
			}

			return outputResult(params.Printer, model, instanceLabel, instance)
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
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs")
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB)")
	cmd.Flags().String(versionFlag, "", "Version")
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
	version := flags.FlagToStringPointer(p, cmd, versionFlag)

	if instanceName == nil && flavorId == nil && cpu == nil && ram == nil && acl == nil &&
		backupSchedule == nil && version == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	if flavorId != nil && (cpu != nil || ram != nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd:     cmd,
			Service: sqlserverflexUtils.ServiceCmd,
			Args:    inputArgs,
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
		Version:         version,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient sqlserverflex.DefaultAPI) (sqlserverflex.ApiPartialUpdateInstanceRequest, error) {
	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.Region, model.InstanceId)

	var flavorId *string
	var err error

	flavors, err := apiClient.ListFlavors(ctx, model.ProjectId, model.Region).Execute()
	if err != nil {
		return req, fmt.Errorf("get SQLServer Flex flavors: %w", err)
	}

	if model.FlavorId == nil && (model.RAM != nil || model.CPU != nil) {
		ram := model.RAM
		cpu := model.CPU
		if model.RAM == nil || model.CPU == nil {
			currentInstance, err := apiClient.GetInstance(ctx, model.ProjectId, model.Region, model.InstanceId).Execute()
			if err != nil {
				return req, fmt.Errorf("get SQLServer Flex instance: %w", err)
			}
			var currentFlavor *sqlserverflex.ListFlavors
			for _, flavor := range flavors.Flavors {
				if flavor.Id == currentInstance.FlavorId {
					currentFlavor = &flavor
				}
			}
			if currentFlavor == nil {
				return req, fmt.Errorf("can't find flavor %s in flavors list", currentInstance.FlavorId)
			}
			if model.RAM == nil {
				ram = &currentFlavor.Memory
			}
			if model.CPU == nil {
				cpu = &currentFlavor.Cpu
			}
		}
		loadedId, err := sqlserverflexUtils.LoadFlavorId(*cpu, *ram, flavors.Flavors)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load flavor ID: %w", err)
			}
			return req, err
		}
		flavorId = &loadedId
	} else if model.FlavorId != nil {
		err := sqlserverflexUtils.ValidateFlavorId(*model.FlavorId, flavors.Flavors)
		if err != nil {
			return req, err
		}
		flavorId = model.FlavorId
	}

	var network *sqlserverflex.PartialUpdateInstancePayloadNetwork
	if model.ACL != nil {
		network = &sqlserverflex.PartialUpdateInstancePayloadNetwork{
			Acl: model.ACL,
		}
	}

	req = req.PartialUpdateInstancePayload(sqlserverflex.PartialUpdateInstancePayload{
		Name:           model.InstanceName,
		Network:        network,
		BackupSchedule: model.BackupSchedule,
		FlavorId:       flavorId,
		Version:        (*sqlserverflex.InstanceVersionOpt)(model.Version),
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *sqlserverflex.GetInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("instance response is empty")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Info("%s instance %q\n", operationState, instanceLabel)
		return nil
	})
}
