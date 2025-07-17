package reconcile

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	clusterNameArg = "CLUSTER_NAME"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ClusterName string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("reconcile %s", clusterNameArg),
		Short: "Trigger reconcile for a SKE cluster",
		Long:  "Trigger reconcile for a STACKIT Kubernetes Engine (SKE) cluster.",
		Args:  args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Trigger reconcile for a SKE cluster with name "my-cluster"`,
				"$ stackit ske cluster reconcile my-cluster"),
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
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("reconcile SKE cluster: %w", err)
			}

			params.Printer.Info("Reconcile got triggered for %q\n", model.ClusterName)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	clusterName := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiTriggerReconcileRequest {
	req := apiClient.TriggerReconcile(ctx, model.ProjectId, model.Region, model.ClusterName)
	return req
}
