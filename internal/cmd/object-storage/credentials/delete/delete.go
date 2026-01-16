package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	objectStorageUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	credentialsIdArg       = "CREDENTIALS_ID" //nolint:gosec // linter false positive
	credentialsGroupIdFlag = "credentials-group-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialsGroupId string
	CredentialsId      string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", credentialsIdArg),
		Short: "Deletes credentials of an Object Storage credentials group",
		Long:  "Deletes credentials of an Object Storage credentials group",
		Args:  args.SingleArg(credentialsIdArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete a credential with ID "xxx" of credentials group with ID "yyy"`,
				"$ stackit object-storage credentials delete xxx --credentials-group-id yyy"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			credentialsGroupLabel, err := objectStorageUtils.GetCredentialsGroupName(ctx, apiClient, model.ProjectId, model.CredentialsGroupId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get credentials group name: %v", err)
				credentialsGroupLabel = model.CredentialsGroupId
			}

			credentialsLabel, err := objectStorageUtils.GetCredentialsName(ctx, apiClient, model.ProjectId, model.CredentialsGroupId, model.CredentialsId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get credentials name: %v", err)
				credentialsLabel = model.CredentialsId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete credentials %q of credentials group  %q? (This cannot be undone)", credentialsLabel, credentialsGroupLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete Object Storage credentials: %w", err)
			}

			params.Printer.Info("Deleted credentials %q of credentials group %q\n", credentialsLabel, credentialsGroupLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), credentialsGroupIdFlag, "Credentials Group ID")

	err := flags.MarkFlagsRequired(cmd, credentialsGroupIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	credentialsId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		CredentialsGroupId: flags.FlagToStringValue(p, cmd, credentialsGroupIdFlag),
		CredentialsId:      credentialsId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiDeleteAccessKeyRequest {
	req := apiClient.DeleteAccessKey(ctx, model.ProjectId, model.Region, model.CredentialsId)
	req = req.CredentialsGroup(model.CredentialsGroupId)
	return req
}
