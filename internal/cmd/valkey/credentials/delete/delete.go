package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/client"
	valkeyUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	credentialsIdArg = "CREDENTIALS_ID"

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
		Short: "Deletes credentials of a Valkey instance",
		Long:  "Deletes credentials of a Valkey instance.",
		Args:  args.SingleArg(credentialsIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete credentials with ID "xxx" of a Valkey instance with ID "yyy"`,
				"$ stackit valkey credentials delete xxx --instance-id yyy"),
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

			instanceLabel, err := valkeyUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.InstanceId, model.Region)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			credentialsLabel, err := valkeyUtils.GetCredentialsUsername(ctx, apiClient.DefaultAPI, model.ProjectId, model.InstanceId, model.CredentialsId, model.Region)
			if err != nil || credentialsLabel == "" {
				params.Printer.Debug(print.ErrorLevel, "get credentials username: %v", err)
				credentialsLabel = model.CredentialsId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete credentials %s of instance %q? (This cannot be undone)", credentialsLabel, instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete Valkey credentials: %w", err)
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
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		CredentialsId:   credentialsId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient valkey.DefaultAPI) valkey.ApiDeleteCredentialsRequest {
	return apiClient.DeleteCredentials(ctx, model.ProjectId, model.Region, model.InstanceId, model.CredentialsId)
}
