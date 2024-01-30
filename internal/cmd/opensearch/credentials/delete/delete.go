package delete

import (
	"context"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/confirm"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/services/opensearch/client"
	opensearchUtils "stackit/internal/pkg/services/opensearch/utils"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/opensearch"
)

const (
	credentialsIdArg = "CREDENTIAL_ID" //nolint:gosec // linter false positive

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId    string
	CredentialsId string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", credentialsIdArg),
		Short: "Delete credentials of an OpenSearch instance",
		Long:  "Delete credentials of an OpenSearch instance",
		Args:  args.SingleArg(credentialsIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete credentials with ID "xxx" of OpenSearch instance with ID "yyy"`,
				"$ stackit opensearch credentials delete xxx --instance-id yyy"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			instanceLabel, err := opensearchUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				instanceLabel = model.InstanceId
			}

			credentialLabel, err := opensearchUtils.GetCredentialUsername(ctx, apiClient, model.ProjectId, model.InstanceId, model.CredentialsId)
			if err != nil {
				credentialLabel = model.CredentialsId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete credential %s of instance %s? (This cannot be undone)", credentialLabel, instanceLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete OpenSearch credentials: %w", err)
			}

			cmd.Printf("Deleted credential %s of instance %s\n", credentialLabel, instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	credentialsId := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(cmd, instanceIdFlag),
		CredentialsId:   credentialsId,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *opensearch.APIClient) opensearch.ApiDeleteCredentialsRequest {
	req := apiClient.DeleteCredentials(ctx, model.ProjectId, model.InstanceId, model.CredentialsId)
	return req
}
