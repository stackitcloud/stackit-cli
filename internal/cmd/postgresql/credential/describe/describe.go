package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
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
	GlobalFlags  *globalflags.Model
	InstanceId   string
	CredentialId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "describe",
		Short:   "Get details of a PostgreSQL instance credential",
		Long:    "Get details of a PostgreSQL instance credential",
		Example: `$ stackit postgresql credential describe --project-id xxx --instance-id xxx --credential-id xxx`,
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
				return fmt.Errorf("describe PostgreSQL credential: %w", err)
			}

			// Show details
			details, err := json.MarshalIndent(resp, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal PostgreSQL credential: %w", err)
			}
			cmd.Println(string(details))

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
	globalFlags := globalflags.Parse()
	if globalFlags.ProjectId == "" {
		return nil, fmt.Errorf("project ID not set")
	}

	return &flagModel{
		GlobalFlags:  globalFlags,
		InstanceId:   utils.FlagToStringValue(cmd, instanceIdFlag),
		CredentialId: utils.FlagToStringValue(cmd, credentialIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiGetCredentialsRequest {
	req := apiClient.GetCredentials(ctx, model.GlobalFlags.ProjectId, model.InstanceId, model.CredentialId)
	return req
}
