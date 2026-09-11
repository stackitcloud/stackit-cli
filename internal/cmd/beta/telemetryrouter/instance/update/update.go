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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"
)

const (
	argInstanceID = "INSTANCE_ID"

	displayNameFlag     = "display-name"
	descriptionFlag     = "description"
	filterAttributeFlag = "filter-attribute"
	filterFlag          = "filter"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceID  string
	DisplayName *string
	Description *string
	Filter      *telemetryrouter.ConfigFilter
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", argInstanceID),
		Short: "Updates a TelemetryRouter instance",
		Long:  "Updates a TelemetryRouter instance.",
		Args:  args.SingleArg(argInstanceID, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of the TelemetryRouter instance with ID "xxx"`,
				"$ stackit beta telemetryrouter instance update xxx --display-name new-name"),
			examples.NewExample(
				`Update the description of the TelemetryRouter instance with ID "xxx"`,
				"$ stackit beta telemetryrouter instance update xxx --description new-description"),
			examples.NewExample(
				`Replace the filter attributes of the TelemetryRouter instance with ID "xxx" so it only routes telemetry data for HTTP GET or HEAD requests, matched at the resource level`,
				`$ stackit beta telemetryrouter instance update xxx --filter-attribute "key=http.method;level=resource;matcher==;values=GET,HEAD"`),
			examples.NewExample(
				`Replace the filter attributes of the TelemetryRouter instance with ID "xxx" with many filter attributes sourced from a JSON file (recommended over repeating --filter-attribute for more than a handful of attributes); `+
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

			instanceLabel, err := telemetryrouterUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceID)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceID
			}

			prompt := fmt.Sprintf("Are you sure you want to update instance %s?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)

			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update TelemetryRouter instance: %w", err)
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
	cmd.Flags().StringArray(filterAttributeFlag, nil,
		`Adds a filter attribute that restricts which telemetry data is routed. `+
			`Must be of the form "key=<attribute-key>;level=<resource|scope|logRecord>;matcher=<=|!=>;values=<v1>,<v2>,...", `+
			`e.g. "key=http.method;level=resource;matcher==;values=GET,HEAD" only routes data whose "http.method" resource `+
			`attribute equals "GET" or "HEAD". Can be repeated to add multiple filter attributes, which are combined with AND. `+
			`Replaces any existing filter attributes on the instance. Impractical for many filter attributes at once; `+
			`use "--filter" instead. Mutually exclusive with "--filter". The API does not support clearing all filter `+
			`attributes on update; at least one filter attribute must remain set.`)
	cmd.Flags().Var(flags.ReadFromFileFlag(), filterFlag,
		`Filter configuration as JSON, of the form {"attributes": [{"key": ..., "level": ..., "matcher": ..., "values": [...]}, ...]}. `+
			`Can be a JSON string or a file path, if prefixed with "@" (example: @./filter.json). Replaces any existing filter `+
			`attributes on the instance. Recommended over "--filter-attribute" when setting many filter attributes at once. `+
			`Run "stackit beta telemetryrouter instance generate-filter --instance-id <id>" to generate a starting point from the `+
			`instance's current filter. Mutually exclusive with "--filter-attribute". The API does not support clearing all `+
			`filter attributes on update; at least one filter attribute must remain set.`)

	cmd.MarkFlagsMutuallyExclusive(filterAttributeFlag, filterFlag)
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
	filterAttributes := flags.FlagToStringArrayValue(p, cmd, filterAttributeFlag)

	// --filter-attribute and --filter are mutually exclusive; this is enforced by
	// cmd.MarkFlagsMutuallyExclusive in configureFlags (validated by cobra before parseInput
	// runs), so both being set at once should not happen here in practice.
	var filter *telemetryrouter.ConfigFilter
	switch {
	case filterValue != nil:
		parsedFilter, err := telemetryrouterUtils.ParseFilterJSON(*filterValue)
		if err != nil {
			return nil, err
		}
		filter = parsedFilter
	case len(filterAttributes) > 0:
		builtFilter, err := telemetryrouterUtils.BuildConfigFilter(filterAttributes)
		if err != nil {
			return nil, err
		}
		filter = builtFilter
	}

	if displayName == nil && description == nil && filter == nil {
		return nil, &errors.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceID:      instanceId,
		DisplayName:     displayName,
		Description:     description,
		Filter:          filter,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiUpdateTelemetryRouterRequest {
	req := apiClient.DefaultAPI.UpdateTelemetryRouter(ctx, model.ProjectId, model.Region, model.InstanceID)

	req = req.UpdateTelemetryRouterPayload(telemetryrouter.UpdateTelemetryRouterPayload{
		DisplayName: model.DisplayName,
		Description: model.Description,
		Filter:      model.Filter,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, instance *telemetryrouter.TelemetryRouterResponse) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	}
	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("input model is nil")
	}
	return p.OutputResult(model.OutputFormat, instance, func() error {
		p.Outputf("Updated instance %q for project %q.\n", instance.DisplayName, projectLabel)
		return nil
	})
}
