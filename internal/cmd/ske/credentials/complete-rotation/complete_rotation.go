package completerotation

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
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

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("complete-rotation %s", clusterNameArg),
		Short: "Completes the rotation of the credentials associated to a SKE cluster",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Completes the rotation of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster.",
			"This is step 2 of a two-step process, if you haven't, please start the process by running:",
			"  $ stackit ske credentials start-rotation my-cluster"),
		Args: args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Complete the rotation of the credentials associated to the SKE cluster with name "my-cluster"`,
				"$ stackit ske credentials complete-rotation my-cluster"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to complete the rotation of the credentials for SKE cluster %q?", model.ClusterName)
				err = confirm.PromptForConfirmation(cmd, prompt)
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
				s := spinner.New(cmd)
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
			cmd.Printf("%s for cluster %q\n", operationState, model.ClusterName)
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
