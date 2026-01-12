package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/opensearch/client"
	opensearchUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/opensearch/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/opensearch"
)

const (
	credentialsIdArg = "CREDENTIALS_ID" //nolint:gosec // linter false positive

	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId    string
	CredentialsId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", credentialsIdArg),
		Short: "Deletes credentials of an OpenSearch instance",
		Long:  "Deletes credentials of an OpenSearch instance.",
		Args:  args.SingleArg(credentialsIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete credentials with ID "xxx" of OpenSearch instance with ID "yyy"`,
				"$ stackit opensearch credentials delete xxx --instance-id yyy"),
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

			instanceLabel, err := opensearchUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			credentialsLabel, err := opensearchUtils.GetCredentialsUsername(ctx, apiClient, model.ProjectId, model.InstanceId, model.CredentialsId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get credentials user name: %v", err)
				credentialsLabel = model.CredentialsId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete credentials %s of instance %q? (This cannot be undone)", credentialsLabel, instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete OpenSearch credentials: %w", err)
			}

			params.Printer.Info("Deleted credentials %s of instance %q\n", credentialsLabel, instanceLabel)
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

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	credentialsId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		CredentialsId:   credentialsId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *opensearch.APIClient) opensearch.ApiDeleteCredentialsRequest {
	req := apiClient.DeleteCredentials(ctx, model.ProjectId, model.InstanceId, model.CredentialsId)
	return req
}
