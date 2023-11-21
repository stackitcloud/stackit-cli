package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
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
	Use:     "describe",
	Short:   "Get details of a SKE cluster",
	Long:    "Get details of a SKE cluster",
	Example: `$ stackit ske cluster describe --project-id xxx --name xxx`,
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
		resp, err := req.Execute()
		if err != nil {
			return fmt.Errorf("read SKE cluster: %w", err)
		}

		// Show details
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE cluster: %w", err)
		}
		cmd.Println(string(details))

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

func buildRequest(ctx context.Context, model *flagModel, apiClient *ske.APIClient) ske.ApiGetClusterRequest {
	req := apiClient.GetCluster(ctx, model.ProjectId, model.ClusterName)
	return req
}
