package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	telemetryrouterUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api/wait"
)

const (
	instanceIdArg = "INSTANCE_ID"

	displayNameFlag = "display-name"
	descriptionFlag = "description"
	filterFlag      = "filter"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId  string
	DisplayName *string
	Description *string
	Filter      *telemetryrouter.ConfigFilterPatch
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a Telemetry Router instance",
		Long:  "Updates a Telemetry Router instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of the Telemetry Router instance with ID "xxx"`,
				"$ stackit beta telemetryrouter instance update xxx --display-name new-name"),
			examples.NewExample(
				`Update the description of the Telemetry Router instance with ID "xxx"`,
				"$ stackit beta telemetryrouter instance update xxx --description new-description"),
			examples.NewExample(
				`Replace the filter attributes of the Telemetry Router instance with ID "xxx" with a filter configuration sourced from a JSON file so it only routes telemetry data for HTTP GET or HEAD requests, matched at the resource level; `+
					`use "stackit beta telemetryrouter instance generate-filter --instance-id xxx" to generate a starting point from the instance's current filter`,
				`$ stackit beta telemetryrouter instance update xxx --filter @./filter.json`),
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

			prompt := fmt.Sprintf("Are you sure you want to update instance %s?", instanceLabel)

			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update Telemetry Router instance: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating instance", func() error {
					_, err = wait.PartialUpdateTelemetryRouterWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for Telemetry Router instance update: %w", err)
				}
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)

	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().Var(flags.ReadFromFileFlag(), filterFlag,
		`Filter configuration as JSON, of the form {"attributes": [{"key": ..., "level": ..., "matcher": ..., "values": [...]}, ...], "disabled": ...}. `+
			`Can be a JSON string or a file path, if prefixed with "@" (example: @./filter.json). Only the provided fields are updated: `+
			`if "attributes" is set, it replaces any existing filter attributes on the instance; if "disabled" is set, it enables or disables `+
			`the filter without changing its attributes (example: {"disabled": true}). Run "stackit beta telemetryrouter instance generate-filter --instance-id <id>" `+
			`to generate a starting point from the instance's current filter. Clearing all filter attributes is not supported; `+
			`an empty "attributes" list leaves the existing attributes unchanged.`)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	displayName := flags.FlagToStringPointer(p, cmd, displayNameFlag)
	description := flags.FlagToStringPointer(p, cmd, descriptionFlag)
	filterValue := flags.FlagToStringPointer(p, cmd, filterFlag)

	var filter *telemetryrouter.ConfigFilterPatch
	if filterValue != nil {
		parsedFilter, err := telemetryrouterUtils.ParseFilterPatchJSON(*filterValue)
		if err != nil {
			return nil, err
		}

		filter = parsedFilter
	}

	if displayName == nil && description == nil && filter == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		DisplayName:     displayName,
		Description:     description,
		Filter:          filter,
	}

	p.DebugInputModel(model)

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient telemetryrouter.DefaultAPI) telemetryrouter.ApiPartialUpdateTelemetryRouterRequest {
	req := apiClient.PartialUpdateTelemetryRouter(ctx, model.ProjectId, model.Region, model.InstanceId)

	// Partial update: fields which are not set (nil) are left unchanged on API side
	return req.PartialUpdateTelemetryRouterPayload(telemetryrouter.PartialUpdateTelemetryRouterPayload{
		DisplayName: model.DisplayName,
		Description: model.Description,
		Filter:      model.Filter,
	})
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, instance *telemetryrouter.TelemetryRouterResponse) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}

	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("input model is nil")
	}

	return p.OutputResult(model.OutputFormat, instance, func() error {
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}

		p.Outputf("%s instance %q for project %q.\n", operationState, instance.DisplayName, projectLabel)

		return nil
	})
}
