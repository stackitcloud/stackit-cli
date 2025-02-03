package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/client"
	sqlserverflexUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/sqlserverflex/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex"
	"github.com/stackitcloud/stackit-sdk-go/services/sqlserverflex/wait"
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
	editionFlag        = "edition"
	retentionDaysFlag  = "retention-days"
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
	Edition        *string
	RetentionDays  *int64
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a SQLServer Flex instance",
		Long:  "Creates a SQLServer Flex instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a SQLServer Flex instance with name "my-instance" and specify flavor by CPU and RAM. Other parameters are set to default values`,
				`$ stackit beta sqlserverflex instance create --name my-instance --cpu 1 --ram 4`),
			examples.NewExample(
				`Create a SQLServer Flex instance with name "my-instance" and specify flavor by ID. Other parameters are set to default values`,
				`$ stackit beta sqlserverflex instance create --name my-instance --flavor-id xxx`),
			examples.NewExample(
				`Create a SQLServer Flex instance with name "my-instance", specify flavor by CPU and RAM, set storage size to 20 GB, and restrict access to a specific range of IP addresses. Other parameters are set to default values`,
				`$ stackit beta sqlserverflex instance create --name my-instance --cpu 1 --ram 4 --storage-size 20  --acl 1.2.3.0/24`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a SQLServer Flex instance for project %q?", projectLabel)
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
				return fmt.Errorf("create SQLServer Flex instance: %w", err)
			}
			instanceId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Creating instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SQLServer Flex instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc.")
	cmd.Flags().String(backupScheduleFlag, "", "Backup schedule")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor")
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs")
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB)")
	cmd.Flags().Int64(storageSizeFlag, 0, "Storage size (in GB)")
	cmd.Flags().String(storageClassFlag, "", "Storage class")
	cmd.Flags().String(versionFlag, "", "SQLServer version")
	cmd.Flags().String(editionFlag, "", "Edition of the SQLServer instance")
	cmd.Flags().Int64(retentionDaysFlag, 0, "The days for how long the backup files should be stored before being cleaned up")

	err := flags.MarkFlagsRequired(cmd, instanceNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
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
		InstanceName:    flags.FlagToStringPointer(p, cmd, instanceNameFlag),
		ACL:             flags.FlagToStringSlicePointer(p, cmd, aclFlag),
		BackupSchedule:  flags.FlagToStringPointer(p, cmd, backupScheduleFlag),
		FlavorId:        flavorId,
		CPU:             cpu,
		RAM:             ram,
		StorageClass:    flags.FlagToStringPointer(p, cmd, storageClassFlag),
		StorageSize:     flags.FlagToInt64Pointer(p, cmd, storageSizeFlag),
		Version:         flags.FlagToStringPointer(p, cmd, versionFlag),
		Edition:         flags.FlagToStringPointer(p, cmd, editionFlag),
		RetentionDays:   flags.FlagToInt64Pointer(p, cmd, retentionDaysFlag),
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

type sqlServerFlexClient interface {
	CreateInstance(ctx context.Context, projectId string) sqlserverflex.ApiCreateInstanceRequest
	ListFlavorsExecute(ctx context.Context, projectId string) (*sqlserverflex.ListFlavorsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, flavorId string) (*sqlserverflex.ListStoragesResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient sqlServerFlexClient) (sqlserverflex.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId)

	var flavorId *string
	var err error

	flavors, err := apiClient.ListFlavorsExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get SQLServer Flex flavors: %w", err)
	}

	if model.FlavorId == nil {
		flavorId, err = sqlserverflexUtils.LoadFlavorId(*model.CPU, *model.RAM, flavors.Flavors)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load flavor ID: %w", err)
			}
			return req, err
		}
	} else {
		err := sqlserverflexUtils.ValidateFlavorId(*model.FlavorId, flavors.Flavors)
		if err != nil {
			return req, err
		}
		flavorId = model.FlavorId
	}

	storages, err := apiClient.ListStoragesExecute(ctx, model.ProjectId, *flavorId)
	if err != nil {
		return req, fmt.Errorf("get SQLServer Flex storages: %w", err)
	}
	err = sqlserverflexUtils.ValidateStorage(model.StorageClass, model.StorageSize, storages, *flavorId)
	if err != nil {
		return req, err
	}

	var retentionDays *string
	if model.RetentionDays != nil {
		retentionDays = utils.Ptr(fmt.Sprintf("%d", *model.RetentionDays))
	}

	req = req.CreateInstancePayload(sqlserverflex.CreateInstancePayload{
		Name:           model.InstanceName,
		Acl:            &sqlserverflex.CreateInstancePayloadAcl{Items: model.ACL},
		BackupSchedule: model.BackupSchedule,
		FlavorId:       flavorId,
		Storage: &sqlserverflex.CreateInstancePayloadStorage{
			Class: model.StorageClass,
			Size:  model.StorageSize,
		},
		Version: model.Version,
		Options: &sqlserverflex.CreateInstancePayloadOptions{
			Edition:       model.Edition,
			RetentionDays: retentionDays,
		},
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *sqlserverflex.CreateInstanceResponse) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SQLServerFlex instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal SQLServerFlex instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, utils.PtrString(resp.Id))
		return nil
	}
}
