package update

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/redis/client"
	redisUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/redis/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/redis"
	"github.com/stackitcloud/stackit-sdk-go/services/redis/wait"
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
	SgwAcl               *[]string
	Syslog               *[]string
	PlanId               *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a Redis instance",
		Long:  "Updates a Redis instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the plan of a Redis instance with ID "xxx"`,
				"$ stackit redis instance update xxx --plan-id yyy"),
			examples.NewExample(
				`Update the range of IPs allowed to access a Redis instance with ID "xxx"`,
				"$ stackit redis instance update xxx --acl 1.2.3.0/24"),
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

			instanceLabel, err := redisUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
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
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
				if !errors.As(err, &dsaInvalidPlanError) {
					return fmt.Errorf("build Redis instance update request: %w", err)
				}
				return err
			}
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Redis instance: %w", err)
			}
			instanceId := model.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Updating instance")
				_, err = wait.PartialUpdateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Redis instance update: %w", err)
				}
				s.Stop()
			}

			operationState := "Updated"
			if model.Async {
				operationState = "Triggered update of"
			}
			params.Printer.Info("%s instance %q\n", operationState, instanceLabel)
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
	cmd.Flags().Var(flags.CIDRSliceFlag(), sgwAclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")
	cmd.Flags().StringSlice(syslogFlag, []string{}, "Syslog")
	cmd.Flags().Var(flags.UUIDFlag(), planIdFlag, "Plan ID")
	cmd.Flags().String(planNameFlag, "", "Plan name")
	cmd.Flags().String(versionFlag, "", "Instance Redis version")
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
	metricsFrequency := flags.FlagToInt64Pointer(p, cmd, metricsFrequencyFlag)
	metricsPrefix := flags.FlagToStringPointer(p, cmd, metricsPrefixFlag)
	sgwAcl := flags.FlagToStringSlicePointer(p, cmd, sgwAclFlag)
	syslog := flags.FlagToStringSlicePointer(p, cmd, syslogFlag)
	planId := flags.FlagToStringPointer(p, cmd, planIdFlag)
	planName := flags.FlagToStringValue(p, cmd, planNameFlag)
	version := flags.FlagToStringValue(p, cmd, versionFlag)

	if planId != nil && (planName != "" || version != "") {
		return nil, &cliErr.DSAInputPlanError{
			Cmd:  cmd,
			Args: inputArgs,
		}
	}

	if enableMonitoring == nil && monitoringInstanceId == nil && graphite == nil &&
		metricsFrequency == nil && metricsPrefix == nil &&
		sgwAcl == nil && syslog == nil && planId == nil &&
		planName == "" && version == "" {
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
	}

	p.DebugInputModel(model)
	return &model, nil
}

type redisClient interface {
	PartialUpdateInstance(ctx context.Context, projectId, instanceId string) redis.ApiPartialUpdateInstanceRequest
	ListOfferingsExecute(ctx context.Context, projectId string) (*redis.ListOfferingsResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient redisClient) (redis.ApiPartialUpdateInstanceRequest, error) {
	req := apiClient.PartialUpdateInstance(ctx, model.ProjectId, model.InstanceId)

	var planId *string
	var err error

	offerings, err := apiClient.ListOfferingsExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get Redis offerings: %w", err)
	}

	if model.PlanId == nil && model.PlanName != "" && model.Version != "" {
		planId, err = redisUtils.LoadPlanId(model.PlanName, model.Version, offerings)
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
			err := redisUtils.ValidatePlanId(*model.PlanId, offerings)
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

	req = req.PartialUpdateInstancePayload(redis.PartialUpdateInstancePayload{
		Parameters: &redis.InstanceParameters{
			EnableMonitoring:     model.EnableMonitoring,
			Graphite:             model.Graphite,
			MonitoringInstanceId: model.MonitoringInstanceId,
			MetricsFrequency:     model.MetricsFrequency,
			MetricsPrefix:        model.MetricsPrefix,
			SgwAcl:               sgwAcl,
			Syslog:               model.Syslog,
		},
		PlanId: planId,
	})
	return req, nil
}
