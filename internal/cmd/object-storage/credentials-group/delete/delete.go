package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/client"
	objectStorageUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/object-storage/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/objectstorage"
)

const (
	credentialsGroupIdArg = "CREDENTIALS_GROUP_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialsGroupId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", credentialsGroupIdArg),
		Short: "Deletes a credentials group that holds Object Storage access credentials",
		Long:  "Deletes a credentials group that holds Object Storage access credentials. Only possible if there are no valid credentials (access-keys) left in the group, otherwise it will throw an error.",
		Args:  args.SingleArg(credentialsGroupIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a credentials group with ID "xxx"`,
				"$ stackit object-storage credentials-group delete xxx"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete credentials group %q? (This cannot be undone)", credentialsGroupLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete Object Storage credentials group: %w", err)
			}

			params.Printer.Info("Deleted credentials group %q\n", credentialsGroupLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	credentialsGroupId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		CredentialsGroupId: credentialsGroupId,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *objectstorage.APIClient) objectstorage.ApiDeleteCredentialsGroupRequest {
	req := apiClient.DeleteCredentialsGroup(ctx, model.ProjectId, model.Region, model.CredentialsGroupId)
	return req
}
