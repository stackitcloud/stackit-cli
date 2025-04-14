package add

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const (
	usernameFlag    = "username"
	displaynameFlag = "displayname"
	passwordFlag    = "password"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Username    *string
	Displayname *string
	Password    *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Adds observability credentials to an application load balancer",
		Long:  "Adds observability credentials (username and password) to an application load balancer.  The credentials can be for Observability or another monitoring tool.",
		Args:  cobra.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Add observability credentials to a load balancer with username "xxx" and display name "yyy", providing the path to a file with the password as flag`,
				"$ stackit beta alb observability-credentials add --username xxx --password @./password.txt --display-name yyy"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := "Are your sure you want to add credentials?"
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("add credential: %w", err)
			}

			return outputResult(p, model.GlobalFlagModel.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(usernameFlag, "u", "", "Username for the credentials")
	cmd.Flags().StringP(displaynameFlag, "d", "", "Displayname for the credentials")
	cmd.Flags().Var(flags.ReadFromFileFlag(), passwordFlag, `Password. Can be a string or a file path, if prefixed with "@" (example: @./password.txt).`)

	cobra.CheckErr(flags.MarkFlagsRequired(cmd, usernameFlag, displaynameFlag))
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Username:        flags.FlagToStringPointer(p, cmd, usernameFlag),
		Displayname:     flags.FlagToStringPointer(p, cmd, displaynameFlag),
		Password:        flags.FlagToStringPointer(p, cmd, passwordFlag),
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
		if item.Credential != nil {
			p.Outputf("Created credential %s\n",
				utils.PtrString(item.Credential.CredentialsRef),
			)
		}
	}
	return nil
}
