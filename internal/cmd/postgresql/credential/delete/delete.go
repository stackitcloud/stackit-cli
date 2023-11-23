package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/commonflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/postgresql/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/postgresql"
)

const (
	instanceIdFlag   = "instance-id"
	credentialIdFlag = "credential-id" //nolint:gosec // linter false positive
)

type flagModel struct {
	ProjectId    string
	InstanceId   string
	CredentialId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete a PostgreSQL instance credential",
		Long:    "Delete a PostgreSQL instance credential",
		Example: `$ stackit postgresql credential delete --project-id xxx --instance-id xxx --credential-id xxx`,
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
				return fmt.Errorf("delete PostgreSQL credential: %w", err)
			}

			cmd.Printf("Deleted credential with ID %s\n", model.CredentialId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Var(flags.UUIDFlag(), credentialIdFlag, "Credentials ID")

	err := utils.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
	err = utils.MarkFlagsRequired(cmd, credentialIdFlag)
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	projectId := commonflags.GetString(commonflags.ProjectIdFlag)
	if projectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		ProjectId:    projectId,
		InstanceId:   utils.FlagToStringValue(cmd, instanceIdFlag),
		CredentialId: utils.FlagToStringValue(cmd, credentialIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiDeleteCredentialsRequest {
	req := apiClient.DeleteCredentials(ctx, model.ProjectId, model.InstanceId, model.CredentialId)
	return req
}
