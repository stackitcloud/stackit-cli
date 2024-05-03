package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	loadbalancerUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/utils"

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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", credentialsRefArg),
		Short: "Deletes observability credentials for Load Balancer",
		Long:  "Deletes observability credentials for Load Balancer.",
		Args:  args.SingleArg(credentialsRefArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete credentials with reference "credentials-xxx" for Load Balancer`,
				"$ stackit loadbalancer credentials delete credentials-xxx"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			credentialsLabel, err := loadbalancerUtils.GetCredentialsDisplayName(ctx, apiClient, model.ProjectId, model.CredentialsRef)
			if err != nil {
				p.Debug(print.ErrorLevel, "get credentials display name: %v", err)
				credentialsLabel = model.CredentialsRef
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete observability credentials %q? (This cannot be undone)", credentialsLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete Load Balancer observability credentials: %w", err)
			}

			p.Info("Deleted credentials %q\n", credentialsLabel)
			return nil
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

	return &inputModel{
		GlobalFlagModel: globalFlags,
		CredentialsRef:  credentialsRef,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiDeleteCredentialsRequest {
	req := apiClient.DeleteCredentials(ctx, model.ProjectId, model.CredentialsRef)
	return req
}
