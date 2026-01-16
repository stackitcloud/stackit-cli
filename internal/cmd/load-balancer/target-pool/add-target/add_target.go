package addtarget

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

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
	ipArg = "TARGET_IP"

	lbNameFlag         = "lb-name"
	targetNameFlag     = "target-name"
	targetPoolNameFlag = "target-pool-name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	TargetPoolName string
	LBName         string
	TargetName     string
	IP             string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("add-target %s", ipArg),
		Short: "Adds a target to a target pool",
		Long: fmt.Sprintf("%s\n%s",
			"Adds a target to a target pool.",
			"The target IP must by unique within a target pool and must be a valid IPv4 or IPv6."),
		Args: args.SingleArg(ipArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Add a target with IP 1.2.3.4 and name "my-new-target" to target pool "my-target-pool" of load balancer with name "my-load-balancer"`,
				"$ stackit load-balancer target-pool add-target 1.2.3.4 --target-name my-new-target --target-pool-name my-target-pool --lb-name my-load-balancer"),
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

			prompt := fmt.Sprintf("Are you sure you want to add a target with IP %q to target pool %q of load balancer %q?", model.IP, model.TargetPoolName, model.LBName)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
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

			params.Printer.Info("Added target to target pool of load balancer %q\n", model.LBName)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(targetNameFlag, "n", "", "Target name")
	cmd.Flags().String(targetPoolNameFlag, "", "Target pool name")
	cmd.Flags().String(lbNameFlag, "", "Load balancer name")

	err := flags.MarkFlagsRequired(cmd, lbNameFlag, targetNameFlag, targetPoolNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	ip := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		TargetPoolName:  cmd.Flag(targetPoolNameFlag).Value.String(),
		LBName:          cmd.Flag(lbNameFlag).Value.String(),
		TargetName:      cmd.Flag(targetNameFlag).Value.String(),
		IP:              ip,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient utils.LoadBalancerClient) (loadbalancer.ApiUpdateTargetPoolRequest, error) {
	req := apiClient.UpdateTargetPool(ctx, model.ProjectId, model.Region, model.LBName, model.TargetPoolName)

	targetPool, err := utils.GetLoadBalancerTargetPool(ctx, apiClient, model.ProjectId, model.Region, model.LBName, model.TargetPoolName)
	if err != nil {
		return req, fmt.Errorf("get load balancer target pool: %w", err)
	}

	newTarget := &loadbalancer.Target{
		DisplayName: &model.TargetName,
		Ip:          &model.IP,
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
