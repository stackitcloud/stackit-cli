package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
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
		Use:   fmt.Sprintf("describe %s", credentialsIdArg),
		Short: "Shows details of credentials of a Valkey instance",
		Long:  "Shows details of credentials of a Valkey instance. The password will be shown in plain text in the output.",
		Args:  args.SingleArg(credentialsIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of credentials with ID "xxx" from a Valkey instance with ID "yyy"`,
				"$ stackit valkey credentials describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of credentials with ID "xxx" from a Valkey instance with ID "yyy" in JSON format`,
				"$ stackit valkey credentials describe xxx --instance-id yyy --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe Valkey credentials: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient valkey.DefaultAPI) valkey.ApiGetCredentialsRequest {
	return apiClient.GetCredentials(ctx, model.ProjectId, model.Region, model.InstanceId, model.CredentialsId)
}

func outputResult(p *print.Printer, outputFormat string, credentials *valkey.CredentialsResponse) error {
	return p.OutputResult(outputFormat, credentials, func() error {
		if credentials == nil {
			return fmt.Errorf("no credentials found")
		}
		table := tables.NewTable()
		table.AddRow("ID", credentials.Id)
		table.AddSeparator()
		if credentials.Raw != nil {
			if username := credentials.Raw.Credentials.Username; username != "" {
				table.AddRow("USERNAME", username)
				table.AddSeparator()
			}
			table.AddRow("PASSWORD", credentials.Raw.Credentials.Password)
			table.AddSeparator()
			table.AddRow("URI", utils.PtrString(credentials.Raw.Credentials.Uri))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
