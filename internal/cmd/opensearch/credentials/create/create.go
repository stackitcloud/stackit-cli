package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	opensearch "github.com/stackitcloud/stackit-sdk-go/services/opensearch/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/opensearch/client"
	opensearchUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/opensearch/utils"
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

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates credentials for an OpenSearch instance",
		Long:  "Creates credentials (username and password) for an OpenSearch instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials for an OpenSearch instance`,
				"$ stackit opensearch credentials create --instance-id xxx"),
			examples.NewExample(
				`Create credentials for an OpenSearch instance and show the password in the output`,
				"$ stackit opensearch credentials create --instance-id xxx --show-password"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			instanceLabel, err := opensearchUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to create credentials for instance %q?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create OpenSearch credentials: %w", err)
			}

			return outputResult(params.Printer, *model, instanceLabel, resp)
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

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *opensearch.APIClient) opensearch.ApiCreateCredentialsRequest {
	req := apiClient.DefaultAPI.CreateCredentials(ctx, model.ProjectId, model.Region, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, model inputModel, instanceLabel string, resp *opensearch.CredentialsResponse) error {
	if model.GlobalFlagModel == nil {
		return fmt.Errorf("no global flags defined")
	}
	if resp == nil || resp.Raw == nil {
		return fmt.Errorf("response or response content is nil")
	}

	if !model.ShowPassword {
		resp.Raw.Credentials.Password = "hidden"
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Created credentials for instance %q. Credentials ID: %s\n\n", instanceLabel, resp.Id)
		// The username field cannot be set by the user so we only display it if it's not returned empty
		if resp.HasRaw() {
			if username := resp.Raw.Credentials.Username; username != "" {
				p.Outputf("Username: %s\n", username)
			}
			if !model.ShowPassword {
				p.Outputf("Password: <hidden>\n")
			} else {
				p.Outputf("Password: %s\n", resp.Raw.Credentials.Password)
			}
			p.Outputf("Host: %s\n", resp.Raw.Credentials.Host)
			p.Outputf("Port: %d\n", resp.Raw.Credentials.Port)
		}
		p.Outputf("URI: %s\n", resp.Uri)
		return nil
	})
}
