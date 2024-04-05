package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/redis/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/redis"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", credentialsIdArg),
		Short: "Shows details of credentials of a Redis instance",
		Long:  "Shows details of credentials of a Redis instance. The password will be shown in plain text in the output.",
		Args:  args.SingleArg(credentialsIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of credentials with ID "xxx" from instance with ID "yyy"`,
				"$ stackit redis credentials describe xxx --instance-id yyy"),
			examples.NewExample(
				`Get details of credentials with ID "xxx" from instance with ID "yyy" in a table format`,
				"$ stackit redis credentials describe xxx --instance-id yyy --output-format pretty"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("describe Redis credentials: %w", err)
			}

			return outputResult(cmd, model.OutputFormat, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *redis.APIClient) redis.ApiGetCredentialsRequest {
	req := apiClient.GetCredentials(ctx, model.ProjectId, model.InstanceId, model.CredentialsId)
	return req
}

func outputResult(cmd *cobra.Command, outputFormat string, credentials *redis.CredentialsResponse) error {
	switch outputFormat {
	case globalflags.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("ID", *credentials.Id)
		table.AddSeparator()
		// The username field cannot be set by the user so we only display it if it's not returned empty
		username := *credentials.Raw.Credentials.Username
		if username != "" {
			table.AddRow("USERNAME", *credentials.Raw.Credentials.Username)
			table.AddSeparator()
		}
		table.AddRow("PASSWORD", *credentials.Raw.Credentials.Password)
		table.AddSeparator()
		table.AddRow("URI", *credentials.Raw.Credentials.Uri)
		err := table.Display(cmd)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(credentials, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Redis credentials: %w", err)
		}
		cmd.Println(string(details))

		return nil
	}
}
