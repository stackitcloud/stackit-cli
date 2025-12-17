package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	credentialsRefArg = "CREDENTIALS_REF" //nolint:gosec // linter false positive
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialsRef string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", credentialsRefArg),
		Short: "Shows details of observability credentials for Load Balancer",
		Long:  "Shows details of observability credentials for Load Balancer.",
		Args:  args.SingleArg(credentialsRefArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get details of observability credentials with reference "credentials-xxx"`,
				"$ stackit load-balancer observability-credentials describe credentials-xxx"),
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
				return fmt.Errorf("describe Load Balancer observability credentials: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	credentialsRef := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		CredentialsRef:  credentialsRef,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiGetCredentialsRequest {
	req := apiClient.GetCredentials(ctx, model.ProjectId, model.Region, model.CredentialsRef)
	return req
}

func outputResult(p *print.Printer, outputFormat string, credentials *loadbalancer.GetCredentialsResponse) error {
	return p.OutputResult(outputFormat, credentials, func() error {
		if credentials == nil || credentials.Credential == nil {
			return fmt.Errorf("credentials response is empty")
		}

		table := tables.NewTable()
		table.AddRow("REFERENCE", utils.PtrString(credentials.Credential.CredentialsRef))
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", utils.PtrString(credentials.Credential.DisplayName))
		table.AddSeparator()
		table.AddRow("USERNAME", utils.PtrString(credentials.Credential.Username))
		table.AddSeparator()
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
