package create

import (
	"context"
	"fmt"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	telemetryrouterUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/utils"

	"github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const (
	instanceIdFlag  = "instance-id"
	displayNameFlag = "display-name"
	descriptionFlag = "description"
	ttlFlag         = "ttl"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceId  string
	DisplayName string
	Description *string
	Ttl         *int32
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a TelemetryRouter access token",
		Long:  "Creates a TelemetryRouter access token.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an access token with the display name "access-token-1" for the TelemetryRouter instance "xxx"`,
				`$ stackit beta telemetryrouter access-token create --display-name access-token-1 --instance-id xxx`,
			),
			examples.NewExample(
				`Create an access token with a description`,
				`$ stackit beta telemetryrouter access-token create --display-name access-token-1 --instance-id xxx --description "Access token for service"`,
			),
			examples.NewExample(
				`Create an access token which expires in 30 days`,
				`$ stackit beta telemetryrouter access-token create --display-name access-token-1 --instance-id xxx --ttl 30`,
			),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			instanceLabel, err := telemetryrouterUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to create an access token for the TelemetryRouter instance %q in project %q?", instanceLabel, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create TelemetryRouter access token: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating access token", func() error {
					_, err = wait.CreateAccessTokenWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, resp.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for TelemetryRouter access token creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the TelemetryRouter instance")
	cmd.Flags().String(displayNameFlag, "", "Display name for the access token")
	cmd.Flags().String(descriptionFlag, "", "Description of the access token")
	cmd.Flags().Int32(ttlFlag, 0, "Time-to-live (TTL) in days for the access token. If not set, the token will not expire")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, displayNameFlag)
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
		DisplayName:     flags.FlagToStringValue(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Ttl:             flags.FlagToInt32Pointer(p, cmd, ttlFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiCreateAccessTokenRequest {
	req := apiClient.DefaultAPI.CreateAccessToken(ctx, model.ProjectId, model.Region, model.InstanceId)
	return req.CreateAccessTokenPayload(telemetryrouter.CreateAccessTokenPayload{
		DisplayName: model.DisplayName,
		Description: model.Description,
		Ttl:         *telemetryrouter.NewNullableInt32(model.Ttl),
	})
}

func outputResult(p *print.Printer, model *inputModel, instanceLabel string, resp *telemetryrouter.CreateAccessTokenResponse) error {
	if resp == nil {
		return fmt.Errorf("create TelemetryRouter access token response is empty")
	}
	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("input model is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s access token for TelemetryRouter instance %q.\n\nID: %s\nToken: %s\n", operationState, instanceLabel, resp.Id, resp.AccessToken)
		return nil
	})
}
