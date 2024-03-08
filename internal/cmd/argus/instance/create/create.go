package create

import (
	"context"
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/client"
	argusUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/argus/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/argus"
	"github.com/stackitcloud/stackit-sdk-go/services/argus/wait"
)

const (
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
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PlanName string

	InstanceName         *string
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates an Argus instance",
		Long:  "Creates an Argus instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an Argus instance with name "my-instance" and specify plan by name`,
				"$ stackit argus instance create --name my-instance --plan-name Monitoring-Starter-EU01"),
			examples.NewExample(
				`Create an Argus instance with name "my-instance" and specify plan by ID`,
				"$ stackit argus instance create --name my-instance --plan-id xxx"),
			examples.NewExample(
				`Create an Argus instance with name "my-instance" and specify IP range which is allowed to access it`,
				"$ stackit argus instance create --name my-instance --plan-id xxx --acl 1.2.3.0/24"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create an Argus instance for project %q?", projectLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
				if !errors.As(err, &dsaInvalidPlanError) {
					return fmt.Errorf("build Argus instance creation request: %w", err)
				}
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Argus instance: %w", err)
			}
			instanceId := *resp.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(cmd)
				s.Start("Creating instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for Argus instance creation: %w", err)
				}
				s.Stop()
			}

			operationState := "Created"
			if model.Async {
				operationState = "Triggered creation of"
			}
			cmd.Printf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, instanceId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
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

	err := flags.MarkFlagsRequired(cmd, instanceNameFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	planId := flags.FlagToStringPointer(cmd, planIdFlag)
	planName := flags.FlagToStringValue(cmd, planNameFlag)

	if planId == nil && (planName == "") {
		return nil, &cliErr.DSAInputPlanError{
			Cmd: cmd,
		}
	}
	if planId != nil && (planName != "" ){
		return nil, &cliErr.DSAInputPlanError{
			Cmd: cmd,
		}
	}

	return &inputModel{
		GlobalFlagModel:      globalFlags,
		InstanceName:         flags.FlagToStringPointer(cmd, instanceNameFlag),
		EnableMonitoring:     flags.FlagToBoolPointer(cmd, enableMonitoringFlag),
		MonitoringInstanceId: flags.FlagToStringPointer(cmd, monitoringInstanceIdFlag),
		Graphite:             flags.FlagToStringPointer(cmd, graphiteFlag),
		MetricsFrequency:     flags.FlagToInt64Pointer(cmd, metricsFrequencyFlag),
		MetricsPrefix:        flags.FlagToStringPointer(cmd, metricsPrefixFlag),
		Plugin:               flags.FlagToStringSlicePointer(cmd, pluginFlag),
		SgwAcl:               flags.FlagToStringSlicePointer(cmd, sgwAclFlag),
		Syslog:               flags.FlagToStringSlicePointer(cmd, syslogFlag),
		PlanId:               planId,
		PlanName:             planName,
	}, nil
}

type argusClient interface {
	CreateInstance(ctx context.Context, projectId string) argus.ApiCreateInstanceRequest
	ListPlansExecute(ctx context.Context, projectId string) (*argus.PlansResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient argusClient) (argus.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId)

	var planId *string
	var err error

	plans, err := apiClient.ListPlansExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get Argus offerings: %w", err)
	}

	if model.PlanId == nil {
		planId, err = argusUtils.LoadPlanId(model.PlanName, plans)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load plan ID: %w", err)
			}
			return req, err
		}
	} else {
		err := argusUtils.ValidatePlanId(*model.PlanId, plans)
		if err != nil {
			return req, err
		}
		planId = model.PlanId
	}

	req = req.CreateInstancePayload(argus.CreateInstancePayload{
		Name: model.InstanceName,
		PlanId: planId,
	})
	return req, nil
}
