package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	lbUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	targetPoolNameArg = "TARGET_POOL_NAME"

	lbNameFlag = "lb-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	TargetPoolName string
	LBName         string
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
				"$ stackit load-balancer target-pool describe pool --lb-name my-load-balancer"),
			examples.NewExample(
				`Get details of a target pool with name "pool" in load balancer with name "my-load-balancer in JSON output"`,
				"$ stackit load-balancer target-pool describe pool --lb-name my-load-balancer --output-format json"),
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

			targetPool := lbUtils.FindLoadBalancerTargetPoolByName(*resp.TargetPools, model.TargetPoolName)
			if targetPool == nil {
				return fmt.Errorf("target pool not found")
			}

			listener := lbUtils.FindLoadBalancerListenerByTargetPool(*resp.Listeners, *targetPool.Name)

			return outputResult(p, model.OutputFormat, *targetPool, listener)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(lbNameFlag, "", "Name of the load balancer")

	err := flags.MarkFlagsRequired(cmd, lbNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	targetPoolName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		TargetPoolName:  targetPoolName,
		LBName:          cmd.Flag(lbNameFlag).Value.String(),
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
	req := apiClient.GetLoadBalancer(ctx, model.ProjectId, model.LBName)
	return req
}

func outputResult(p *print.Printer, outputFormat string, targetPool loadbalancer.TargetPool, listener *loadbalancer.Listener) error {
	output := struct {
		*loadbalancer.TargetPool
		Listener *loadbalancer.Listener `json:"attached_listener"`
	}{
		&targetPool,
		listener,
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal load balancer: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(output, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal load balancer: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		return outputResultAsTable(p, targetPool, listener)
	}
}

func outputResultAsTable(p *print.Printer, targetPool loadbalancer.TargetPool, listener *loadbalancer.Listener) error {
	sessionPersistence := "None"
	if targetPool.SessionPersistence != nil && targetPool.SessionPersistence.UseSourceIpAddress != nil && *targetPool.SessionPersistence.UseSourceIpAddress {
		sessionPersistence = "Use Source IP"
	}

	healthCheckInterval := "-"
	healthCheckUnhealthyThreshold := "-"
	healthCheckHealthyThreshold := "-"
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
	}

	targets := "-"
	if targetPool.Targets != nil {
		var targetsSlice []string
		for _, target := range *targetPool.Targets {
			targetStr := fmt.Sprintf("%s (%s)", *target.DisplayName, *target.Ip)
			targetsSlice = append(targetsSlice, targetStr)
		}
		targets = strings.Join(targetsSlice, "\n")
	}

	listenerStr := "-"
	if listener != nil {
		listenerStr = fmt.Sprintf("%s (Port:%s, Protocol: %s)",
			utils.PtrString(listener.Name),
			utils.PtrString(listener.Port),
			utils.PtrString(listener.Protocol),
		)
	}

	table := tables.NewTable()
	table.AddRow("NAME", *targetPool.Name)
	table.AddSeparator()
	table.AddRow("TARGET PORT", *targetPool.TargetPort)
	table.AddSeparator()
	table.AddRow("ATTACHED LISTENER", listenerStr)
	table.AddSeparator()
	table.AddRow("TARGETS", targets)
	table.AddSeparator()
	table.AddRow("SESSION PERSISTENCE", sessionPersistence)
	table.AddSeparator()
	table.AddRow("HEALTH CHECK INTERVAL", healthCheckInterval)
	table.AddSeparator()
	table.AddRow("HEALTH CHECK DOWN AFTER", healthCheckUnhealthyThreshold)
	table.AddSeparator()
	table.AddRow("HEALTH CHECK UP AFTER", healthCheckHealthyThreshold)
	table.AddSeparator()

	err := p.PagerDisplay(table.Render())
	if err != nil {
		return fmt.Errorf("display output: %w", err)
	}

	return nil
}
