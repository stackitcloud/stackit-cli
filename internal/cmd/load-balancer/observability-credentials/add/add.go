package add

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	displayNameFlag = "display-name"
	usernameFlag    = "username"
	passwordFlag    = "password"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	DisplayName *string
	Username    *string
	Password    *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds observability credentials to Load Balancer",
		Long:  "Adds existing observability credentials (username and password) to Load Balancer. The credentials can be for Argus or another monitoring tool.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Add credentials to a load balancer with username "xxx" and display name "yyy". The password is entered using the terminal`,
				"$ stackit load-balancer observability-credentials add --username xxx --display-name yyy"),
			examples.NewExample(
				`Add credentials to a load balancer with username "xxx" and display name "yyy", providing the password as flag`,
				"$ stackit load-balancer observability-credentials add --username xxx --password pwd --display-name yyy"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to add observability credentials for your Load Balancers on project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Prompt for password if not passed in as a flag
			if model.Password == nil {
				pwd, err := p.PromptForPassword("Enter password: ")
				if err != nil {
					return fmt.Errorf("prompt for password: %w", err)
				}
				model.Password = utils.Ptr(pwd)
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("add Load Balancer credentials: %w", err)
			}

			return outputResult(p, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "Credentials name")
	cmd.Flags().String(usernameFlag, "", "Username")
	cmd.Flags().String(passwordFlag, "", "Password")

	err := flags.MarkFlagsRequired(cmd, displayNameFlag, usernameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Username:        flags.FlagToStringPointer(p, cmd, usernameFlag),
		Password:        flags.FlagToStringPointer(p, cmd, passwordFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiCreateCredentialsRequest {
	req := apiClient.CreateCredentials(ctx, model.ProjectId)
	req = req.XRequestID(uuid.NewString())

	req = req.CreateCredentialsPayload(loadbalancer.CreateCredentialsPayload{
		DisplayName: model.DisplayName,
		Username:    model.Username,
		Password:    model.Password,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *loadbalancer.CreateCredentialsResponse) error {
	if resp.Credential == nil {
		return fmt.Errorf("nil credentials response")
	}

	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Load Balancer credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Added load balancer observability credentials for project %q. Credentials reference: %q\n", projectLabel, *resp.Credential.CredentialsRef)
		return nil
	}
}
