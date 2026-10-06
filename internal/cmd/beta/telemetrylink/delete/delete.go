package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-sdk-go/services/telemetrylink/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetrylink/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	telemetrylink "github.com/stackitcloud/stackit-sdk-go/services/telemetrylink/v1api"
)

const (
	resourceIdFlag   = "resource-id"
	resourceTypeFlag = "resource-type"
)

var (
	resourceTypes = []string{"organization", "folder", "project"}
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ResourceId   string
	ResourceType string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Deletes the given Telemetry Link within the specified resource",
		Long:  "Deletes the given Telemetry Link within the specified resource.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Delete the Telemetry Link for project with ID "xxx"`,
				"$ stackit beta telemetrylink delete --resource-type project --resource-id xxx"),
			examples.NewExample(
				`Delete the Telemetry Link for folder with ID "xxx"`,
				"$ stackit beta telemetrylink delete --resource-type folder --resource-id xxx"),
			examples.NewExample(
				`Delete the Telemetry Link for organization with ID "xxx"`,
				"$ stackit beta telemetrylink delete --resource-type organization --resource-id xxx"),
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
				switch model.ResourceType {
				case "project":
					resourceLabel, err = rmUtils.GetProjectName(ctx, rmApiClient.DefaultAPI, model.ResourceId)
				case "organization":
					resourceLabel, err = rmUtils.GetOrganizationName(ctx, rmApiClient.DefaultAPI, model.ResourceId)
				case "folder":
					resourceLabel, err = rmUtils.GetFolderName(ctx, rmApiClient.DefaultAPI, model.ResourceId)
				default:
					params.Printer.Debug(print.ErrorLevel, "unknown resource type: %v", model.ResourceType)
				}
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get %v name: %v", model.ResourceType, err)
					resourceLabel = model.ResourceId
				}
			} else {
				params.Printer.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
			}

			prompt := fmt.Sprintf("Are you sure you want to delete the Telemetry Link for %q %q?", model.ResourceType, resourceLabel)
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
				err = projectReq.Execute()
			case "organization":
				orgReq, errReq := buildOrganizationRequest(ctx, model, apiClient)
				if errReq != nil {
					return errReq
				}
				err = orgReq.Execute()
			case "folder":
				folderReq, errReq := buildFolderRequest(ctx, model, apiClient)
				if errReq != nil {
					return errReq
				}
				err = folderReq.Execute()
			default:
				params.Printer.Debug(print.ErrorLevel, "unknown resource type: %v", model.ResourceType)
			}

			if err != nil {
				return fmt.Errorf("delete Telemetry Link: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Deleting Telemetry Link", func() error {
					switch model.ResourceType {
					case "project":
						_, err = wait.DeleteProjectTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					case "organization":
						_, err = wait.DeleteOrganizationTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					case "folder":
						_, err = wait.DeleteFolderTelemetryLinkWaitHandler(ctx, apiClient.DefaultAPI, model.ResourceId, model.Region).WaitWithContext(ctx)
					default:
						params.Printer.Debug(print.ErrorLevel, "unknown resource type: %v", model.ResourceType)
					}
					if err != nil {
						params.Printer.Debug(print.ErrorLevel, "error deleting telemetry link: %v", err)
					}
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Telemetry Link deletion: %w", err)
				}
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			params.Printer.Outputf("%s telemetry link for %q %q. \n", operationState, model.ResourceType, resourceLabel)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(resourceTypeFlag, "", fmt.Sprintf("The resource type of the TelemetryLink resource, possible values are %s", resourceTypes))
	cmd.Flags().Var(flags.UUIDFlag(), resourceIdFlag, "STACKIT project ID, folder ID, or organization ID associated with the Telemetry Link resource")

	err := flags.MarkFlagsRequired(cmd, resourceIdFlag, resourceTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.Region == "" {
		return nil, &cliErr.RegionError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ResourceType:    flags.FlagToStringValue(p, cmd, resourceTypeFlag),
		ResourceId:      flags.FlagToStringValue(p, cmd, resourceIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildProjectRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiDeleteProjectTelemetryLinkRequest, error) {
	return apiClient.DefaultAPI.DeleteProjectTelemetryLink(ctx, model.ResourceId, model.Region), nil
}

func buildOrganizationRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiDeleteOrganizationTelemetryLinkRequest, error) {
	return apiClient.DefaultAPI.DeleteOrganizationTelemetryLink(ctx, model.ResourceId, model.Region), nil
}

func buildFolderRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiDeleteFolderTelemetryLinkRequest, error) {
	return apiClient.DefaultAPI.DeleteFolderTelemetryLink(ctx, model.ResourceId, model.Region), nil
}
