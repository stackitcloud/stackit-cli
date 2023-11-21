package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
	"github.com/stackitcloud/stackit-sdk-go/services/ske/wait"
)

const (
	projectIdFlag   = "project-id"
	clusterNameFlag = "name"
)

type flagModel struct {
	ProjectId   string
	ClusterName string
}

var Cmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a SKE cluster",
	Long:    "Delete a SKE cluster",
	Example: `$ stackit ske cluster delete --project-id xxx --name xxx`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		model, err := parseFlags(cmd)
		if err != nil {
			return err
		}
		// Configure API client
		apiClient, err := client.ConfigureClient(cmd)
		if err != nil {
			return fmt.Errorf("authentication failed, please run \"stackit auth login\" or \"stackit auth activate-service-account\"")
		}

		// Call API
		req := buildRequest(ctx, model, apiClient)
		_, err = req.Execute()
		if err != nil {
			return fmt.Errorf("delete SKE cluster: %w", err)
		}

		// Wait for async operation
		_, err = wait.DeleteClusterWaitHandler(ctx, apiClient, model.ProjectId, model.ClusterName).WaitWithContext(ctx)
		if err != nil {
			return fmt.Errorf("wait for SKE cluster deletion: %w", err)
		}

		fmt.Println("Cluster deleted")
		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(clusterNameFlag, "", "Cluster name")

	err := utils.MarkFlagsRequired(cmd, clusterNameFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}
	clusterName := utils.FlagToStringValue(cmd, clusterNameFlag)
	if clusterName == "" {
		return nil, fmt.Errorf("cluster name can't be empty")
	}

	return &flagModel{
		ProjectId:   projectId,
		ClusterName: clusterName,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *ske.APIClient) ske.ApiDeleteClusterRequest {
	req := apiClient.DeleteCluster(ctx, model.ProjectId, model.ClusterName)
	return req
}
