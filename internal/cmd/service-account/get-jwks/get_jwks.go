package getjwks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-account/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
)

const (
	emailArg = "EMAIL"
)

type inputModel struct {
	Email string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("get-jwks %s", emailArg),
		Short: "Shows the JWKS for a service account",
		Long:  "Shows the JSON Web Key set (JWKS) for a service account. Only JSON output is supported.",
		Args:  args.SingleArg(emailArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Get JWKS for the service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account get-jwks my-service-account-1234567@sa.stackit.cloud"),
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
				return fmt.Errorf("get JWKS: %w", err)
			}
			jwks := *resp.Keys
			if len(jwks) == 0 {
				params.Printer.Info("Empty JWKS for service account %s\n", model.Email)
				return nil
			}

			return outputResult(params.Printer, jwks)
		},
	}

	return cmd
}

func parseInput(p *print.Printer, _ *cobra.Command, inputArgs []string) (*inputModel, error) {
	email := inputArgs[0]

	model := inputModel{
		Email: email,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiGetJWKSRequest {
	req := apiClient.GetJWKS(ctx, model.Email)
	return req
}

func outputResult(p *print.Printer, serviceAccounts []serviceaccount.JWK) error {
	details, err := json.MarshalIndent(serviceAccounts, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JWK list: %w", err)
	}
	p.Outputln(string(details))
	return nil
}
