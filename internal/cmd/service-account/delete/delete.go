package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/service-account/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serviceaccount"
)

const (
	emailArg = "EMAIL"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Email string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", emailArg),
		Short: "Deletes a service account",
		Long:  "Deletes a service account.",
		Args:  args.SingleArg(emailArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete a service account with email "my-service-account-1234567@sa.stackit.cloud"`,
				"$ stackit service-account delete my-service-account-1234567@sa.stackit.cloud"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete service account %s? (This cannot be undone)", model.Email)
				err = p.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete service account: %w", err)
			}

			p.Info("Service account %s deleted\n", model.Email)
			return nil
		},
	}
	return cmd
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	email := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Email:           email,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serviceaccount.APIClient) serviceaccount.ApiDeleteServiceAccountRequest {
	req := apiClient.DeleteServiceAccount(ctx, model.ProjectId, model.Email)
	return req
}
