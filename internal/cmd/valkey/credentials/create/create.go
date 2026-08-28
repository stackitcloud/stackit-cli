package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	valkey "github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api"
	"github.com/stackitcloud/stackit-sdk-go/services/valkey/v2api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/client"
	valkeyUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/valkey/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
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
		Short: "Creates credentials for a Valkey instance",
		Long:  "Creates credentials (username and password) for a Valkey instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create credentials for a Valkey instance with ID "xxx"`,
				"$ stackit valkey credentials create --instance-id xxx"),
			examples.NewExample(
				`Create credentials for a Valkey instance and show the password in the output`,
				"$ stackit valkey credentials create --instance-id xxx --show-password"),
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

			instanceLabel, err := valkeyUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.InstanceId, model.Region)
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
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Valkey credentials: %w", err)
			}

			credentialsId := resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating credentials", func() error {
					resp, err = wait.CreateCredentialsWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, credentialsId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Valkey credentials creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, instanceLabel, resp)
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
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		ShowPassword:    flags.FlagToBoolValue(p, cmd, showPasswordFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient valkey.DefaultAPI) valkey.ApiCreateCredentialsRequest {
	// No payload — the empty body in the OpenAPI spec is a known issue
	return apiClient.CreateCredentials(ctx, model.ProjectId, model.Region, model.InstanceId)
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *valkey.CredentialsResponse) error {
	if resp == nil {
		return fmt.Errorf("no response defined")
	}

	if !model.ShowPassword {
		if resp.Raw == nil {
			resp.Raw = &valkey.RawCredentials{Credentials: valkey.Credentials{}}
		}
		resp.Raw.Credentials.Password = "hidden"
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s credentials for instance %q. Credentials ID: %s\n\n", operationState, instanceLabel, resp.Id)
		if resp.Raw != nil {
			if username := resp.Raw.Credentials.Username; username != "" {
				p.Outputf("Username: %s\n", username)
			}
			if !model.ShowPassword {
				p.Outputf("Password: <hidden>\n")
			} else {
				p.Outputf("Password: %s\n", resp.Raw.Credentials.Password)
			}
			p.Outputf("Host: %s\n", resp.Raw.Credentials.Host)
			p.Outputf("Port: %s\n", utils.PtrString(resp.Raw.Credentials.Port))
		}
		p.Outputf("URI: %s\n", resp.Uri)
		return nil
	})
}
