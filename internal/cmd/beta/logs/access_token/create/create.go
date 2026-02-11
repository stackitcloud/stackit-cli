package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/client"
	logsUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

	"github.com/spf13/cobra"
)

const (
	displayNameFlag = "display-name"
	instanceIdFlag  = "instance-id"
	lifetimeFlag    = "lifetime"
	descriptionFlag = "description"
	permissionsFlag = "permissions"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId  string
	Description *string
	DisplayName string
	Lifetime    *int64
	Permissions []string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a log access token",
		Long:  "Creates a log access token.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a access token with the display name "access-token-1" for the instance "xxx" with read and write permissions`,
				`$ stackit logs access-token create --display-name access-token-1 --instance-id xxx --permissions read,write`,
			),
			examples.NewExample(
				`Create a write only access token with a description`,
				`$ stackit logs access-token create --display-name access-token-2 --instance-id xxx --permissions write --description "Access token for service"`,
			),
			examples.NewExample(
				`Create a read only access token which expires in 30 days`,
				`$ stackit logs access-token create --display-name access-token-3 --instance-id xxx --permissions read --lifetime 30`,
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

			instanceLabel, err := logsUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a access token for the log instance %q in the project %q?", instanceLabel, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create log access-token : %w", err)
			}

			if resp == nil || resp.Id == nil {
				return fmt.Errorf("create log access-token : empty response")
			}

			return outputResult(params.Printer, model.OutputFormat, instanceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the logs instance")
	cmd.Flags().String(displayNameFlag, "", "Display name for the access token")
	cmd.Flags().String(descriptionFlag, "", "Description of the access token")
	cmd.Flags().Int64(lifetimeFlag, 0, "Lifetime of the access token in days [1 - 180]")
	cmd.Flags().StringSlice(permissionsFlag, []string{}, `Permissions of the access token ["read" "write"]`)

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, displayNameFlag, permissionsFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DisplayName:     flags.FlagToStringValue(p, cmd, displayNameFlag),
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Lifetime:        flags.FlagToInt64Pointer(p, cmd, lifetimeFlag),
		Permissions:     flags.FlagToStringSliceValue(p, cmd, permissionsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logs.APIClient) logs.ApiCreateAccessTokenRequest {
	req := apiClient.CreateAccessToken(ctx, model.ProjectId, model.Region, model.InstanceId)

	return req.CreateAccessTokenPayload(logs.CreateAccessTokenPayload{
		Description: model.Description,
		DisplayName: &model.DisplayName,
		Lifetime:    model.Lifetime,
		Permissions: &model.Permissions,
	})
}

func outputResult(p *print.Printer, outputFormat, instanceLabel string, accessToken *logs.AccessToken) error {
	if accessToken == nil {
		return fmt.Errorf("access token cannot be nil")
	}
	return p.OutputResult(outputFormat, accessToken, func() error {
		p.Outputf("Created access token for log instance %q.\n\nID: %s\nToken: %s\n", instanceLabel, utils.PtrValue(accessToken.Id), utils.PtrValue(accessToken.AccessToken))
		return nil
	})
}
