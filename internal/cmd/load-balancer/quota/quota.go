package quota

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quota",
		Short: "Shows the configured Load Balancer quota",
		Long:  "Shows the configured Load Balancer quota for the project. If you want to change the quota, please create a support ticket in the STACKIT Help Center (https://support.stackit.cloud)",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get the configured load balancer quota for the project`,
				"$ stackit load-balancer quota"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
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
				return fmt.Errorf("get load balancer quota: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiGetQuotaRequest {
	req := apiClient.GetQuota(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, quota *loadbalancer.GetQuotaResponse) error {
	if quota == nil {
		return fmt.Errorf("quota response is empty")
	}
	switch outputFormat {
	case print.PrettyOutputFormat:

		maxLoadBalancers := "Unlimited"
		if quota.MaxLoadBalancers != nil && *quota.MaxLoadBalancers != -1 {
			maxLoadBalancers = strconv.FormatInt(*quota.MaxLoadBalancers, 10)
		}

		p.Outputf("Maximum number of load balancers allowed: %s\n", maxLoadBalancers)

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(quota, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal quota: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		details, err := json.MarshalIndent(quota, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal quota: %w", err)
		}

		p.Outputln(string(details))

		return nil
	}
}
