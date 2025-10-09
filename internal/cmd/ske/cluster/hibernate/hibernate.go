package hibernate

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	"github.com/stackitcloud/stackit-sdk-go/services/ske/wait"
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
		Use:   fmt.Sprintf("hibernate %s", clusterNameArg),
		Short: "Trigger hibernate for a SKE cluster",
		Long:  "Trigger hibernate for a STACKIT Kubernetes Engine (SKE) cluster.",
		Args:  args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Trigger hibernate for a SKE cluster with name "my-cluster"`,
				"$ stackit ske cluster hibernate my-cluster"),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to trigger hibernate for %q in project %q?", model.ClusterName, projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("hibernate SKE cluster: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Hibernating cluster")
				_, err = wait.TriggerClusterHibernationWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.ClusterName).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SKE cluster hibernation: %w", err)
				}
				s.Stop()
			}

			operationState := "Hibernated"
			if model.Async {
				operationState = "Triggered hibernation of"
			}
			params.Printer.Outputf("%s cluster %q\n", operationState, model.ClusterName)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiTriggerHibernateRequest {
	req := apiClient.TriggerHibernate(ctx, model.ProjectId, model.Region, model.ClusterName)
	return req
}
