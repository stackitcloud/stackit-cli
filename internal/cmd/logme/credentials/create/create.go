package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logme/client"
	logmeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/logme/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/logme"
)

const (
	instanceIdFlag   = "instance-id"
	showPasswordFlag = "show-password"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId   string
	ShowPassword bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates credentials for a LogMe instance",
		Long:  "Creates credentials (username and password) for a LogMe instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials for a LogMe instance`,
				"$ stackit logme credentials create --instance-id xxx"),
			examples.NewExample(
				`Create credentials for a LogMe instance and show the password in the output`,
				"$ stackit logme credentials create --instance-id xxx --show-password"),
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

			instanceLabel, err := logmeUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create credentials for instance %q?", instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create LogMe credentials: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.ShowPassword, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().BoolP(showPasswordFlag, "s", false, "Show password in output")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		ShowPassword:    flags.FlagToBoolValue(p, cmd, showPasswordFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logme.APIClient) logme.ApiCreateCredentialsRequest {
	req := apiClient.CreateCredentials(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, showPassword bool, instanceLabel string, resp *logme.CredentialsResponse) error {
	if resp == nil {
		return fmt.Errorf("credentials response is empty")
	}

	if !showPassword && resp.HasRaw() && resp.Raw.Credentials != nil {
		resp.Raw.Credentials.Password = utils.Ptr("hidden")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal LogMe credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal LogMe credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created credentials for instance %q. Credentials ID: %s\n\n", instanceLabel, utils.PtrString(resp.Id))
		// The username field cannot be set by the user so we only display it if it's not returned empty
		if resp.HasRaw() && resp.Raw.Credentials != nil {
			if username := resp.Raw.Credentials.Username; username != nil && *username != "" {
				p.Outputf("Username: %s\n", utils.PtrString(username))
			}
			if !showPassword {
				p.Outputf("Password: <hidden>\n")
			} else {
				p.Outputf("Password: %s\n", utils.PtrString(resp.Raw.Credentials.Password))
			}
			p.Outputf("Host: %s\n", utils.PtrString(resp.Raw.Credentials.Host))
			p.Outputf("Port: %s\n", utils.PtrString(resp.Raw.Credentials.Port))
		}
		p.Outputf("URI: %s\n", utils.PtrString(resp.Uri))
		return nil
	}
}
