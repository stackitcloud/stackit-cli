package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetrylink/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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
		Use:   "describe",
		Short: "Shows details of a Telemetry Link within the specified resource",
		Long:  "Shows details of a Telemetry Link within the specified resource.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Telemetry Link for project with ID "xxx"`,
				"$ stackit beta telemetrylink describe --resource-type project --resource-id xxx"),
			examples.NewExample(
				`Get details of a Telemetry Link for folder with ID "xxx"`,
				"$ stackit beta telemetrylink describe --resource-type folder --resource-id xxx"),
			examples.NewExample(
				`Get details of a Telemetry Link for organization with ID "xxx"`,
				"$ stackit beta telemetrylink describe --resource-type organization --resource-id xxx"),
			examples.NewExample(
				`Get details of a Telemetry Link for organization with ID "xxx" in JSON format`,
				"$ stackit beta telemetrylink describe --resource-type organization --resource-id xxx --output-format json"),
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

			// Call API
			var response *telemetrylink.TelemetryLinkResponse
			switch model.ResourceType {
			case "project":
				projectReq, errReq := buildProjectRequest(ctx, model, apiClient)
				if errReq != nil {
					return errReq
				}
				response, err = projectReq.Execute()
			case "organization":
				orgReq, errReq := buildOrganizationRequest(ctx, model, apiClient)
				if errReq != nil {
					return errReq
				}
				response, err = orgReq.Execute()
			case "folder":
				folderReq, errReq := buildFolderRequest(ctx, model, apiClient)
				if errReq != nil {
					return errReq
				}
				response, err = folderReq.Execute()
			default:
				params.Printer.Debug(print.ErrorLevel, "unknown resource type: %v", model.ResourceType)
			}

			if err != nil {
				return fmt.Errorf("describe Telemetry Link: %w", err)
			}

			return outputResult(params.Printer, model, response)
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

func buildProjectRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiGetProjectTelemetryLinkRequest, error) {
	return apiClient.DefaultAPI.GetProjectTelemetryLink(ctx, model.ResourceId, model.Region), nil
}

func buildOrganizationRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiGetOrganizationTelemetryLinkRequest, error) {
	return apiClient.DefaultAPI.GetOrganizationTelemetryLink(ctx, model.ResourceId, model.Region), nil
}

func buildFolderRequest(ctx context.Context, model *inputModel, apiClient *telemetrylink.APIClient) (telemetrylink.ApiGetFolderTelemetryLinkRequest, error) {
	return apiClient.DefaultAPI.GetFolderTelemetryLink(ctx, model.ResourceId, model.Region), nil
}

func outputResult(p *print.Printer, model *inputModel, resp *telemetrylink.TelemetryLinkResponse) error {
	if resp == nil {
		return fmt.Errorf("received nil response, could not display details")
	}

	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("input model is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		table := tables.NewTable()
		table.AddRow("ID", resp.Id)
		table.AddSeparator()
		table.AddRow("NAME", resp.DisplayName)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(resp.Description))
		table.AddSeparator()
		table.AddRow("TELEMETRY ROUTER ID", resp.TelemetryRouterId)
		table.AddSeparator()
		table.AddRow("STATUS", resp.Status)
		table.AddSeparator()
		table.AddRow("ENABLED", resp.Enabled)
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
