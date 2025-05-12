package add

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds observability credentials to Load Balancer",
		Long:  "Adds existing observability credentials (username and password) to Load Balancer. The credentials can be for Observability or another monitoring tool.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Add observability credentials to a load balancer with username "xxx" and display name "yyy". The password is entered using the terminal`,
				"$ stackit load-balancer observability-credentials add --username xxx --display-name yyy"),
			examples.NewExample(
				`Add observability credentials to a load balancer with username "xxx" and display name "yyy", providing the path to a file with the password as flag`,
				"$ stackit load-balancer observability-credentials add --username xxx --password @./password.txt --display-name yyy"),
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

			// Prompt for password if not passed in as a flag
			if model.Password == nil {
				pwd, err := params.Printer.PromptForPassword("Enter user password: ")
				if err != nil {
					return fmt.Errorf("prompt for password: %w", err)
				}
				model.Password = utils.Ptr(pwd)
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to add observability credentials for Load Balancer on project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("add Load Balancer observability credentials: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "Credentials display name")
	cmd.Flags().String(usernameFlag, "", "Username")
	cmd.Flags().Var(flags.ReadFromFileFlag(), passwordFlag, `Password. Can be a string or a file path, if prefixed with "@" (example: @./password.txt).`)

	err := flags.MarkFlagsRequired(cmd, displayNameFlag, usernameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Username:        flags.FlagToStringPointer(p, cmd, usernameFlag),
		Password:        flags.FlagToStringPointer(p, cmd, passwordFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiCreateCredentialsRequest {
	req := apiClient.CreateCredentials(ctx, model.ProjectId, model.Region)
	req = req.XRequestID(uuid.NewString())

	req = req.CreateCredentialsPayload(loadbalancer.CreateCredentialsPayload{
		DisplayName: model.DisplayName,
		Username:    model.Username,
		Password:    model.Password,
	})
	return req
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *loadbalancer.CreateCredentialsResponse) error {
	if resp == nil || resp.Credential == nil {
		return fmt.Errorf("nil observability credentials response")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Load Balancer observability credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Load Balancer observability credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Added Load Balancer observability credentials on project %q. Credentials reference: %q\n", projectLabel, utils.PtrString(resp.Credential.CredentialsRef))
		return nil
	}
}
