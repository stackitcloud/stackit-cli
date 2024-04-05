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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	credentialsGroupNameFlag = "name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialsGroupName string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a credentials group to hold Object Storage access credentials",
		Long:  "Creates a credentials group to hold Object Storage access credentials.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials group to hold Object Storage access credentials`,
				"$ stackit object-storage credentials-group create --name example"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a credentials group with name %q?", model.CredentialsGroupName)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Object Storage credentials group: %w", err)
			}

			cmd.Printf("Created credentials group %q. Credentials group ID: %s\n\n", *resp.CredentialsGroup.DisplayName, *resp.CredentialsGroup.CredentialsGroupId)
			cmd.Printf("URN: %s\n", *resp.CredentialsGroup.Urn)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(credentialsGroupNameFlag, "", "Name of the group holding credentials")

	err := flags.MarkFlagsRequired(cmd, credentialsGroupNameFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel:      globalFlags,
		CredentialsGroupName: flags.FlagToStringValue(cmd, credentialsGroupNameFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiCreateCredentialsGroupRequest {
	req := apiClient.CreateCredentialsGroup(ctx, model.ProjectId)
	req = req.CreateCredentialsGroupPayload(objectstorage.CreateCredentialsGroupPayload{
		DisplayName: utils.Ptr(model.CredentialsGroupName),
	})
	return req
}
