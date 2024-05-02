package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/rabbitmq/client"
	rabbitmqUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/rabbitmq/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/rabbitmq"
)

const (
	instanceIdFlag   = "instance-id"
	hidePasswordFlag = "hide-password"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId   string
	HidePassword bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates credentials for a RabbitMQ instance",
		Long:  "Creates credentials (username and password) for a RabbitMQ instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials for a RabbitMQ instance`,
				"$ stackit rabbitmq credentials create --instance-id xxx"),
			examples.NewExample(
				`Create credentials for a RabbitMQ instance and hide the password in the output`,
				"$ stackit rabbitmq credentials create --instance-id xxx --hide-password"),
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

			instanceLabel, err := rabbitmqUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create credentials for instance %q?", instanceLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create RabbitMQ credentials: %w", err)
			}

			return outputResult(p, model, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Bool(hidePasswordFlag, false, "Hide password in output")

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
		HidePassword:    flags.FlagToBoolValue(p, cmd, hidePasswordFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *rabbitmq.APIClient) rabbitmq.ApiCreateCredentialsRequest {
	req := apiClient.CreateCredentials(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *rabbitmq.CredentialsResponse) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal RabbitMQ credentials: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created credentials for instance %q. Credentials ID: %s\n\n", instanceLabel, *resp.Id)
		// The username field cannot be set by the user so we only display it if it's not returned empty
		username := *resp.Raw.Credentials.Username
		if username != "" {
			p.Outputf("Username: %s\n", *resp.Raw.Credentials.Username)
		}
		if model.HidePassword {
			p.Outputf("Password: <hidden>\n")
		} else {
			p.Outputf("Password: %s\n", *resp.Raw.Credentials.Password)
		}
		p.Outputf("Host: %s\n", *resp.Raw.Credentials.Host)
		p.Outputf("Port: %d\n", *resp.Raw.Credentials.Port)
		p.Outputf("URI: %s\n", *resp.Uri)
		return nil
	}
}
