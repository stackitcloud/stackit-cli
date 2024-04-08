package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logme/client"
	logmeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/logme/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/logme"
	"github.com/stackitcloud/stackit-sdk-go/services/logme/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	instanceNameFlag         = "name"
	enableMonitoringFlag     = "enable-monitoring"
	graphiteFlag             = "graphite"
	metricsFrequencyFlag     = "metrics-frequency"
	metricsPrefixFlag        = "metrics-prefix"
	monitoringInstanceIdFlag = "monitoring-instance-id"
	pluginFlag               = "plugin"
	sgwAclFlag               = "acl"
	syslogFlag               = "syslog"
	planIdFlag               = "plan-id"
	planNameFlag             = "plan-name"
	versionFlag              = "version"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	PlanName   string
	Version    string

	EnableMonitoring     *bool
	Graphite             *string
	MetricsFrequency     *int64
	MetricsPrefix        *string
	MonitoringInstanceId *string
	Plugin               *[]string
	SgwAcl               *[]string
	Syslog               *[]string
	PlanId               *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a LogMe instance",
		Long:  "Updates a LogMe instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the plan of a LogMe instance with ID "xxx"`,
				"$ stackit logme instance update xxx --plan-id yyy"),
			examples.NewExample(
				`Update the range of IPs allowed to access a LogMe instance with ID "xxx"`,
				"$ stackit logme instance update xxx --acl 1.2.3.0/24"),
		),
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

			instanceLabel, err := logmeUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update instance %q?", instanceLabel)
				err = p.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
				if !errors.As(err, &dsaInvalidPlanError) {
					return fmt.Errorf("build LogMe instance update request: %w", err)
				}
				return err
			}
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update LogMe instance: %w", err)
			}
			instanceId := model.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Updating instance")
				_, err = wait.PartialUpdateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for LogMe instance update: %w", err)
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
	cmd.Flags().Bool(enableMonitoringFlag, false, "Enable monitoring")
	cmd.Flags().String(graphiteFlag, "", "Graphite host")
	cmd.Flags().Int64(metricsFrequencyFlag, 0, "Metrics frequency")
	cmd.Flags().String(metricsPrefixFlag, "", "Metrics prefix")
	cmd.Flags().Var(flags.UUIDFlag(), monitoringInstanceIdFlag, "Monitoring instance ID")
	cmd.Flags().StringSlice(pluginFlag, []string{}, "Plugin")
	cmd.Flags().Var(flags.CIDRSliceFlag(), sgwAclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")
	cmd.Flags().StringSlice(syslogFlag, []string{}, "Syslog")
	cmd.Flags().Var(flags.UUIDFlag(), planIdFlag, "Plan ID")
	cmd.Flags().String(planNameFlag, "", "Plan name")
	cmd.Flags().String(versionFlag, "", "Instance LogMe version")
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	enableMonitoring := flags.FlagToBoolPointer(cmd, enableMonitoringFlag)
	monitoringInstanceId := flags.FlagToStringPointer(cmd, monitoringInstanceIdFlag)
	graphite := flags.FlagToStringPointer(cmd, graphiteFlag)
	metricsFrequency := flags.FlagToInt64Pointer(cmd, metricsFrequencyFlag)
	metricsPrefix := flags.FlagToStringPointer(cmd, metricsPrefixFlag)
	plugin := flags.FlagToStringSlicePointer(cmd, pluginFlag)
	sgwAcl := flags.FlagToStringSlicePointer(cmd, sgwAclFlag)
	syslog := flags.FlagToStringSlicePointer(cmd, syslogFlag)
	planId := flags.FlagToStringPointer(cmd, planIdFlag)
	planName := flags.FlagToStringValue(cmd, planNameFlag)
	version := flags.FlagToStringValue(cmd, versionFlag)

	if planId != nil && (planName != "" || version != "") {
		return nil, &cliErr.DSAInputPlanError{
			Cmd:  cmd,
			Args: inputArgs,
		}
	}

	if enableMonitoring == nil && monitoringInstanceId == nil && graphite == nil &&
		metricsFrequency == nil && metricsPrefix == nil && plugin == nil &&
		sgwAcl == nil && syslog == nil && planId == nil &&
		planName == "" && version == "" {
		return nil, &cliErr.EmptyUpdateError{}
	}

	return &inputModel{
		GlobalFlagModel:      globalFlags,
		InstanceId:           instanceId,
		EnableMonitoring:     enableMonitoring,
		MonitoringInstanceId: monitoringInstanceId,
		Graphite:             graphite,
		MetricsFrequency:     metricsFrequency,
		MetricsPrefix:        metricsPrefix,
		Plugin:               plugin,
		SgwAcl:               sgwAcl,
		Syslog:               syslog,
		PlanId:               planId,
		PlanName:             planName,
		Version:              version,
	}, nil
}

type logMeClient interface {
	PartialUpdateInstance(ctx context.Context, projectId, instanceId string) logme.ApiPartialUpdateInstanceRequest
	ListOfferingsExecute(ctx context.Context, projectId string) (*logme.ListOfferingsResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient logMeClient) (logme.ApiPartialUpdateInstanceRequest, error) {
	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.InstanceId)

	var planId *string
	var err error

	offerings, err := apiClient.ListOfferingsExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get LogMe offerings: %w", err)
	}

	if model.PlanId == nil && model.PlanName != "" && model.Version != "" {
		planId, err = logmeUtils.LoadPlanId(model.PlanName, model.Version, offerings)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load plan ID: %w", err)
			}
			return req, err
		}
	} else {
		// planId is not required for update operation
		if model.PlanId != nil {
			err := logmeUtils.ValidatePlanId(*model.PlanId, offerings)
			if err != nil {
				return req, err
			}
		}
		planId = model.PlanId
	}

	var sgwAcl *string
	if model.SgwAcl != nil {
		sgwAcl = utils.Ptr(strings.Join(*model.SgwAcl, ","))
	}

	req = req.PartialUpdateInstancePayload(logme.PartialUpdateInstancePayload{
		Parameters: &logme.InstanceParameters{
			EnableMonitoring:     model.EnableMonitoring,
			Graphite:             model.Graphite,
			MonitoringInstanceId: model.MonitoringInstanceId,
			MetricsFrequency:     model.MetricsFrequency,
			MetricsPrefix:        model.MetricsPrefix,
			Plugins:              model.Plugin,
			SgwAcl:               sgwAcl,
			Syslog:               model.Syslog,
		},
		PlanId: planId,
	})
	return req, nil
}
