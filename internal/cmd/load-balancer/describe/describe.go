package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	loadBalancerNameArg = "LOAD_BALANCER_NAME"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	LoadBalancerName string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", loadBalancerNameArg),
		Short: "Shows details of a Load Balancer",
		Long:  "Shows details of a Load Balancer.",
		Args:  args.SingleArg(loadBalancerNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a load balancer with name "my-load-balancer"`,
				"$ stackit load-balancer describe my-load-balancer"),
			examples.NewExample(
				`Get details of a load-balancer with name "my-load-balancer" in a JSON format`,
				"$ stackit load-balancer describe my-load-balancer --output-format json"),
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

			return outputResult(p, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	loadBalancerName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		LoadBalancerName: loadBalancerName,
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

func outputResult(p *print.Printer, outputFormat string, loadBalancer *loadbalancer.LoadBalancer) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(loadBalancer, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal load balancer: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		return outputResultAsTable(p, loadBalancer)
	}
}

func outputResultAsTable(p *print.Printer, loadBalancer *loadbalancer.LoadBalancer) error {
	content := renderLoadBalancer(loadBalancer)

	if loadBalancer.Listeners != nil {
		content += renderListeners(*loadBalancer.Listeners)
	}

	if loadBalancer.TargetPools != nil {
		content += renderTargetPools(*loadBalancer.TargetPools)
	}

	err := p.PagerDisplay(content)
	if err != nil {
		return fmt.Errorf("display output: %w", err)
	}

	return nil
}

func renderLoadBalancer(loadBalancer *loadbalancer.LoadBalancer) string {
	acl := []string{}
	privateAccessOnly := false
	if loadBalancer.Options != nil {
		if loadBalancer.Options.AccessControl != nil && loadBalancer.Options.AccessControl.AllowedSourceRanges != nil {
			acl = *loadBalancer.Options.AccessControl.AllowedSourceRanges
		}

		if loadBalancer.Options.PrivateNetworkOnly != nil {
			privateAccessOnly = *loadBalancer.Options.PrivateNetworkOnly
		}
	}

	networkId := "-"
	if loadBalancer.Networks != nil && len(*loadBalancer.Networks) > 0 {
		networks := *loadBalancer.Networks
		networkId = *networks[0].NetworkId
	}

	externalAdress := "-"
	if loadBalancer.ExternalAddress != nil {
		externalAdress = *loadBalancer.ExternalAddress
	}

	table := tables.NewTable()
	table.AddRow("NAME", *loadBalancer.Name)
	table.AddSeparator()
	table.AddRow("STATE", *loadBalancer.Status)
	table.AddSeparator()
	table.AddRow("PRIVATE ACCESS ONLY", privateAccessOnly)
	table.AddSeparator()
	table.AddRow("ATTACHED PUBLIC IP", externalAdress)
	table.AddSeparator()
	table.AddRow("ATTACHED NETWORK ID", networkId)
	table.AddSeparator()
	table.AddRow("ACL", acl)
	return table.Render()
}

func renderListeners(listeners []loadbalancer.Listener) string {
	table := tables.NewTable()
	table.SetHeader("LISTENER NAME", "PORT", "PROTOCOL", "TARGET POOL")
	for i := range listeners {
		listener := listeners[i]
		table.AddRow(*listener.Name, *listener.Port, *listener.Protocol, *listener.TargetPool)
	}
	return table.Render()
}

func renderTargetPools(targetPools []loadbalancer.TargetPool) string {
	table := tables.NewTable()
	table.SetHeader("TARGET POOL NAME", "PORT", "TARGETS")
	for _, targetPool := range targetPools {
		table.AddRow(*targetPool.Name, *targetPool.TargetPort, len(*targetPool.Targets))
	}
	return table.Render()
}
