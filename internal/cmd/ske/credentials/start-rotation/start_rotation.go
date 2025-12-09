package startrotation

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

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

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("start-rotation %s", clusterNameArg),
		Short: "Starts the rotation of the credentials associated to a SKE cluster",
		Long: fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n\n%s\n%s\n%s\n%s\n%s",
			"Starts the rotation of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster.",
			"This is step 1 of a 2-step process to rotate all SKE cluster credentials. Tasks accomplished in this phase include:",
			"  - Rolling recreation of all worker nodes",
			"  - A new Certificate Authority (CA) will be established and incorporated into the existing CA bundle.",
			"  - A new etcd encryption key is generated and added to the Certificate Authority (CA) bundle.",
			"  - A new signing key will be generated for the service account and added to the Certificate Authority (CA) bundle.",
			"  - The kube-apiserver will rewrite all secrets in the cluster, encrypting them with the new encryption key.",
			"The old CA, encryption key and signing key will be retained until the rotation is completed.",
			"After completing the rotation of credentials, you can generate a new kubeconfig file by running:",
			"  $ stackit ske kubeconfig create my-cluster",
			"Complete the rotation by running:",
			"  $ stackit ske credentials complete-rotation my-cluster",
			"For more information, visit: https://docs.stackit.cloud/products/runtime/kubernetes-engine/how-tos/rotate-ske-credentials/",
		),
		Args: args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Start the rotation of the credentials associated to the SKE cluster with name "my-cluster"`,
				"$ stackit ske credentials start-rotation my-cluster"),
			examples.NewExample(
				`Flow of the 2-step process to rotate all SKE cluster credentials, including generating a new kubeconfig file`,
				"$ stackit ske credentials start-rotation my-cluster",
				"$ stackit ske kubeconfig create my-cluster",
				"$ stackit ske credentials complete-rotation my-cluster",
			),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to start the rotation of the credentials for SKE cluster %q?", model.ClusterName)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("start rotation of SKE credentials: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Starting credentials rotation")
				_, err = wait.StartCredentialsRotationWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.ClusterName).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for start SKE credentials rotation %w", err)
				}
				s.Stop()
			}

			operationState := "Rotation of credentials is ready to be completed"
			if model.Async {
				operationState = "Triggered start of credentials rotation"
			}
			params.Printer.Info("%s for cluster %q\n", operationState, model.ClusterName)
			params.Printer.Info("Complete the rotation by running:\n  $ stackit ske credentials complete-rotation %s\n", model.ClusterName)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiStartCredentialsRotationRequest {
	req := apiClient.StartCredentialsRotation(ctx, model.ProjectId, model.Region, model.ClusterName)
	return req
}
