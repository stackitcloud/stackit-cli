package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresql/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql"
)

const (
	projectIdFlag  = "project-id"
	instanceIdFlag = "instance-id"
)

type flagModel struct {
	ProjectId  string
	InstanceId string
}

var Cmd = &cobra.Command{
	Use:     "describe",
	Short:   "Get details of a PostgreSQL instance",
	Long:    "Get details of a PostgreSQL instance",
	Example: `$ stackit postgresql instance describe --project-id xxx --instance-id xxx`,
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
			return fmt.Errorf("read PostgreSQL instance: %w", err)
		}

		// Show details
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal PostgreSQL instance: %w", err)
		}
		cmd.Println(string(details))

		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := utils.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:  projectId,
		InstanceId: utils.FlagToStringValue(cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}
