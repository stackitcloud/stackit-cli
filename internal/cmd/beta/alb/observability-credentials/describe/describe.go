package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const (
	credentialRefArg = "CREDENTIAL_REF" // nolint:gosec // false alert, these are not valid credentials
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialRef string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", credentialRefArg),
		Short: "Describes observability credentials for the Application Load Balancer",
		Long:  "Describes observability credentials for the Application Load Balancer.",
		Args:  args.SingleArg(credentialRefArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details about credentials with name "credential-12345"`,
				"$ stackit beta alb observability-credentials describe credential-12345",
			),
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
				return fmt.Errorf("read credentials: %w", err)
			}

			if credential := resp; credential != nil && credential.Credential != nil {
				return outputResult(params.Printer, model.OutputFormat, *credential.Credential)
			}
			params.Printer.Outputln("No credentials found.")
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	credentialRef := inputArgs[0]
	model := inputModel{
		GlobalFlagModel: globalFlags,
		CredentialRef:   credentialRef,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) alb.ApiGetCredentialsRequest {
	return apiClient.GetCredentials(ctx, model.ProjectId, model.Region, model.CredentialRef)
}

func outputResult(p *print.Printer, outputFormat string, response alb.CredentialsResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(response, "", "  ")

		if err != nil {
			return fmt.Errorf("marshal credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(response, yaml.IndentSequence(true), yaml.UseJSONMarshaler())

		if err != nil {
			return fmt.Errorf("marshal credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("CREDENTIAL REF", utils.PtrString(response.CredentialsRef))
		table.AddSeparator()
		table.AddRow("DISPLAYNAME", utils.PtrString(response.DisplayName))
		table.AddSeparator()
		table.AddRow("UESRNAME", utils.PtrString(response.Username))
		table.AddSeparator()
		table.AddRow("REGION", utils.PtrString(response.Region))
		table.AddSeparator()

		p.Outputln(table.Render())
	}

	return nil
}
