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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	sqlserverflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/wait"
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
	ACL            *[]string
	BackupSchedule *string
	FlavorId       *string
	CPU            *int64
	RAM            *int64
	Version        *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates an SQLServer Flex instance",
		Long:  "Updates an SQLServer Flex instance.",
		Example: examples.Build(
			examples.NewExample(
				`Update the name of an SQLServer Flex instance with ID "xxx"`,
				"$ stackit sqlserverflex instance update xxx --name my-new-name"),
			examples.NewExample(
				`Update the backup schedule of an SQLServer Flex instance with ID "xxx"`,
				`$ stackit sqlserverflex instance update xxx --backup-schedule "30 0 * * *"`),
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

			instanceLabel, err := sqlserverflexUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
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
				return fmt.Errorf("update SQLServer Flex instance: %w", err)
			}
			instanceId := *resp.Item.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Updating instance")
				_, err = wait.PartialUpdateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SQLServer Flex instance update: %w", err)
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
	acl := flags.FlagToStringSlicePointer(p, cmd, aclFlag)
	backupSchedule := flags.FlagToStringPointer(p, cmd, backupScheduleFlag)
	version := flags.FlagToStringPointer(p, cmd, versionFlag)

	if instanceName == nil && flavorId == nil && cpu == nil && ram == nil && acl == nil &&
		backupSchedule == nil && version == nil {
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
		Version:         version,
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

type SQLServerFlexClient interface {
	PartialUpdateInstance(ctx context.Context, projectId, instanceId string) sqlserverflex.ApiPartialUpdateInstanceRequest
	GetInstanceExecute(ctx context.Context, projectId, instanceId string) (*sqlserverflex.GetInstanceResponse, error)
	ListFlavorsExecute(ctx context.Context, projectId string) (*sqlserverflex.ListFlavorsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, flavorId string) (*sqlserverflex.ListStoragesResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient SQLServerFlexClient) (sqlserverflex.ApiPartialUpdateInstanceRequest, error) {
	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.InstanceId)

	var flavorId *string
	var err error

	flavors, err := apiClient.ListFlavorsExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get SQLServer Flex flavors: %w", err)
	}

	if model.FlavorId == nil && (model.RAM != nil || model.CPU != nil) {
		ram := model.RAM
		cpu := model.CPU
		if model.RAM == nil || model.CPU == nil {
			currentInstance, err := apiClient.GetInstanceExecute(ctx, model.ProjectId, model.InstanceId)
			if err != nil {
				return req, fmt.Errorf("get SQLServer Flex instance: %w", err)
			}
			if model.RAM == nil {
				ram = currentInstance.Item.Flavor.Memory
			}
			if model.CPU == nil {
				cpu = currentInstance.Item.Flavor.Cpu
			}
		}
		flavorId, err = sqlserverflexUtils.LoadFlavorId(*cpu, *ram, flavors.Flavors)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load flavor ID: %w", err)
			}
			return req, err
		}
	} else if model.FlavorId != nil {
		err := sqlserverflexUtils.ValidateFlavorId(*model.FlavorId, flavors.Flavors)
		if err != nil {
			return req, err
		}
		flavorId = model.FlavorId
	}

	var payloadAcl *sqlserverflex.CreateInstancePayloadAcl
	if model.ACL != nil {
		payloadAcl = &sqlserverflex.CreateInstancePayloadAcl{Items: model.ACL}
	}

	req = req.PartialUpdateInstancePayload(sqlserverflex.PartialUpdateInstancePayload{
		Name:           model.InstanceName,
		Acl:            payloadAcl,
		BackupSchedule: model.BackupSchedule,
		FlavorId:       flavorId,
		Version:        model.Version,
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *sqlserverflex.UpdateInstanceResponse) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal update SQLServerFlex instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal update SQLServerFlex instance: %w", err)
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
