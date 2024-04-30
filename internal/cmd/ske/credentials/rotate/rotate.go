package rotate

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("rotate %s", clusterNameArg),
		Short: "Rotates credentials associated to a SKE cluster",
		Long:  "Rotates credentials associated to a STACKIT Kubernetes Engine (SKE) cluster. The old credentials will be invalid after the operation.",
		Args:  args.SingleArg(clusterNameArg, nil),
		Deprecated: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
			"and will be removed in a future release.",
			"Please use the 2-step credential rotation flow instead, by running the commands:",
			" $ stackit ske credentials start-rotation CLUSTER_NAME",
			" $ stackit ske credentials complete-rotation CLUSTER_NAME",
			"For more information, visit: https://docs.stackit.cloud/stackit/en/how-to-rotate-ske-credentials-200016334.html",
		),
		Example: examples.Build(
			examples.NewExample(
				`Rotate credentials associated to the SKE cluster with name "my-cluster"`,
				"$ stackit ske credentials rotate my-cluster"),
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
				prompt := fmt.Sprintf("Are you sure you want to rotate credentials for SKE cluster %q? (The old credentials will be invalid after this operation)", model.ClusterName)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("rotate SKE credentials: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Rotating credentials")
				_, err = wait.RotateCredentialsWaitHandler(ctx, apiClient, model.ProjectId, model.ClusterName).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for SKE credentials rotation: %w", err)
				}
				s.Stop()
			}

			operationState := "Rotated"
			if model.Async {
				operationState = "Triggered rotation of"
			}
			p.Info("%s credentials for cluster %q\n", operationState, model.ClusterName)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiTriggerRotateCredentialsRequest {
	req := apiClient.TriggerRotateCredentials(ctx, model.ProjectId, model.ClusterName)
	return req
}
