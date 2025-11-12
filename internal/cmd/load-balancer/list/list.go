package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	limitFlag = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Load Balancers",
		Long:  "Lists all Load Balancers.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all load balancers`,
				"$ stackit load-balancer list"),
			examples.NewExample(
				`List all loadbalancers in JSON format`,
				"$ stackit load-balancer list --output-format json"),
			examples.NewExample(
				`List up to 10 load balancers `,
				"$ stackit load-balancer list --limit 10"),
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
				return fmt.Errorf("get load balancers: %w", err)
			}

			if resp.LoadBalancers == nil || (resp.LoadBalancers != nil && len(*resp.LoadBalancers) == 0) {
				projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
					projectLabel = model.ProjectId
				}
				params.Printer.Info("No load balancers found for project %q\n", projectLabel)
				return nil
			}

			loadBalancers := *resp.LoadBalancers
			// Truncate output
			if model.Limit != nil && len(loadBalancers) > int(*model.Limit) {
				loadBalancers = loadBalancers[:*model.Limit]
			}

			return outputResult(params.Printer, model.OutputFormat, loadBalancers)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           flags.FlagToInt64Pointer(p, cmd, limitFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiListLoadBalancersRequest {
	req := apiClient.ListLoadBalancers(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, loadBalancers []loadbalancer.LoadBalancer) error {
	return p.OutputResult(outputFormat, loadBalancers, func() error {
		table := tables.NewTable()
		table.SetHeader("NAME", "STATE", "IP ADDRESS", "LISTENERS", "TARGET POOLS")
		for i := range loadBalancers {
			l := loadBalancers[i]
			var numListeners, numTargetPools int
			if l.Listeners != nil {
				numListeners = len(*l.Listeners)
			}
			if l.TargetPools != nil {
				numTargetPools = len(*l.TargetPools)
			}

			externalAddress := utils.PtrStringDefault(l.ExternalAddress, "-")
			table.AddRow(
				utils.PtrString(l.Name),
				utils.PtrString(l.Status),
				externalAddress,
				numListeners,
				numTargetPools,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
