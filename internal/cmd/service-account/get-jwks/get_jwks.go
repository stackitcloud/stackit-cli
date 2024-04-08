package getjwks

import (
	"context"
	"encoding/json"
	"fmt"

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

func NewCmd(p *print.Printer) *cobra.Command {
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
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
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
				p.Info("Empty JWKS for service account %s\n", model.Email)
				return nil
			}

			return outputResult(p, jwks)
		},
	}

	return cmd
}

func parseInput(_ *cobra.Command, inputArgs []string) (*inputModel, error) {
	email := inputArgs[0]

	return &inputModel{
		Email: email,
	}, nil
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
