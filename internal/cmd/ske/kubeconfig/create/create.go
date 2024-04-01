package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	clusterNameArg = "CLUSTER_NAME"

	expirationFlag = "expiration"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ClusterName    string
	ExpirationTime *string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("create %s", clusterNameArg),
		Short: "Creates a kubeconfig for an SKE cluster",
		Long:  "Creates a kubeconfig for a STACKIT Kubernetes Engine (SKE) cluster.",
		Args:  args.SingleArg(clusterNameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Create a kubeconfig for the SKE cluster with name "my-cluster"`,
				"$ stackit ske kubeconfig create my-cluster"),
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
				prompt := fmt.Sprintf("Are you sure you want to start the rotation of the credentials for SKE cluster %q?", model.ClusterName)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create kubeconfig for SKE cluster: %w", err)
			}

			// Output kubeconfig to stdout
			fmt.Println(*resp.Kubeconfig)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(expirationFlag, "e", "", "Expiration time for the kubeconfig")
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
		ExpirationTime:  flags.FlagToStringPointer(cmd, expirationFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *ske.APIClient) ske.ApiCreateKubeconfigRequest {
	req := apiClient.CreateKubeconfig(ctx, model.ProjectId, model.ClusterName)

	payload := ske.CreateKubeconfigPayload{}

	if model.ExpirationTime != nil {
		payload.ExpirationSeconds = model.ExpirationTime
	}

	return req.CreateKubeconfigPayload(payload)
}
