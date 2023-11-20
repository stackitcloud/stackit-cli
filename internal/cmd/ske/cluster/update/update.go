package update

import (
	"context"
	"fmt"
	"os"

	"github.com/stackitcloud/stackit-cli/internal/cmd/ske/cluster/create"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske/wait"
)

var Cmd = &cobra.Command{
	Use:     "update",
	Short:   "Updates a Kubernetes cluster",
	Long:    "Updates a Kubernetes cluster",
	Example: `$ stackit ske cluster update --project-id xxx --payload "{\"field\":\value\"}"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		model, err := create.ParseFlags(cmd, os.ReadFile)
		if err != nil {
			return err
		}

		// Configure API client
		apiClient, err := client.ConfigureClient(cmd)
		if err != nil {
			return fmt.Errorf("authentication failed, please run \"stackit auth login\" or \"stackit auth activate-service-account\"")
		}

		// Check if cluster exists
		clusters, err := apiClient.GetClustersExecute(ctx, model.ProjectId)
		if err != nil {
			return fmt.Errorf("get SKE cluster: %w", err)
		}
		exists := false
		for _, cl := range *clusters.Items {
			if cl.Name != nil && *cl.Name == model.Name {
				exists = true
				break
			}
		}
		if !exists {
			return fmt.Errorf("cluster with name %s does not exist", model.Name)
		}

		// Call API
		req, err := create.BuildRequest(ctx, model, apiClient)
		if err != nil {
			return fmt.Errorf("build SKE cluster update request: %w", err)
		}
		resp, err := req.Execute()
		if err != nil {
			return fmt.Errorf("update SKE cluster: %w", err)
		}

		// Wait for async operation
		name := *resp.Name
		_, err = wait.CreateOrUpdateClusterWaitHandler(ctx, apiClient, model.ProjectId, name).WaitWithContext(ctx)
		if err != nil {
			return fmt.Errorf("wait for SKE cluster update: %w", err)
		}

		fmt.Printf("Updated cluster with name %s\n", name)
		return nil
	},
}

func init() {
	create.ConfigureFlags(Cmd)
}
