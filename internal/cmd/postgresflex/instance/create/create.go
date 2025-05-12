package create

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresflex/wait"
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

	defaultBackupSchedule = "0 0 * * *"
	defaultStorageClass   = "premium-perf2-stackit"
	defaultStorageSize    = 10
	defaultType           = "Replica"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a PostgreSQL Flex instance",
		Long:  "Creates a PostgreSQL Flex instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a PostgreSQL Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access) and specify flavor by CPU and RAM. Other parameters are set to default values`,
				`$ stackit postgresflex instance create --name my-instance --cpu 2 --ram 4 --acl 0.0.0.0/0`),
			examples.NewExample(
				`Create a PostgreSQL Flex instance with name "my-instance", ACL 0.0.0.0/0 (open access) and specify flavor by ID. Other parameters are set to default values`,
				`$ stackit postgresflex instance create --name my-instance --flavor-id xxx --acl 0.0.0.0/0`),
			examples.NewExample(
				`Create a PostgreSQL Flex instance with name "my-instance", allow access to a specific range of IP addresses, specify flavor by CPU and RAM and set storage size to 20 GB. Other parameters are set to default values`,
				`$ stackit postgresflex instance create --name my-instance --cpu 2 --ram 4 --acl 1.2.3.0/24 --storage-size 20`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			model, err := parseInput(params.Printer, cmd)
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a PostgreSQL Flex instance for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Fill in version, if needed
			if model.Version == nil {
				version, err := postgresflexUtils.GetLatestPostgreSQLVersion(ctx, apiClient, model.ProjectId, model.Region)
				if err != nil {
					return fmt.Errorf("get latest PostgreSQL version: %w", err)
				}
				model.Version = &version
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create PostgreSQL Flex instance: %w", err)
			}
			instanceId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, model.ProjectId, model.Region, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for PostgreSQL Flex instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, instanceId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	typeFlagOptions := postgresflexUtils.AvailableInstanceTypes()

	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "The access control list (ACL). Must contain at least one valid subnet, for instance '0.0.0.0/0' for open access (discouraged), '1.2.3.0/24 for a public IP range of an organization, '1.2.3.4/32' for a single IP range, etc.")
	cmd.Flags().String(backupScheduleFlag, defaultBackupSchedule, "Backup schedule")
	cmd.Flags().String(flavorIdFlag, "", "ID of the flavor")
	cmd.Flags().Int64(cpuFlag, 0, "Number of CPUs")
	cmd.Flags().Int64(ramFlag, 0, "Amount of RAM (in GB)")
	cmd.Flags().String(storageClassFlag, defaultStorageClass, "Storage class")
	cmd.Flags().Int64(storageSizeFlag, defaultStorageSize, "Storage size (in GB)")
	cmd.Flags().String(versionFlag, "", "PostgreSQL version. Defaults to the latest version available")
	cmd.Flags().Var(flags.EnumFlag(false, defaultType, typeFlagOptions...), typeFlag, fmt.Sprintf("Instance type, one of %q", typeFlagOptions))

	err := flags.MarkFlagsRequired(cmd, instanceNameFlag, aclFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	storageSize := flags.FlagWithDefaultToInt64Value(p, cmd, storageSizeFlag)

	flavorId := flags.FlagToStringPointer(p, cmd, flavorIdFlag)
	cpu := flags.FlagToInt64Pointer(p, cmd, cpuFlag)
	ram := flags.FlagToInt64Pointer(p, cmd, ramFlag)

	if flavorId == nil && (cpu == nil || ram == nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd: cmd,
		}
	}
	if flavorId != nil && (cpu != nil || ram != nil) {
		return nil, &cliErr.DatabaseInputFlavorError{
			Cmd: cmd,
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceName:    flags.FlagToStringPointer(p, cmd, instanceNameFlag),
		ACL:             flags.FlagToStringSlicePointer(p, cmd, aclFlag),
		BackupSchedule:  utils.Ptr(flags.FlagWithDefaultToStringValue(p, cmd, backupScheduleFlag)),
		FlavorId:        flavorId,
		CPU:             cpu,
		RAM:             ram,
		StorageClass:    utils.Ptr(flags.FlagWithDefaultToStringValue(p, cmd, storageClassFlag)),
		StorageSize:     &storageSize,
		Version:         flags.FlagToStringPointer(p, cmd, versionFlag),
		Type:            utils.Ptr(flags.FlagWithDefaultToStringValue(p, cmd, typeFlag)),
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

type PostgreSQLFlexClient interface {
	CreateInstance(ctx context.Context, projectId, region string) postgresflex.ApiCreateInstanceRequest
	ListFlavorsExecute(ctx context.Context, projectId, region string) (*postgresflex.ListFlavorsResponse, error)
	ListStoragesExecute(ctx context.Context, projectId, region, flavorId string) (*postgresflex.ListStoragesResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient PostgreSQLFlexClient) (postgresflex.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId, model.Region)

	var flavorId *string
	var err error

	flavors, err := apiClient.ListFlavorsExecute(ctx, model.ProjectId, model.Region)
	if err != nil {
		return req, fmt.Errorf("get PostgreSQL Flex flavors: %w", err)
	}

	if model.FlavorId == nil {
		flavorId, err = postgresflexUtils.LoadFlavorId(*model.CPU, *model.RAM, flavors.Flavors)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load flavor ID: %w", err)
			}
			return req, err
		}
	} else {
		err := postgresflexUtils.ValidateFlavorId(*model.FlavorId, flavors.Flavors)
		if err != nil {
			return req, err
		}
		flavorId = model.FlavorId
	}

	storages, err := apiClient.ListStoragesExecute(ctx, model.ProjectId, model.Region, *flavorId)
	if err != nil {
		return req, fmt.Errorf("get PostgreSQL Flex storages: %w", err)
	}
	err = postgresflexUtils.ValidateStorage(model.StorageClass, model.StorageSize, storages, *flavorId)
	if err != nil {
		return req, err
	}

	replicas, err := postgresflexUtils.GetInstanceReplicas(*model.Type)
	if err != nil {
		return req, fmt.Errorf("get PostgreSQL Flex instance type: %w", err)
	}

	req = req.CreateInstancePayload(postgresflex.CreateInstancePayload{
		Name:           model.InstanceName,
		Acl:            &postgresflex.ACL{Items: model.ACL},
		BackupSchedule: model.BackupSchedule,
		FlavorId:       flavorId,
		Replicas:       &replicas,
		Storage: &postgresflex.Storage{
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

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel, instanceId string, resp *postgresflex.CreateInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("no response passed")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgresFlex instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal PostgresFlex instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, instanceId)
		return nil
	}
}
