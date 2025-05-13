package update

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const (
	usernameFlag     = "username"
	displaynameFlag  = "displayname"
	passwordFlag     = "password"
	credentialRefArg = "CREDENTIAL_REF_ARG" //nolint:gosec // false alert, these are not valid credentials
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Username       *string
	Displayname    *string
	Password       *string
	CredentialsRef *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", credentialRefArg),
		Short: "Update credentials",
		Long:  "Update credentials.",
		Args:  args.SingleArg(credentialRefArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update the password of observability credentials of Application Load Balancer with credentials reference "credentials-xxx", by providing the path to a file with the new password as flag`,
				"$ stackit beta alb observability-credentials update credentials-xxx --username user1 --displayname user1 --password @./new-password.txt"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model := parseInput(params.Printer, cmd, args)

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req, err := buildRequest(ctx, &model, apiClient)
			if err != nil {
				return err
			}
			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}
			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update credential %q for %q?", *model.CredentialsRef, projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return fmt.Errorf("update credential: %w", err)
				}
			}

			// Call API
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update credential: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("response is nil")
			}

			return outputResult(params.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(usernameFlag, "u", "", "Username for the credentials")
	cmd.Flags().StringP(displaynameFlag, "d", "", "Displayname for the credentials")
	cmd.Flags().Var(flags.ReadFromFileFlag(), passwordFlag, `Password. Can be a string or a file path, if prefixed with "@" (example: @./password.txt).`)

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, displaynameFlag, usernameFlag))
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) (req alb.ApiUpdateCredentialsRequest, err error) {
	req = apiClient.UpdateCredentials(ctx, model.ProjectId, model.Region, *model.CredentialsRef)

	payload := alb.UpdateCredentialsPayload{
		DisplayName: model.Displayname,
		Password:    model.Password,
		Username:    model.Username,
	}

	if model.Displayname == nil && model.Username == nil {
		return req, fmt.Errorf("no attribute to change passed")
	}

	return req.UpdateCredentialsPayload(payload), nil
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) inputModel {
	model := inputModel{
		GlobalFlagModel: globalflags.Parse(p, cmd),
		Username:        flags.FlagToStringPointer(p, cmd, usernameFlag),
		Displayname:     flags.FlagToStringPointer(p, cmd, displaynameFlag),
		CredentialsRef:  &inputArgs[0],
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

	return model
}

func outputResult(p *print.Printer, model inputModel, response *alb.UpdateCredentialsResponse) error {
	var outputFormat string
	if model.GlobalFlagModel != nil {
		outputFormat = model.GlobalFlagModel.OutputFormat
	}
	if response == nil {
		return fmt.Errorf("no response passewd")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(response.Credential, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal credential: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(response.Credential, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal credential: %w", err)
		}
		p.Outputln(string(details))
	default:
		p.Outputf("Updated credential %q\n", utils.PtrString(model.CredentialsRef))
	}
	return nil
}
