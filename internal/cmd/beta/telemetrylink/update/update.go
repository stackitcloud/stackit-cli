package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/telemetrylink/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetrylink/utils"

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

	DisplayName       *string
	ResourceId        string
	ResourceType      string
	Description       *string
	TelemetryRouterId *string
	AccessToken       *string
	Enabled           *bool
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates a Telemetry Link within the specified resource",
		Long:  "Updates a Telemetry Link within the specified resource.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Update a Telemetry Link with name "my-new-link" for project with ID "xxx"`,
				"$ stackit beta telemetrylink update --display-name my-new-link --resource-type project --resource-id xxx"),
			examples.NewExample(
				`Update a Telemetry Link with new access token for folder with ID "xxx"`,
				"$ stackit beta telemetrylink update --resource-type folder --resource-id xxx --access-token my-new-token"),
			examples.NewExample(
				`Update a Telemetry Link with new telemetry router ID "xxx" for organization with ID "yyy"`,
				"$ stackit beta telemetrylink update --resource-type organization --resource-id yyy --telemetry-router-id xxx"),
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

			prompt := fmt.Sprintf("Are you sure you want to update the Telemetry Link for %q %q?", model.ResourceType, resourceLabel)
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
				return fmt.Errorf("update Telemetry Link: %w", err)
			}

			var response *telemetrylink.TelemetryLinkResponse
			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating Telemetry Link", func() error {
					switch model.ResourceType {
					case "project":
						response, err = wait.PartialUpdateProjectTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					case "organization":
						response, err = wait.PartialUpdateOrganizationTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					case "folder":
						response, err = wait.PartialUpdateFolderTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					default:
						params.Printer.Debug(print.ErrorLevel, "unknown resource type: %v", model.ResourceType)
					}
					if err != nil {
						params.Printer.Debug(print.ErrorLevel, "error updating telemetry link: %v", err)
					}
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Telemetry Link update: %w", err)
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

	err := flags.MarkFlagsRequired(cmd, resourceIdFlag, resourceTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		DisplayName:       flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:       flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Enabled:           flags.FlagToBoolPointer(p, cmd, enabledFlag),
		AccessToken:       flags.FlagToStringPointer(p, cmd, accessTokenFlag),
		ResourceType:      flags.FlagToStringValue(p, cmd, resourceTypeFlag),
		ResourceId:        flags.FlagToStringValue(p, cmd, resourceIdFlag),
		TelemetryRouterId: flags.FlagToStringPointer(p, cmd, telemetryRouterIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildProjectRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiPartialUpdateProjectTelemetryLinkRequest, error) {
	req := apiClient.DefaultAPI.PartialUpdateProjectTelemetryLink(ctx, model.ResourceId, model.Region).PartialUpdateProjectTelemetryLinkPayload(
		telemetrylink.PartialUpdateProjectTelemetryLinkPayload{
			AccessToken:       model.AccessToken,
			Description:       model.Description,
			DisplayName:       model.DisplayName,
			Enabled:           model.Enabled,
			TelemetryRouterId: model.TelemetryRouterId,
		})

	return req, nil
}

func buildOrganizationRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiPartialUpdateOrganizationTelemetryLinkRequest, error) {
	req := apiClient.DefaultAPI.PartialUpdateOrganizationTelemetryLink(ctx, model.ResourceId, model.Region).PartialUpdateOrganizationTelemetryLinkPayload(
		telemetrylink.PartialUpdateOrganizationTelemetryLinkPayload{
			AccessToken:       model.AccessToken,
			Description:       model.Description,
			DisplayName:       model.DisplayName,
			Enabled:           model.Enabled,
			TelemetryRouterId: model.TelemetryRouterId,
		})

	return req, nil
}

func buildFolderRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiPartialUpdateFolderTelemetryLinkRequest, error) {
	req := apiClient.DefaultAPI.PartialUpdateFolderTelemetryLink(ctx, model.ResourceId, model.Region).PartialUpdateFolderTelemetryLinkPayload(
		telemetrylink.PartialUpdateFolderTelemetryLinkPayload{
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
		operationState := "Updated"
		if async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s telemetry link for %q %q. Link ID: %s\n", operationState, model.ResourceType, resourceLabel, resp.Id)
		return nil
	})
}
