package cleanup

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "Deletes observability credentials unused by any Load Balancer",
		Long:  "Deletes observability credentials unused by any Load Balancer.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Delete observability credentials unused by any Load Balancer`,
				"$ stackit load-balancer observability-credentials cleanup"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			listReq := buildListCredentialsRequest(ctx, model, apiClient)
			resp, err := listReq.Execute()
			if err != nil {
				return fmt.Errorf("list Load Balancer observability credentials: %w", err)
			}

			var credentials []loadbalancer.CredentialsResponse
			if resp.Credentials != nil && len(*resp.Credentials) > 0 {
				credentials, err = utils.FilterCredentials(ctx, apiClient, *resp.Credentials, model.ProjectId, model.Region, utils.OP_FILTER_UNUSED)
				if err != nil {
					return fmt.Errorf("filter Load Balancer observability credentials: %w", err)
				}
			}

			if len(credentials) == 0 {
				params.Printer.Info("No unused observability credentials found on project %q\n", projectLabel)
				return nil
			}

			if !model.AssumeYes {
				prompt := "Will delete the following unused observability credentials: \n"
				for _, credential := range credentials {
					if credential.DisplayName == nil || credential.Username == nil {
						return fmt.Errorf("list unused Load Balancer observability credentials: credentials %q missing display name or username", *credential.CredentialsRef)
					}
					name := *credential.DisplayName
					username := *credential.Username
					prompt += fmt.Sprintf("  - %s (username: %q)\n", name, username)
				}
				prompt += fmt.Sprintf("Are you sure you want to delete unused observability credentials on project %q? (This cannot be undone)", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			for _, credential := range credentials {
				if credential.CredentialsRef == nil {
					return fmt.Errorf("delete Load Balancer observability credentials: missing credentials reference")
				}
				credentialsRef := *credential.CredentialsRef
				// Call API
				req := buildDeleteCredentialRequest(ctx, model, apiClient, credentialsRef)
				_, err = req.Execute()
				if err != nil {
					return fmt.Errorf("delete Load Balancer observability credentials: %w", err)
				}
			}

			params.Printer.Info("Deleted unused Load Balancer observability credentials on project %q\n", projectLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
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

func buildDeleteCredentialRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient, credentialsRef string) loadbalancer.ApiDeleteCredentialsRequest {
	req := apiClient.DeleteCredentials(ctx, model.ProjectId, model.Region, credentialsRef)
	return req
}

func buildListCredentialsRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiListCredentialsRequest {
	req := apiClient.ListCredentials(ctx, model.ProjectId, model.Region)
	return req
}
