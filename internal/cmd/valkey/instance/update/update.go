package update

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"
	"github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/client"
	valkeyUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	instanceIdArg = "INSTANCE_ID"

	enableMonitoringFlag     = "enable-monitoring"
	graphiteFlag             = "graphite"
	metricsFrequencyFlag     = "metrics-frequency"
	metricsPrefixFlag        = "metrics-prefix"
	monitoringInstanceIdFlag = "monitoring-instance-id"
	sgwAclFlag               = "acl"
	syslogFlag               = "syslog"
	planIdFlag               = "plan-id"
	planNameFlag             = "plan-name"
	versionFlag              = "version"

	minReplicasToWriteFlag = "min-replicas-to-write"
	replBacklogSizeFlag    = "repl-backlog-size"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	PlanName   string
	Version    string

	EnableMonitoring     *bool
	Graphite             *string
	MetricsFrequency     *int32
	MetricsPrefix        *string
	MonitoringInstanceId *string
	SgwAcl               *[]string
	Syslog               []string
	PlanId               *string

	MinReplicasToWrite *int32
	ReplBacklogSize    *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a Valkey instance",
		Long:  "Updates a Valkey instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the plan of a Valkey instance with ID "xxx" by plan ID`,
				"$ stackit valkey instance update xxx --plan-id yyy"),
			examples.NewExample(
				`Update the plan of a Valkey instance with ID "xxx" by name and version`,
				"$ stackit valkey instance update xxx --plan-name stackit-keyvalue-1.2.10-replica --version 8"),
			examples.NewExample(
				`Update the range of IPs allowed to access a Valkey instance with ID "xxx"`,
				"$ stackit valkey instance update xxx --acl 1.2.3.0/24"),
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

			instanceLabel, err := valkeyUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.InstanceId, model.Region)
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
				if _, ok := errors.AsType[*cliErr.DSAInvalidPlanError](err); !ok {
					return fmt.Errorf("build Valkey instance update request: %w", err)
				}
				return err
			}
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Valkey instance: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating instance", func() error {
					_, err = wait.PartialUpdateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Valkey instance update: %w", err)
				}
			}

			return outputResult(params.Printer, model, instanceLabel)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(enableMonitoringFlag, false, "Enable monitoring")
	cmd.Flags().String(graphiteFlag, "", "Graphite host")
	cmd.Flags().Int32(metricsFrequencyFlag, 0, "Metrics frequency in seconds")
	cmd.Flags().String(metricsPrefixFlag, "", "Metrics prefix")
	cmd.Flags().Var(flags.UUIDFlag(), monitoringInstanceIdFlag, "Monitoring instance ID")
	cmd.Flags().Var(flags.CIDRSliceFlag(), sgwAclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")
	cmd.Flags().StringSlice(syslogFlag, []string{}, "Syslog")
	cmd.Flags().Var(flags.UUIDFlag(), planIdFlag, "Plan ID")
	cmd.Flags().String(planNameFlag, "", "Plan name")
	cmd.Flags().String(versionFlag, "", "Instance Valkey version")

	cmd.Flags().Int32(minReplicasToWriteFlag, 0, "Minimum number of replicas that must acknowledge a write for it to be accepted")
	cmd.Flags().String(replBacklogSizeFlag, "", "Replication backlog size (e.g. \"1mb\")")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	enableMonitoring := flags.FlagToBoolPointer(p, cmd, enableMonitoringFlag)
	monitoringInstanceId := flags.FlagToStringPointer(p, cmd, monitoringInstanceIdFlag)
	graphite := flags.FlagToStringPointer(p, cmd, graphiteFlag)
	metricsFrequency := flags.FlagToInt32Pointer(p, cmd, metricsFrequencyFlag)
	metricsPrefix := flags.FlagToStringPointer(p, cmd, metricsPrefixFlag)
	sgwAcl := flags.FlagToStringSlicePointer(p, cmd, sgwAclFlag)
	syslog := flags.FlagToStringSliceValue(p, cmd, syslogFlag)
	planId := flags.FlagToStringPointer(p, cmd, planIdFlag)
	planName := flags.FlagToStringValue(p, cmd, planNameFlag)
	version := flags.FlagToStringValue(p, cmd, versionFlag)
	minReplicasToWrite := flags.FlagToInt32Pointer(p, cmd, minReplicasToWriteFlag)
	replBacklogSize := flags.FlagToStringPointer(p, cmd, replBacklogSizeFlag)

	if planId != nil && (planName != "" || version != "") {
		return nil, &cliErr.DSAInputPlanError{
			Cmd: cmd,
		}
	}
	if planId == nil && (planName == "") != (version == "") {
		return nil, &cliErr.DSAInputPlanError{
			Cmd: cmd,
		}
	}

	if enableMonitoring == nil && monitoringInstanceId == nil &&
		graphite == nil && metricsFrequency == nil && metricsPrefix == nil &&
		sgwAcl == nil && len(syslog) == 0 && planId == nil && planName == "" && version == "" &&
		minReplicasToWrite == nil && replBacklogSize == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		InstanceId:           instanceId,
		EnableMonitoring:     enableMonitoring,
		MonitoringInstanceId: monitoringInstanceId,
		Graphite:             graphite,
		MetricsFrequency:     metricsFrequency,
		MetricsPrefix:        metricsPrefix,
		SgwAcl:               sgwAcl,
		Syslog:               syslog,
		PlanId:               planId,
		PlanName:             planName,
		Version:              version,
		MinReplicasToWrite:   minReplicasToWrite,
		ReplBacklogSize:      replBacklogSize,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient valkey.DefaultAPI) (valkey.ApiPartialUpdateInstanceRequest, error) {
	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.Region, model.InstanceId)

	var planId *string

	// Only call ListOfferings when plan selection is requested
	if model.PlanId != nil || model.PlanName != "" {
		offerings, err := apiClient.ListOfferings(ctx, model.ProjectId, model.Region).Execute()
		if err != nil {
			return req, fmt.Errorf("get Valkey offerings: %w", err)
		}

		if model.PlanId != nil {
			if err := valkeyUtils.ValidatePlanId(*model.PlanId, offerings); err != nil {
				return req, err
			}
			planId = model.PlanId
		} else {
			foundPlanId, err := valkeyUtils.LoadPlanId(model.PlanName, model.Version, offerings)
			if err != nil {
				if _, ok := errors.AsType[*cliErr.DSAInvalidPlanError](err); !ok {
					return req, fmt.Errorf("load plan ID: %w", err)
				}
				return req, err
			}
			planId = foundPlanId
		}
	}

	var sgwAcl *string
	if model.SgwAcl != nil {
		sgwAcl = new(strings.Join(*model.SgwAcl, ","))
	}

	req = req.PartialUpdateInstancePayload(valkey.PartialUpdateInstancePayload{
		Parameters: &valkey.InstanceParameters{
			EnableMonitoring:     model.EnableMonitoring,
			Graphite:             model.Graphite,
			MonitoringInstanceId: model.MonitoringInstanceId,
			MetricsFrequency:     model.MetricsFrequency,
			MetricsPrefix:        model.MetricsPrefix,
			SgwAcl:               sgwAcl,
			Syslog:               model.Syslog,
			MinReplicasToWrite:   model.MinReplicasToWrite,
			ReplBacklogSize:      model.ReplBacklogSize,
		},
		PlanId: planId,
	})
	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string) error {
	operationState := "Updated"
	if model.Async {
		operationState = "Triggered update of"
	}
	p.Info("%s instance %q\n", operationState, instanceLabel)
	return nil
}
