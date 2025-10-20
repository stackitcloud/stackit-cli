package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

const (
	instanceIdFlag = "instance-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates credentials for an Observability instance.",
		Long: fmt.Sprintf("%s\n%s",
			"Creates credentials (username and password) for an Observability instance.",
			"The credentials will be generated and included in the response. You won't be able to retrieve the password later."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials for Observability instance with ID "xxx"`,
				"$ stackit observability credentials create --instance-id xxx"),
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

			instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
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
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create credentials for Observability instance: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiCreateCredentialsRequest {
	req := apiClient.CreateCredentials(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat, instanceLabel string, resp *observability.CreateCredentialsResponse) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created credentials for instance %q.\n\n", instanceLabel)

		if resp.Credentials != nil {
			// The username field cannot be set by the user, so we only display it if it's not returned empty
			username := *resp.Credentials.Username
			if username != "" {
				p.Outputf("Username: %s\n", username)
			}

			p.Outputf("Password: %s\n", utils.PtrString(resp.Credentials.Password))
		}
		return nil
	})
}
