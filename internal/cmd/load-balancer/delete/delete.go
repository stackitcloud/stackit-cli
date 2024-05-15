package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer/wait"
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
		Use:   fmt.Sprintf("delete %s", loadBalancerNameArg),
		Short: "Deletes a Load Balancer",
		Long:  "Deletes a Load Balancer.",
		Args:  args.SingleArg(loadBalancerNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Deletes a load balancer with name "my-load-balancer"`,
				"$ stackit load-balancer delete my-load-balancer"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete load balancer %q? (This cannot be undone)", model.LoadBalancerName)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete load balancer: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Deleting load balancer")
				_, err = wait.DeleteLoadBalancerWaitHandler(ctx, apiClient, model.ProjectId, model.LoadBalancerName).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for load balancer deletion: %w", err)
				}
				s.Stop()
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			p.Info("%s load balancer %q\n", operationState, model.LoadBalancerName)
			return nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiDeleteLoadBalancerRequest {
	req := apiClient.DeleteLoadBalancer(ctx, model.ProjectId, model.LoadBalancerName)
	return req
}
