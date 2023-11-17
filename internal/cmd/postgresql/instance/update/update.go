package update

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresql/client"
	postgresqlUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresql/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql/wait"
)

const (
	projectIdFlag            = "project-id"
	instanceIdFlag           = "instance-id"
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

type flagModel struct {
	ProjectId  string
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

var Cmd = &cobra.Command{
	Use:     "update",
	Short:   "Updates a PostgreSQL instance",
	Long:    "Updates a PostgreSQL instance",
	Example: `$ stackit postgresql instance update --project-id xxx --instance-id xxx --plan-id xxx --acl xx.xx.xx.xx/xx`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		model, err := parseFlags(cmd)
		if err != nil {
			return err
		}

		// Configure API client
		apiClient, err := client.ConfigureClient(cmd)
		if err != nil {
			return fmt.Errorf("authentication failed, please run \"stackit auth login\" or \"stackit auth activate-service-account\"")
		}

		// Call API
		req, err := buildRequest(ctx, model, apiClient)
		if err != nil {
			return fmt.Errorf("build PostgreSQL instance update request: %w", err)
		}
		err = req.Execute()
		if err != nil {
			return fmt.Errorf("update PostgreSQL instance: %w", err)
		}

		// Wait for async operation
		instanceId := model.InstanceId
		_, err = wait.UpdateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
		if err != nil {
			return fmt.Errorf("wait for PostgreSQL instance update: %w", err)
		}

		fmt.Printf("Updated instance with ID %s\n", instanceId)
		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
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
	cmd.Flags().String(versionFlag, "", "Instance PostgreSQL version")

	cmd.MarkFlagsMutuallyExclusive(planIdFlag, planNameFlag)
	cmd.MarkFlagsMutuallyExclusive(planIdFlag, versionFlag)
	cmd.MarkFlagsRequiredTogether(planNameFlag, versionFlag)

	err := cmd.MarkFlagRequired(instanceIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	planId := utils.FlagToStringPointer(cmd, planIdFlag)
	planName := utils.FlagToStringValue(cmd, planNameFlag)
	version := utils.FlagToStringValue(cmd, versionFlag)

	if planId != nil && (planName != "" || version != "") {
		return nil, fmt.Errorf("please specify either plan-id or plan-name and version but not both")
	}

	return &flagModel{
		ProjectId:            projectId,
		InstanceId:           utils.FlagToStringValue(cmd, instanceIdFlag),
		EnableMonitoring:     utils.FlagToBoolPointer(cmd, enableMonitoringFlag),
		MonitoringInstanceId: utils.FlagToStringPointer(cmd, monitoringInstanceIdFlag),
		Graphite:             utils.FlagToStringPointer(cmd, graphiteFlag),
		MetricsFrequency:     utils.FlagToInt64Pointer(cmd, metricsFrequencyFlag),
		MetricsPrefix:        utils.FlagToStringPointer(cmd, metricsPrefixFlag),
		Plugin:               utils.FlagToStringSlicePointer(cmd, pluginFlag),
		SgwAcl:               utils.FlagToStringSlicePointer(cmd, sgwAclFlag),
		Syslog:               utils.FlagToStringSlicePointer(cmd, syslogFlag),
		PlanId:               planId,
		PlanName:             planName,
		Version:              version,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient postgresqlUtils.PostgreSQLClient) (postgresql.ApiUpdateInstanceRequest, error) {
	req := apiClient.UpdateInstance(ctx, model.ProjectId, model.InstanceId)

	var planId *string
	var err error
	if model.PlanId == nil && model.PlanName != "" && model.Version != "" {
		planId, err = postgresqlUtils.LoadPlanId(ctx, apiClient, model.ProjectId, model.PlanName, model.Version)
		if err != nil {
			return req, fmt.Errorf("load plan ID: %w", err)
		}
	} else {
		planId = model.PlanId
	}

	var sgwAcl *string
	if model.SgwAcl != nil {
		sgwAcl = utils.Ptr(strings.Join(*model.SgwAcl, ","))
	}

	req = req.UpdateInstancePayload(postgresql.UpdateInstancePayload{
		Parameters: &postgresql.InstanceParameters{
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
