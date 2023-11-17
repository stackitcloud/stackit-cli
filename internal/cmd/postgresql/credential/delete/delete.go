package delete

import (
	"context"
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
	projectIdFlag     = "project-id"
	instanceIdFlag    = "instance-id"
	credentialsIdFlag = "credentials-id"
)

type flagModel struct {
	ProjectId     string
	InstanceId    string
	CredentialsId string
}

var Cmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a PostgreSQL instance credential",
	Long:    "Delete a PostgreSQL instance credential",
	Example: `$ stackit postgresql credential delete --project-id xxx --instance-id xxx --credentials-id xxx`,
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
		err = req.Execute()
		if err != nil {
			return fmt.Errorf("delete PostgreSQL credentials: %w", err)
		}

		fmt.Printf("Deleted credentials with ID %s\n", model.CredentialsId)
		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Var(flags.UUIDFlag(), credentialsIdFlag, "Credentials ID")

	err := utils.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
	err = utils.MarkFlagsRequired(cmd, credentialsIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:     projectId,
		InstanceId:    utils.FlagToStringValue(cmd, instanceIdFlag),
		CredentialsId: utils.FlagToStringValue(cmd, credentialsIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiDeleteCredentialsRequest {
	req := apiClient.DeleteCredentials(ctx, model.ProjectId, model.InstanceId, model.CredentialsId)
	return req
}
