package create

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
	hidePasswordFlag = "hide-password"
)

type flagModel struct {
	ProjectId    string
	InstanceId   string
	HidePassword bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Create credentials for a PostgreSQL instance",
		Long:    "Create credentials for a PostgreSQL instance",
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
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create PostgreSQL credential: %w", err)
			}

			cmd.Printf("Created credential with ID %s\n\nUsername: %s\n", *resp.Id, *resp.Raw.Credentials.Username)
			if model.HidePassword {
				cmd.Printf("Password: <hidden>\n")
			} else {
				cmd.Printf("Password: %s\n", *resp.Raw.Credentials.Password)
			}
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Bool(hidePasswordFlag, false, "Hide password in output")

	err := utils.MarkFlagsRequired(cmd, instanceIdFlag)
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
		HidePassword: utils.FlagToBoolValue(cmd, hidePasswordFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *flagModel, apiClient *postgresql.APIClient) postgresql.ApiCreateCredentialsRequest {
	req := apiClient.CreateCredentials(ctx, model.ProjectId, model.InstanceId)
	return req
}
