package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/telemetrylink/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetrylink/utils"

	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetrylink/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	telemetrylink "github.com/stackitcloud/stackit-sdk-go/services/telemetrylink/v1api"

	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
)

const (
	displayNameFlag       = "display-name"
	resourceIdFlag        = "resource-id"
	resourceTypeFlag      = "resource-type"
	descriptionFlag       = "description"
	telemetryRouterIdFlag = "telemetry-router-id"
	accessTokenFlag       = "access-token"
	enabledFlag           = "enabled"
)

var (
	resourceTypes = []string{"organization", "folder", "project"}
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	DisplayName       string
	ResourceId        string
	ResourceType      string
	Description       *string
	TelemetryRouterId string
	AccessToken       string
	Enabled           *bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Telemetry Link with the specified resource",
		Long:  "Creates a Telemetry Link with the specified resource.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Telemetry Link with name "my-link" for project with ID "xxx" and for telemetry router with ID "yyy"`,
				"$ stackit beta telemetrylink create --display-name my-link --resource-type project --resource-id xxx --telemetry-router-id yyy --access-token my-token"),
			examples.NewExample(
				`Create a Telemetry Link with name "my-link" for folder with ID "xxx" and for telemetry router with ID "yyy"`,
				"$ stackit beta telemetrylink create --display-name my-link --resource-type folder --resource-id xxx --telemetry-router-id yyy --access-token my-token"),
			examples.NewExample(
				`Create a Telemetry Link with name "my-link" for organization with ID "xxx" and for telemetry router with ID "yyy"`,
				"$ stackit beta telemetrylink create --display-name my-link --resource-type organization --resource-id xxx --telemetry-router-id yyy --access-token my-token"),
			examples.NewExample(
				`Create a Telemetry Link with name "my-link" for organization with ID "xxx" and for telemetry router with ID "yyy", and disable routing`,
				"$ stackit beta telemetrylink create --display-name my-link --resource-type organization --resource-id xxx --telemetry-router-id yyy --enabled=false --access-token my-token"),
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

			var resourceLabel string
			rmApiClient, err := rmClient.ConfigureClient(params.Printer, params.CliVersion)
			if err == nil {
				resourceLabel, err = utils.GetResourceLabel(ctx, rmApiClient.DefaultAPI, model.ResourceId, model.ResourceType)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get %v name: %v", model.ResourceType, err)
					resourceLabel = model.ResourceId
				}
			} else {
				params.Printer.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
			}

			if resourceLabel == "" {
				resourceLabel = model.ResourceId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a Telemetry Link for %q %q?", model.ResourceType, resourceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			switch model.ResourceType {
			case "project":
				projectReq, errReq := buildProjectRequest(ctx, model, apiClient)
				if errReq != nil {
					return errReq
				}
				_, err = projectReq.Execute()
			case "organization":
				orgReq, errReq := buildOrganizationRequest(ctx, model, apiClient)
				if errReq != nil {
					return errReq
				}
				_, err = orgReq.Execute()
			case "folder":
				folderReq, errReq := buildFolderRequest(ctx, model, apiClient)
				if errReq != nil {
					return errReq
				}
				_, err = folderReq.Execute()
			default:
				params.Printer.Debug(print.ErrorLevel, "unknown resource type: %v", model.ResourceType)
			}

			if err != nil {
				return fmt.Errorf("create Telemetry Link: %w", err)
			}

			var response *telemetrylink.TelemetryLinkResponse
			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating Telemetry Link", func() error {
					switch model.ResourceType {
					case "project":
						response, err = wait.CreateProjectTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					case "organization":
						response, err = wait.CreateOrganizationTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					case "folder":
						response, err = wait.CreateFolderTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					default:
						params.Printer.Debug(print.ErrorLevel, "unknown resource type: %v", model.ResourceType)
					}
					if err != nil {
						params.Printer.Debug(print.ErrorLevel, "error creating telemetry link: %v", err)
					}
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Telemetry Link creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, model.Async, resourceLabel, response)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "The displayed name of the Telemetry Link resource")
	cmd.Flags().String(descriptionFlag, "", "The description of the Telemetry Link resource")
	cmd.Flags().String(accessTokenFlag, "", "The access token of the Telemetry Router instance")
	cmd.Flags().String(resourceTypeFlag, "", fmt.Sprintf("The resource type of the TelemetryLink resource, possible values are %s", resourceTypes))
	cmd.Flags().Var(flags.UUIDFlag(), resourceIdFlag, "STACKIT project ID, folder ID, or organization ID associated with the Telemetry Link resource")
	cmd.Flags().Var(flags.UUIDFlag(), telemetryRouterIdFlag, "The ID of the telemetry-router to route the telemetry data")
	cmd.Flags().Bool(enabledFlag, true, "Enable routing through the link to a telemetry router")

	err := flags.MarkFlagsRequired(cmd, accessTokenFlag, displayNameFlag, telemetryRouterIdFlag, resourceIdFlag, resourceTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		DisplayName:       flags.FlagToStringValue(p, cmd, displayNameFlag),
		Description:       flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Enabled:           flags.FlagToBoolPointer(p, cmd, enabledFlag),
		AccessToken:       flags.FlagToStringValue(p, cmd, accessTokenFlag),
		ResourceType:      flags.FlagToStringValue(p, cmd, resourceTypeFlag),
		ResourceId:        flags.FlagToStringValue(p, cmd, resourceIdFlag),
		TelemetryRouterId: flags.FlagToStringValue(p, cmd, telemetryRouterIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildProjectRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiCreateOrUpdateProjectTelemetryLinkRequest, error) {
	req := apiClient.DefaultAPI.CreateOrUpdateProjectTelemetryLink(ctx, model.ResourceId, model.Region).CreateOrUpdateProjectTelemetryLinkPayload(
		telemetrylink.CreateOrUpdateProjectTelemetryLinkPayload{
			AccessToken:       model.AccessToken,
			Description:       model.Description,
			DisplayName:       model.DisplayName,
			Enabled:           model.Enabled,
			TelemetryRouterId: model.TelemetryRouterId,
		})

	return req, nil
}

func buildOrganizationRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiCreateOrUpdateOrganizationTelemetryLinkRequest, error) {
	req := apiClient.DefaultAPI.CreateOrUpdateOrganizationTelemetryLink(ctx, model.ResourceId, model.Region).CreateOrUpdateOrganizationTelemetryLinkPayload(
		telemetrylink.CreateOrUpdateOrganizationTelemetryLinkPayload{
			AccessToken:       model.AccessToken,
			Description:       model.Description,
			DisplayName:       model.DisplayName,
			Enabled:           model.Enabled,
			TelemetryRouterId: model.TelemetryRouterId,
		})

	return req, nil
}

func buildFolderRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiCreateOrUpdateFolderTelemetryLinkRequest, error) {
	req := apiClient.DefaultAPI.CreateOrUpdateFolderTelemetryLink(ctx, model.ResourceId, model.Region).CreateOrUpdateFolderTelemetryLinkPayload(
		telemetrylink.CreateOrUpdateFolderTelemetryLinkPayload{
			AccessToken:       model.AccessToken,
			Description:       model.Description,
			DisplayName:       model.DisplayName,
			Enabled:           model.Enabled,
			TelemetryRouterId: model.TelemetryRouterId,
		})

	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, async bool, resourceLabel string, resp *telemetrylink.TelemetryLinkResponse) error {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}

	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("input model is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s telemetry link for %q %q. Link ID: %s\n", operationState, model.ResourceType, resourceLabel, resp.Id)
		return nil
	})
}
