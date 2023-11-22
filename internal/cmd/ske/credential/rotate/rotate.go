package rotate

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	clusterNameFlag = "cluster-name"
)

type flagModel struct {
	*globalflags.GlobalFlagModel
	ClusterName string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "rotate",
		Short:   "Rotate credential associated to a SKE cluster",
		Long:    "Rotate credential associated to a SKE cluster. The old credential will be invalid after the operation",
		Example: `$ stackit postgresql credential create --project-id xxx --instance-id xxx`,
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
				return fmt.Errorf("rotate SKE credential: %w", err)
			}

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(clusterNameFlag, "", "Cluster name")

	err := utils.MarkFlagsRequired(cmd, clusterNameFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	globalFlags := globalflags.Parse()
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	clusterName := utils.FlagToStringValue(cmd, clusterNameFlag)
	if clusterName == "" {
		return nil, fmt.Errorf("cluster name can't be empty")
	}

	return &flagModel{
		GlobalFlagModel: globalFlags,
		ClusterName:     clusterName,
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *ske.APIClient) ske.ApiTriggerRotateCredentialsRequest {
	req := apiClient.TriggerRotateCredentials(ctx, model.ProjectId, model.ClusterName)
	return req
}
