package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	targetPoolNameArg = "TARGET_POOL_NAME"

	loadBalancerNameFlag = "load-balancer"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TargetPoolName   string
	LoadBalancerName string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", targetPoolNameArg),
		Short: "Shows details of a target pool in a Load Balancer",
		Long:  "Shows details of a target pool in a Load Balancer.",
		Args:  args.SingleArg(targetPoolNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a target pool with name "pool" in load balancer with name "my-load-balancer"`,
				"$ stackit load-balancer target-pool describe pool --load-balancer my-load-balancer"),
			examples.NewExample(
				`Get details of a target pool with name "pool" in load balancer with name "my-load-balancer in JSON output"`,
				"$ stackit load-balancer target-pool describe pool --load-balancer my-load-balancer --output-format json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read load balancer: %w", err)
			}

			return outputResult(p, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(loadBalancerNameFlag, "", "Name of the load balancer")

	err := flags.MarkFlagsRequired(cmd, loadBalancerNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	loadBalancerName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		TargetPoolName:   loadBalancerName,
		LoadBalancerName: cmd.Flag(loadBalancerNameFlag).Value.String(),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiGetLoadBalancerRequest {
	req := apiClient.GetLoadBalancer(ctx, model.ProjectId, model.LoadBalancerName)
	return req
}

func outputResult(p *print.Printer, model *inputModel, loadBalancer *loadbalancer.LoadBalancer) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(loadBalancer, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal load balancer: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		return outputResultAsTable(p, model, loadBalancer)
	}
}

func outputResultAsTable(p *print.Printer, model *inputModel, loadBalancer *loadbalancer.LoadBalancer) error {
	targetPool := utils.FindLoadBalancerTargetPoolByName(loadBalancer.TargetPools, model.TargetPoolName)

	sessionPersistence := "None"
	if targetPool.SessionPersistence != nil && targetPool.SessionPersistence.UseSourceIpAddress != nil && *targetPool.SessionPersistence.UseSourceIpAddress {
		sessionPersistence = "Use Source IP"
	}

	healthCheckInterval := "-"
	healthCheckUnhealthyThreshold := "-"
	healthCheckHealthyThreshold := "-"
	healthCheckTimeout := "-"
	if targetPool.ActiveHealthCheck != nil {
		if targetPool.ActiveHealthCheck.Interval != nil {
			healthCheckInterval = *targetPool.ActiveHealthCheck.Interval
		}
		if targetPool.ActiveHealthCheck.UnhealthyThreshold != nil {
			healthCheckUnhealthyThreshold = strconv.FormatInt(*targetPool.ActiveHealthCheck.UnhealthyThreshold, 10)
		}
		if targetPool.ActiveHealthCheck.HealthyThreshold != nil {
			healthCheckHealthyThreshold = strconv.FormatInt(*targetPool.ActiveHealthCheck.HealthyThreshold, 10)
		}
		if targetPool.ActiveHealthCheck.Timeout != nil {
			healthCheckTimeout = *targetPool.ActiveHealthCheck.Timeout
		}
	}

	targets := "-"
	if targetPool.Targets != nil {
		var targetsSlice []string
		for _, target := range *targetPool.Targets {
			targetStr := fmt.Sprintf("%s: %s", *target.DisplayName, *target.Ip)
			targetsSlice = append(targetsSlice, targetStr)
		}
		targets = strings.Join(targetsSlice, "\n")
	}

	table := tables.NewTable()
	table.AddRow("NAME", *targetPool.Name)
	table.AddSeparator()
	table.AddRow("SESSION PERSISTENCE", sessionPersistence)
	table.AddSeparator()
	table.AddRow("HEALTH CHECK INTERVAL", healthCheckInterval)
	table.AddSeparator()
	table.AddRow("HEALTH CHECK TIMEOUT", healthCheckTimeout)
	table.AddSeparator()
	table.AddRow("HEALTH CHECK UNHEALTHY THRESHOLD", healthCheckUnhealthyThreshold)
	table.AddSeparator()
	table.AddRow("HEALTH CHECK HEALTHY THRESHOLD", healthCheckHealthyThreshold)
	table.AddSeparator()
	table.AddRow("TARGET PORT", *targetPool.TargetPort)
	table.AddSeparator()
	table.AddRow("TARGETS", targets)

	err := p.PagerDisplay(table.Render())
	if err != nil {
		return fmt.Errorf("display output: %w", err)
	}

	return nil
}
