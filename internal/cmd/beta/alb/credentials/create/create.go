package create

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"golang.org/x/term"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const passwordEnv = "ALB_CREDENTIALS_PASSWORD"

const (
	usernameFlag    = "username"
	displaynameFlag = "displayname"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Username    *string
	Displayname *string
	Password    *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a credential",
		Long:  "Creates a credential.",
		Args:  cobra.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new credential, the password is requested interactively or read from ENV variable `+passwordEnv,
				"$ stackit beta alb credential create --username some.user --displayname master-creds",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			password, err := readPassword()
			if err != nil {
				return err
			}
			model, err := parseInput(p, cmd, password)
			if err != nil {
				return err
			}

			// Configure client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := "Are your sure you want to create a credential?"
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create credential: %w", err)
			}

			return outputResult(p, model.GlobalFlagModel.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(usernameFlag, "u", "", "the username for the credentials")
	cmd.Flags().StringP(displaynameFlag, "d", "", "the displayname for the credentials")

	cobra.CheckErr(cmd.MarkFlagRequired(usernameFlag))
	cobra.CheckErr(cmd.MarkFlagRequired(displaynameFlag))
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
	if bytes.Equal(password, confirmation) {
		return "", fmt.Errorf("the password and the confirmation do not match")
	}

	return string(password), nil
}

func parseInput(p *print.Printer, cmd *cobra.Command, password string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Username:        flags.FlagToStringPointer(p, cmd, usernameFlag),
		Displayname:     flags.FlagToStringPointer(p, cmd, displaynameFlag),
		Password:        &password,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string fo debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) alb.ApiCreateCredentialsRequest {
	req := apiClient.CreateCredentials(ctx, model.ProjectId, model.Region)
	payload := alb.CreateCredentialsPayload{
		DisplayName: model.Displayname,
		Password:    model.Password,
		Username:    model.Username,
	}
	return req.CreateCredentialsPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, item *alb.CreateCredentialsResponse) error {
	if item == nil {
		return fmt.Errorf("no credential found")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal credential: %w", err)
		}
		p.Outputln(string(details))
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(item, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal credential: %w", err)
		}
		p.Outputln(string(details))
	default:
		p.Outputf("Created credential %q",
			utils.PtrString(item.Credential),
		)
	}
	return nil
}
