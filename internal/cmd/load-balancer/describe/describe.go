package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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

func NewCmd(params *types.CmdParams) *cobra.Command {
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
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read load balancer: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiGetLoadBalancerRequest {
	req := apiClient.GetLoadBalancer(ctx, model.ProjectId, model.Region, model.LoadBalancerName)
	return req
}

func outputResult(p *print.Printer, outputFormat string, loadBalancer *loadbalancer.LoadBalancer) error {
	if loadBalancer == nil {
		return fmt.Errorf("loadbalancer response is empty")
	}

	return p.OutputResult(outputFormat, loadBalancer, func() error {
		content := []tables.Table{}
		content = append(content, buildLoadBalancerTable(loadBalancer))

		if loadBalancer.Listeners != nil {
			content = append(content, buildListenersTable(*loadBalancer.Listeners))
		}
		if loadBalancer.TargetPools != nil {
			content = append(content, buildTargetPoolsTable(*loadBalancer.TargetPools))
		}

		err := tables.DisplayTables(p, content)
		if err != nil {
			return fmt.Errorf("display output: %w", err)
		}

		return nil
	})
}

func buildLoadBalancerTable(loadBalancer *loadbalancer.LoadBalancer) tables.Table {
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

	externalAddress := utils.PtrStringDefault(loadBalancer.ExternalAddress, "-")

	errorDescriptions := []string{}
	if loadBalancer.Errors != nil && len((*loadBalancer.Errors)) > 0 {
		for _, err := range *loadBalancer.Errors {
			errorDescriptions = append(errorDescriptions, *err.Description)
		}
	}

	table := tables.NewTable()
	table.SetTitle("Load Balancer")
	table.AddRow("NAME", utils.PtrString(loadBalancer.Name))
	table.AddSeparator()
	table.AddRow("STATE", utils.PtrString(loadBalancer.Status))
	table.AddSeparator()
	if len(errorDescriptions) > 0 {
		table.AddRow("ERROR DESCRIPTIONS", strings.Join(errorDescriptions, "\n"))
		table.AddSeparator()
	}
	table.AddRow("PRIVATE ACCESS ONLY", privateAccessOnly)
	table.AddSeparator()
	table.AddRow("ATTACHED PUBLIC IP", externalAddress)
	table.AddSeparator()
	table.AddRow("ATTACHED NETWORK ID", networkId)
	table.AddSeparator()
	table.AddRow("ACL", acl)
	return table
}

func buildListenersTable(listeners []loadbalancer.Listener) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Listeners")
	table.SetHeader("NAME", "PORT", "PROTOCOL", "TARGET POOL")
	for i := range listeners {
		listener := listeners[i]
		table.AddRow(
			utils.PtrString(listener.Name),
			utils.PtrString(listener.Port),
			utils.PtrString(listener.Protocol),
			utils.PtrString(listener.TargetPool),
		)
	}
	return table
}

func buildTargetPoolsTable(targetPools []loadbalancer.TargetPool) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Target Pools")
	table.SetHeader("NAME", "PORT", "TARGETS")
	for _, targetPool := range targetPools {
		table.AddRow(utils.PtrString(targetPool.Name), utils.PtrString(targetPool.TargetPort), len(*targetPool.Targets))
	}
	return table
}
