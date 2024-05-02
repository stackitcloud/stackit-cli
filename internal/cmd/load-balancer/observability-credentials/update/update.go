package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	loadBalancerUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	displayNameFlag = "display-name"
	usernameFlag    = "username"
	passwordFlag    = "password"

	credentialsRefArg = "CREDENTIALS_REF" //nolint:gosec // linter false positive
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	CredentialsRef string
	DisplayName    *string
	Username       *string
	Password       *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates observability credentials for Load Balancer",
		Long:  "Updates existing observability credentials (username and password) for Load Balancer. The credentials can be for Argus or another monitoring tool.",
		Args:  args.SingleArg(credentialsRefArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update the password of credentials of Load Balancer with credentials reference "credentials-xxx". The password is entered using the terminal`,
				"$ stackit load-balancer observability-credentials update credentials-xxx --password "),
			examples.NewExample(
				`Update the password of credentials of Load Balancer with credentials reference "credentials-xxx", by providing it in the flag`,
				"$ stackit load-balancer observability-credentials update credentials-xxx --password new-pwd"),
			examples.NewExample(
				`Update the display name of credentials of Load Balancer with credentials reference "credentials-xxx".`,
				"$ stackit load-balancer observability-credentials update credentials-xxx --display-name yyy"),
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			credentialsLabel, err := loadBalancerUtils.GetCredentialsDisplayName(ctx, apiClient, model.ProjectId, model.CredentialsRef)
			if err != nil {
				p.Debug(print.ErrorLevel, "get credentials display name: %v", err)
				credentialsLabel = model.CredentialsRef
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update observability credentials %q for your Load Balancers on project %q?", credentialsLabel, projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Prompt for password if not passed in as a flag
			if model.Password != nil && *model.Password == "" {
				pwd, err := p.PromptForPassword("Enter new password: ")
				if err != nil {
					return fmt.Errorf("prompt for password: %w", err)
				}
				model.Password = utils.Ptr(pwd)
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Load Balancer observability credentials: %w", err)
			}

			p.Info("Updated observability credentials %q for Load Balancer\n", credentialsLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "Credentials name")
	cmd.Flags().String(usernameFlag, "", "Username")
	cmd.Flags().String(passwordFlag, "", "Password")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	credentialsRef := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	displayName := flags.FlagToStringPointer(p, cmd, displayNameFlag)
	username := flags.FlagToStringPointer(p, cmd, usernameFlag)
	password := flags.FlagToStringPointer(p, cmd, passwordFlag)

	if displayName == nil && username == nil && password == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		CredentialsRef:  credentialsRef,
		DisplayName:     displayName,
		Username:        username,
		Password:        password,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiUpdateCredentialsRequest {
	req := apiClient.UpdateCredentials(ctx, model.ProjectId, model.CredentialsRef)

	req = req.UpdateCredentialsPayload(loadbalancer.UpdateCredentialsPayload{
		DisplayName: model.DisplayName,
		Username:    model.Username,
		Password:    model.Password,
	})
	return req
}
