package describe

import (
	"context"
	"encoding/json"
	"fmt"

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
		Long:  "Describes an application loadbalancer.",
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
		table := tables.NewTable()
		table.AddRow("EXTERNAL ADDRESS", utils.PtrString(response.ExternalAddress))
		table.AddSeparator()
		var numErrors int
		if response.Errors != nil {
			numErrors = len(*response.Errors)
		}
		table.AddRow("NUMBER OF ERRORS", numErrors)
		table.AddSeparator()
		table.AddRow("PLAN ID", utils.PtrString(response.PlanId))
		table.AddSeparator()
		table.AddRow("REGION", utils.PtrString(response.Region))
		table.AddSeparator()
		table.AddRow("STATUS", utils.PtrString(response.Status))
		table.AddSeparator()
		table.AddRow("VERSION", utils.PtrString(response.Version))

		p.Outputln(table.Render())
	}

	return nil
}
