package addtarget

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"

	"github.com/spf13/cobra"
)

const (
	loadBalancerNameArg = "LOAD_BALANCER_NAME"

	targetPoolNameFlag = "target-pool-name"
	targetNameFlag     = "target-name"
	ipFlag             = "ip"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	LoadBalancerName string
	TargetPoolName   string
	TargetName       string
	Ip               string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("add-target %s", loadBalancerNameArg),
		Short: "Adds a target to a target pool",
		Long:  "Adds a target to a target pool.",
		Args:  args.SingleArg(loadBalancerNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Add a target to target pool "my-target-pool"  of load balancer with name "my-load-balancer"`,
				"$ stackit load-balancer target-pool add-target my-load-balancer --target-pool-name my-target-pool --target-name my-new-target --ip 1.2.3.4/32"),
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
				prompt := fmt.Sprintf("Are you sure you want to add a target with IP %q to target pool %q of load balancer %q?", model.Ip, model.TargetPoolName, model.LoadBalancerName)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("build request: %w", err)
			}
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("add target to target pool: %w", err)
			}

			p.Info("Added target to target pool of load balancer %q\n", model.LoadBalancerName)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(targetPoolNameFlag, "", "Target pool name")
	cmd.Flags().String(targetNameFlag, "", "Target name")
	cmd.Flags().String(ipFlag, "", "Target IP. Must by unique within a target pool")

	err := flags.MarkFlagsRequired(cmd, targetPoolNameFlag, targetNameFlag, ipFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	lbName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		LoadBalancerName: lbName,
		TargetPoolName:   cmd.Flag(targetPoolNameFlag).Value.String(),
		TargetName:       cmd.Flag(targetNameFlag).Value.String(),
		Ip:               cmd.Flag(ipFlag).Value.String(),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient utils.LoadBalancerClient) (loadbalancer.ApiUpdateTargetPoolRequest, error) {
	req := apiClient.UpdateTargetPool(ctx, model.ProjectId, model.LoadBalancerName, model.TargetPoolName)

	targetPool, err := utils.GetLoadBalancerTargetPool(ctx, apiClient, model.ProjectId, model.LoadBalancerName, model.TargetPoolName)
	if err != nil {
		return req, fmt.Errorf("get load balancer target pool: %w", err)
	}

	newTarget := &loadbalancer.Target{
		DisplayName: &model.TargetName,
		Ip:          &model.Ip,
	}
	err = utils.AddTargetToTargetPool(targetPool, newTarget)
	if err != nil {
		return req, fmt.Errorf("add target to target pool: %w", err)
	}

	payload := utils.ToPayloadTargetPool(targetPool)
	if payload == nil {
		return req, fmt.Errorf("nil payload")
	}

	return req.UpdateTargetPoolPayload(*payload), nil
}
