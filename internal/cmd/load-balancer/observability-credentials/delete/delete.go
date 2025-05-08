package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", credentialsRefArg),
		Short: "Deletes observability credentials for Load Balancer",
		Long:  "Deletes observability credentials for Load Balancer.",
		Args:  args.SingleArg(credentialsRefArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete observability credentials with reference "credentials-xxx" for Load Balancer`,
				"$ stackit load-balancer observability-credentials delete credentials-xxx"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			credentialsLabel, err := loadbalancerUtils.GetCredentialsDisplayName(ctx, apiClient, model.ProjectId, model.Region, model.CredentialsRef)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get observability credentials display name: %v", err)
				credentialsLabel = model.CredentialsRef
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to delete observability credentials %q on project %q?(This cannot be undone)", credentialsLabel, projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
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

			params.Printer.Info("Deleted observability credentials %q on project %q\n", credentialsLabel, projectLabel)
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		CredentialsRef:  credentialsRef,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiDeleteCredentialsRequest {
	req := apiClient.DeleteCredentials(ctx, model.ProjectId, model.Region, model.CredentialsRef)
	return req
}
