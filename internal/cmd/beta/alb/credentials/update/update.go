package update

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"golang.org/x/term"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const passwordEnv = "ALB_CREDENTIALS_PASSWORD" //nolint:gosec // false alert, these are not valid credentials

const (
	usernameFlag     = "username"
	displaynameFlag  = "displayname"
	credentialRefArg = "CREDENTIAL_REF_ARG" //nolint:gosec // false alert, these are not valid credentials
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Username       *string
	Displayname    *string
	CredentialsRef *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", credentialRefArg),
		Short: "Update credentials",
		Long:  "Update credentials.",
		Args:  args.SingleArg(credentialRefArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update the username`,
				"$ stackit beta alb credentials update --username test-cred2 credentials-12345",
			),
			examples.NewExample(
				`Update the displayname`,
				"$ stackit beta alb credentials update --displayname new-name credentials-12345",
			),
			examples.NewExample(
				`Update the password (is retrieved interactively or from ENV variable )`,
				"$ stackit beta alb credentials update --password credentials-12345",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model := parseInput(p, cmd, args)

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			req, err := buildRequest(ctx, &model, apiClient, readPassword)
			if err != nil {
				return err
			}
			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}
			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update credential %q for %q?", *model.CredentialsRef, projectLabel)
				err = p.PromptForConfirmation(prompt)
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

			return outputResult(p, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(usernameFlag, "u", "", "the username for the credentials")
	cmd.Flags().StringP(displaynameFlag, "d", "", "the displayname for the credentials")
	cobra.CheckErr(flags.MarkFlagsRequired(cmd, displaynameFlag, usernameFlag, displaynameFlag))
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient, readPassword func() (string, error)) (req alb.ApiUpdateCredentialsRequest, err error) {
	req = apiClient.UpdateCredentials(ctx, model.ProjectId, model.Region, *model.CredentialsRef)

	var password *string
	p, err := readPassword()
	if err != nil {
		return req, err
	}
	password = &p
	payload := alb.UpdateCredentialsPayload{
		DisplayName: model.Displayname,
		Password:    password,
		Username:    model.Username,
	}

	if model.Displayname == nil && model.Username == nil {
		return req, fmt.Errorf("no attribute to change passed")
	}

	return req.UpdateCredentialsPayload(payload), nil
}
func readPassword() (string, error) {
	if password, found := os.LookupEnv(passwordEnv); found {
		return password, nil
	}

	fmt.Printf("please provide the password: ")
	password, err := term.ReadPassword(int(os.Stdout.Fd()))
	if err != nil {
		return "", fmt.Errorf("cannot read password: %w", err)
	}
	fmt.Println()
	fmt.Printf("please confirm the password: ")
	confirmation, err := term.ReadPassword(int(os.Stdout.Fd()))
	if err != nil {
		return "", fmt.Errorf("cannot read password: %w", err)
	}
	fmt.Println()
	if !bytes.Equal(password, confirmation) {
		return "", fmt.Errorf("the password and the confirmation do not match")
	}

	return string(password), nil
}
func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) inputModel {
	model := inputModel{
		GlobalFlagModel: globalflags.Parse(p, cmd),
		Username:        flags.FlagToStringPointer(p, cmd, usernameFlag),
		Displayname:     flags.FlagToStringPointer(p, cmd, displaynameFlag),
		CredentialsRef:  &inputArgs[0],
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
		p.Outputf("Updated labels of credential %q\n", utils.PtrString(model.CredentialsRef))
	}
	return nil
}
