package create

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/mariadb/client"
	mariadbUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/mariadb/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/mariadb"
	"github.com/stackitcloud/stackit-sdk-go/services/mariadb/wait"
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
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	PlanName string
	Version  string

	InstanceName         *string
	EnableMonitoring     *bool
	Graphite             *string
	MetricsFrequency     *int64
	MetricsPrefix        *string
	MonitoringInstanceId *string
	SgwAcl               *[]string
	Syslog               *[]string
	PlanId               *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a MariaDB instance",
		Long:  "Creates a MariaDB instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a MariaDB instance with name "my-instance" and specify plan by name and version`,
				"$ stackit mariadb instance create --name my-instance --plan-name stackit-mariadb-1.2.10-replica --version 10.6"),
			examples.NewExample(
				`Create a MariaDB instance with name "my-instance" and specify plan by ID`,
				"$ stackit mariadb instance create --name my-instance --plan-id xxx"),
			examples.NewExample(
				`Create a MariaDB instance with name "my-instance" and specify IP range which is allowed to access it`,
				"$ stackit mariadb instance create --name my-instance --plan-id xxx --acl 1.2.3.0/24"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a MariaDB instance for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
				if !errors.As(err, &dsaInvalidPlanError) {
					return fmt.Errorf("build MariaDB instance creation request: %w", err)
				}
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create MariaDB instance: %w", err)
			}
			instanceId := *resp.InstanceId

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating instance")
				_, err = wait.CreateInstanceWaitHandler(ctx, apiClient, model.ProjectId, instanceId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for MariaDB instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, projectLabel, resp)
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
	cmd.Flags().Var(flags.CIDRSliceFlag(), sgwAclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")
	cmd.Flags().StringSlice(syslogFlag, []string{}, "Syslog")
	cmd.Flags().Var(flags.UUIDFlag(), planIdFlag, "Plan ID")
	cmd.Flags().String(planNameFlag, "", "Plan name")
	cmd.Flags().String(versionFlag, "", "Instance MariaDB version")

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
		InstanceName:         flags.FlagToStringPointer(p, cmd, instanceNameFlag),
		EnableMonitoring:     flags.FlagToBoolPointer(p, cmd, enableMonitoringFlag),
		MonitoringInstanceId: flags.FlagToStringPointer(p, cmd, monitoringInstanceIdFlag),
		Graphite:             flags.FlagToStringPointer(p, cmd, graphiteFlag),
		MetricsFrequency:     flags.FlagToInt64Pointer(p, cmd, metricsFrequencyFlag),
		MetricsPrefix:        flags.FlagToStringPointer(p, cmd, metricsPrefixFlag),
		SgwAcl:               flags.FlagToStringSlicePointer(p, cmd, sgwAclFlag),
		Syslog:               flags.FlagToStringSlicePointer(p, cmd, syslogFlag),
		PlanId:               planId,
		PlanName:             planName,
		Version:              version,
	}

	p.DebugInputModel(model)
	return &model, nil
}

type mariaDBClient interface {
	CreateInstance(ctx context.Context, projectId string) mariadb.ApiCreateInstanceRequest
	ListOfferingsExecute(ctx context.Context, projectId string) (*mariadb.ListOfferingsResponse, error)
}

func buildRequest(ctx context.Context, model *inputModel, apiClient mariaDBClient) (mariadb.ApiCreateInstanceRequest, error) {
	req := apiClient.CreateInstance(ctx, model.ProjectId)

	var planId *string
	var err error

	offerings, err := apiClient.ListOfferingsExecute(ctx, model.ProjectId)
	if err != nil {
		return req, fmt.Errorf("get MariaDB offerings: %w", err)
	}

	if model.PlanId == nil {
		planId, err = mariadbUtils.LoadPlanId(model.PlanName, model.Version, offerings)
		if err != nil {
			var dsaInvalidPlanError *cliErr.DSAInvalidPlanError
			if !errors.As(err, &dsaInvalidPlanError) {
				return req, fmt.Errorf("load plan ID: %w", err)
			}
			return req, err
		}
	} else {
		err := mariadbUtils.ValidatePlanId(*model.PlanId, offerings)
		if err != nil {
			return req, err
		}
		planId = model.PlanId
	}

	var sgwAcl *string
	if model.SgwAcl != nil {
		sgwAcl = utils.Ptr(strings.Join(*model.SgwAcl, ","))
	}

	req = req.CreateInstancePayload(mariadb.CreateInstancePayload{
		InstanceName: model.InstanceName,
		Parameters: &mariadb.InstanceParameters{
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

func outputResult(p *print.Printer, outputFormat string, async bool, projectLabel string, resp *mariadb.CreateInstanceResponse) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, utils.PtrString(resp.InstanceId))
		return nil
	})
}
