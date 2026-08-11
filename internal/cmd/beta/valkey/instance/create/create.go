package create

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/client"
	valkeyUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"
	"github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api/wait"
)

const (
	instanceNameFlag         = "name"
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
	PlanName string
	Version  string

	InstanceName         string
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
		Use:   "create",
		Short: "Creates a Valkey instance",
		Long:  "Creates a Valkey instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Valkey instance with name "my-instance" and specify plan by name and version`,
				"$ stackit beta valkey instance create --name my-instance --plan-name stackit-keyvalue-1.2.10-replica --version 8"),
			examples.NewExample(
				`Create a Valkey instance with name "my-instance" and specify plan by ID`,
				"$ stackit beta valkey instance create --name my-instance --plan-id xxx"),
			examples.NewExample(
				`Create a Valkey instance with name "my-instance" and specify IP range which is allowed to access it`,
				"$ stackit beta valkey instance create --name my-instance --plan-id xxx --acl 1.2.3.0/24"),
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

			prompt := fmt.Sprintf("Are you sure you want to create a Valkey instance for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient.DefaultAPI)
			if err != nil {
				if _, ok := errors.AsType[*cliErr.DSAInvalidPlanError](err); !ok {
					return fmt.Errorf("build Valkey instance creation request: %w", err)
				}
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Valkey instance: %w", err)
			}
			instanceId := resp.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating instance", func() error {
					_, err = wait.CreateInstanceWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, instanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Valkey instance creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, projectLabel, instanceId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
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

	cmd.Flags().Int32(minReplicasToWriteFlag, 0, "Minimum number of replicas that must acknowledge a write for it to be accepted (Valkey only)")
	cmd.Flags().String(replBacklogSizeFlag, "", "Replication backlog size (e.g. \"1mb\") (Valkey only)")

	err := flags.MarkFlagsRequired(cmd, instanceNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	planId := flags.FlagToStringPointer(p, cmd, planIdFlag)
	planName := flags.FlagToStringValue(p, cmd, planNameFlag)
	version := flags.FlagToStringValue(p, cmd, versionFlag)

	if planId == nil && (planName == "" || version == "") {
		return nil, &cliErr.DSAInputPlanError{
			Cmd: cmd,
		}
	}
	if planId != nil && (planName != "" || version != "") {
		return nil, &cliErr.DSAInputPlanError{
			Cmd: cmd,
		}
	}

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		InstanceName:         flags.FlagToStringValue(p, cmd, instanceNameFlag),
		EnableMonitoring:     flags.FlagToBoolPointer(p, cmd, enableMonitoringFlag),
		MonitoringInstanceId: flags.FlagToStringPointer(p, cmd, monitoringInstanceIdFlag),
		Graphite:             flags.FlagToStringPointer(p, cmd, graphiteFlag),
		MetricsFrequency:     flags.FlagToInt32Pointer(p, cmd, metricsFrequencyFlag),
		MetricsPrefix:        flags.FlagToStringPointer(p, cmd, metricsPrefixFlag),
		SgwAcl:               flags.FlagToStringSlicePointer(p, cmd, sgwAclFlag),
		Syslog:               flags.FlagToStringSliceValue(p, cmd, syslogFlag),
		PlanId:               planId,
		PlanName:             planName,
		Version:              version,

		MinReplicasToWrite: flags.FlagToInt32Pointer(p, cmd, minReplicasToWriteFlag),
		ReplBacklogSize:    flags.FlagToStringPointer(p, cmd, replBacklogSizeFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient valkey.DefaultAPI) (valkey.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId, model.Region)

	var planId string
	var err error

	offerings, err := apiClient.ListOfferings(ctx, model.ProjectId, model.Region).Execute()
	if err != nil {
		return req, fmt.Errorf("get Valkey offerings: %w", err)
	}

	if model.PlanId == nil {
		foundPlanId, err := valkeyUtils.LoadPlanId(model.PlanName, model.Version, offerings)
		if err != nil {
			if _, ok := errors.AsType[*cliErr.DSAInvalidPlanError](err); !ok {
				return req, fmt.Errorf("load plan ID: %w", err)
			}
			return req, err
		}
		planId = *foundPlanId
	} else {
		err := valkeyUtils.ValidatePlanId(*model.PlanId, offerings)
		if err != nil {
			return req, err
		}
		planId = *model.PlanId
	}

	var sgwAcl *string
	if model.SgwAcl != nil {
		sgwAcl = new(strings.Join(*model.SgwAcl, ","))
	}

	req = req.CreateInstancePayload(valkey.CreateInstancePayload{
		InstanceName: model.InstanceName,
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

func outputResult(p *print.Printer, model *inputModel, projectLabel, instanceId string, resp *valkey.CreateInstanceResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if model == nil {
			return fmt.Errorf("no model passed")
		}
		if resp == nil {
			return fmt.Errorf("no response defined")
		}

		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, instanceId)
		return nil
	})
}
