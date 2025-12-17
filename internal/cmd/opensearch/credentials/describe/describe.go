package describe

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
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
		Use:   fmt.Sprintf("describe %s", credentialsIdArg),
		Short: "Shows details of credentials of an OpenSearch instance",
		Long:  "Shows details of credentials of an OpenSearch instance. The password will be shown in plain text in the output.",
		Args:  args.SingleArg(credentialsIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of credentials with ID "xxx" from instance with ID "yyy"`,
				"$ stackit opensearch credentials describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of credentials with ID "xxx" from instance with ID "yyy" in JSON format`,
				"$ stackit opensearch credentials describe xxx --instance-id yyy --output-format json"),
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
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe OpenSearch credentials: %w", err)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *opensearch.APIClient) opensearch.ApiGetCredentialsRequest {
	req := apiClient.GetCredentials(ctx, model.ProjectId, model.InstanceId, model.CredentialsId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, credentials *opensearch.CredentialsResponse) error {
	if credentials == nil {
		return fmt.Errorf("credentials is nil")
	}

	return p.OutputResult(outputFormat, credentials, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(credentials.Id))
		table.AddSeparator()
		// The username field cannot be set by the user so we only display it if it's not returned empty
		if credentials.HasRaw() && credentials.Raw.Credentials != nil {
			if username := credentials.Raw.Credentials.Username; username != nil && *username != "" {
				table.AddRow("USERNAME", *username)
				table.AddSeparator()
			}
			table.AddRow("PASSWORD", utils.PtrString(credentials.Raw.Credentials.Password))
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
