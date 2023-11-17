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

		// Get instance
		instance, err := apiClient.GetInstanceExecute(ctx, model.ProjectId, model.InstanceId)
		if err != nil {
			return fmt.Errorf("get PostgreSQL instance: %w", err)
		}

		// Call API
		req, err := buildRequest(ctx, instance, model, apiClient)
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
	cmd.Flags().StringP(instanceIdFlag, "i", "", "Instance ID")
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

func buildRequest(ctx context.Context, instance *postgresql.Instance, model *flagModel, apiClient postgresqlUtils.PostgreSQLClient) (postgresql.ApiUpdateInstanceRequest, error) {
	req := apiClient.UpdateInstance(ctx, model.ProjectId, model.InstanceId)

	payload, err := buildCurrentPayload(instance)
	if err != nil {
		return req, fmt.Errorf("build payload from the current instance parameters: %w", err)
	}

	// Override payload with the command parameters
	var planId *string
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

	if planId != nil {
		payload.PlanId = planId
	}
	if model.EnableMonitoring != nil {
		payload.Parameters.EnableMonitoring = model.EnableMonitoring
	}
	if model.Graphite != nil {
		payload.Parameters.Graphite = model.Graphite
	}
	if model.MonitoringInstanceId != nil {
		payload.Parameters.MonitoringInstanceId = model.MonitoringInstanceId
	}
	if model.MetricsFrequency != nil {
		payload.Parameters.MetricsFrequency = model.MetricsFrequency
	}
	if model.MetricsPrefix != nil {
		payload.Parameters.MetricsPrefix = model.MetricsPrefix
	}
	if model.Plugin != nil {
		payload.Parameters.Plugins = model.Plugin
	}
	if model.SgwAcl != nil {
		payload.Parameters.SgwAcl = sgwAcl
	}
	if model.Syslog != nil {
		payload.Parameters.Syslog = model.Syslog
	}

	req = req.UpdateInstancePayload(*payload)
	return req, nil
}

// Builds the payload from the current instance parameters
func buildCurrentPayload(instance *postgresql.Instance) (*postgresql.UpdateInstancePayload, error) {
	if instance == nil {
		return nil, fmt.Errorf("instance is nil")
	}

	currentParameters := *instance.Parameters
	var ok bool
	var currentEnableMonitoring bool
	var currentGraphite string
	var currentMonitoringInstanceId string
	var currentMetricsFrequency int64
	var currentMetricsPrefix string
	var currentPlugins []string
	var currentSyslog []string
	var currentSgwAcl string
	if currentParameters != nil {
		if currentParameters["enable_monitoring"] != nil {
			currentEnableMonitoring, ok = currentParameters["enable_monitoring"].(bool)
			if !ok {
				return nil, fmt.Errorf("parse enable_monitoring: type cannot be converted to bool")
			}
		}
		if currentParameters["graphite"] != nil {
			currentGraphite, ok = currentParameters["graphite"].(string)
			if !ok {
				return nil, fmt.Errorf("parse graphite: type cannot be converted to string")
			}
		}
		if currentParameters["monitoring_instance_id"] != nil {
			currentMonitoringInstanceId, ok = currentParameters["monitoring_instance_id"].(string)
			if !ok {
				return nil, fmt.Errorf("parse monitoring_instance_id: type cannot be converted to string")
			}
		}
		if currentParameters["metrics_frequency"] != nil {
			currentMetricsFrequency, ok = currentParameters["metrics_frequency"].(int64)
			if !ok {
				return nil, fmt.Errorf("parse metrics_frequency: type cannot be converted to int64")
			}
		}
		if currentParameters["metrics_prefix"] != nil {
			currentMetricsPrefix, ok = currentParameters["metrics_prefix"].(string)
			if !ok {
				return nil, fmt.Errorf("parse metrics_prefix: type cannot be converted to string")
			}
		}
		if currentParameters["plugins"] != nil {
			currentPlugins, ok = currentParameters["plugins"].([]string)
			if !ok {
				return nil, fmt.Errorf("parse plugins: type cannot be converted to []string")
			}
		}
		if currentParameters["syslog"] != nil {
			currentSyslog, ok = currentParameters["syslog"].([]string)
			if !ok {
				return nil, fmt.Errorf("parse syslog: type cannot be converted to []string")
			}
		}
		if currentParameters["sgw_acl"] != nil {
			currentSgwAcl, ok = currentParameters["sgw_acl"].(string)
			if !ok {
				return nil, fmt.Errorf("parse sgw_acl: type cannot be converted to string")
			}
		}
	}
	payload := &postgresql.UpdateInstancePayload{
		Parameters: &postgresql.InstanceParameters{
			EnableMonitoring:     &currentEnableMonitoring,
			Graphite:             &currentGraphite,
			MonitoringInstanceId: &currentMonitoringInstanceId,
			MetricsFrequency:     &currentMetricsFrequency,
			MetricsPrefix:        &currentMetricsPrefix,
			Plugins:              &currentPlugins,
			Syslog:               &currentSyslog,
			SgwAcl:               &currentSgwAcl,
		},
		PlanId: instance.PlanId,
	}

	return payload, nil
}
