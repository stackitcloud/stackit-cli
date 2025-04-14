package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const (
	loadbalancerNameArg = "LOADBALANCER_NAME_ARG"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", loadbalancerNameArg),
		Short: "Describes an application loadbalancer",
		Long:  "Describes an application alb.",
		Args:  args.SingleArg(loadbalancerNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details about an application loadbalancer with name "my-load-balancer"`,
				"$ stackit beta alb describe my-load-balancer",
			),
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
				return fmt.Errorf("read loadbalancer: %w", err)
			}

			if loadbalancer := resp; loadbalancer != nil {
				return outputResult(p, model.OutputFormat, loadbalancer)
			}
			p.Outputln("No load balancer found.")
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	loadbalancerName := inputArgs[0]
	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            loadbalancerName,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) alb.ApiGetLoadBalancerRequest {
	return apiClient.GetLoadBalancer(ctx, model.ProjectId, model.Region, model.Name)
}

func outputResult(p *print.Printer, outputFormat string, response *alb.LoadBalancer) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(response, "", "  ")

		if err != nil {
			return fmt.Errorf("marshal loadbalancer: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(response, yaml.IndentSequence(true), yaml.UseJSONMarshaler())

		if err != nil {
			return fmt.Errorf("marshal loadbalancer: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		if err := outputResultAsTable(p, response); err != nil {
			return err
		}
	}

	return nil
}

func outputResultAsTable(p *print.Printer, loadbalancer *alb.LoadBalancer) error {
	content := []tables.Table{}

	content = append(content, buildLoadBalancerTable(loadbalancer))

	if loadbalancer.Listeners != nil {
		content = append(content, buildListenersTable(*loadbalancer.Listeners))
	}

	if loadbalancer.TargetPools != nil {
		content = append(content, buildTargetPoolsTable(*loadbalancer.TargetPools))
	}

	err := tables.DisplayTables(p, content)
	if err != nil {
		return fmt.Errorf("display output: %w", err)
	}

	return nil
}

func buildLoadBalancerTable(loadbalancer *alb.LoadBalancer) tables.Table {
	acl := []string{}
	privateAccessOnly := false
	if loadbalancer.Options != nil {
		if loadbalancer.Options.AccessControl != nil && loadbalancer.Options.AccessControl.AllowedSourceRanges != nil {
			acl = *loadbalancer.Options.AccessControl.AllowedSourceRanges
		}

		if loadbalancer.Options.PrivateNetworkOnly != nil {
			privateAccessOnly = *loadbalancer.Options.PrivateNetworkOnly
		}
	}

	networkId := "-"
	if loadbalancer.Networks != nil && len(*loadbalancer.Networks) > 0 {
		networks := *loadbalancer.Networks
		networkId = *networks[0].NetworkId
	}

	externalAddress := utils.PtrStringDefault(loadbalancer.ExternalAddress, "-")

	errorDescriptions := []string{}
	if loadbalancer.Errors != nil && len((*loadbalancer.Errors)) > 0 {
		for _, err := range *loadbalancer.Errors {
			errorDescriptions = append(errorDescriptions, *err.Description)
		}
	}

	table := tables.NewTable()
	table.SetTitle("Load Balancer")
	table.AddRow("NAME", utils.PtrString(loadbalancer.Name))
	table.AddSeparator()
	table.AddRow("STATE", utils.PtrString(loadbalancer.Status))
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

func buildListenersTable(listeners []alb.Listener) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Listeners")
	table.SetHeader("NAME", "PORT", "PROTOCOL", "TARGET POOL")
	for i := range listeners {
		listener := listeners[i]
		table.AddRow(
			utils.PtrString(listener.Name),
			utils.PtrString(listener.Port),
			utils.PtrString(listener.Protocol),
		)
	}
	return table
}

func buildTargetPoolsTable(targetPools []alb.TargetPool) tables.Table {
	table := tables.NewTable()
	table.SetTitle("Target Pools")
	table.SetHeader("NAME", "PORT", "TARGETS")
	for _, targetPool := range targetPools {
		table.AddRow(utils.PtrString(targetPool.Name), utils.PtrString(targetPool.TargetPort), len(*targetPool.Targets))
	}
	return table
}
