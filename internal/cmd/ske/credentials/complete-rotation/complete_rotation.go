package completerotation

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
		Use:   fmt.Sprintf("complete-rotation %s", clusterNameArg),
		Short: "Completes the rotation of the credentials associated to a SKE cluster",
		Long: fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n%s\n\n%s\n%s\n%s",
			"Completes the rotation of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster.",
			"This is step 2 of a 2-step process to rotate all SKE cluster credentials. Tasks accomplished in this phase include:",
			"  - The old certification authority will be dropped from the package.",
			"  - The old signing key for the service account will be dropped from the bundle.",
			"To ensure continued access to the Kubernetes cluster, please update your kubeconfig with the new credentials:",
			"  $ stackit ske kubeconfig create my-cluster",
			"If you haven't, please start the process by running:",
			"  $ stackit ske credentials start-rotation my-cluster",
			"For more information, visit: https://docs.stackit.cloud/stackit/en/how-to-rotate-ske-credentials-200016334.html",
		),
		Args: args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Complete the rotation of the credentials associated to the SKE cluster with name "my-cluster"`,
				"$ stackit ske credentials complete-rotation my-cluster",
			),
			examples.NewExample(
				`Flow of the 2-step process to rotate all SKE cluster credentials, including generating a new kubeconfig file`,
				"$ stackit ske credentials start-rotation my-cluster",
				"$ stackit ske kubeconfig create my-cluster",
				"$ stackit ske credentials complete-rotation my-cluster",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to complete the rotation of the credentials for SKE cluster %q?", model.ClusterName)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("complete rotation of SKE credentials: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Completing credentials rotation")
				_, err = wait.CompleteCredentialsRotationWaitHandler(ctx, apiClient, model.ProjectId, model.ClusterName).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for completing SKE credentials rotation %w", err)
				}
				s.Stop()
			}

			operationState := "Rotation of credentials is completed"
			if model.Async {
				operationState = "Triggered completion of credentials rotation"
			}
			p.Info("%s for cluster %q\n", operationState, model.ClusterName)
			p.Warn("Consider updating your kubeconfig with the new credentials, create a new kubeconfig by running:\n  $ stackit ske kubeconfig create %s\n", model.ClusterName)
			return nil
		},
	}
	return cmd
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	clusterName := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiCompleteCredentialsRotationRequest {
	req := apiClient.CompleteCredentialsRotation(ctx, model.ProjectId, model.ClusterName)
	return req
}
