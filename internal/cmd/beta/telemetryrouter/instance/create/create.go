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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	displayNameFlag     = "display-name"
	descriptionFlag     = "description"
	filterAttributeFlag = "filter-attribute"
	filterFlag          = "filter"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	DisplayName *string
	Description *string
	Filter      *telemetryrouter.ConfigFilter
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a TelemetryRouter instance",
		Long:  "Creates a TelemetryRouter instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a TelemetryRouter instance with name "my-instance"`,
				`$ stackit beta telemetryrouter instance create --display-name "my-instance"`),
			examples.NewExample(
				`Create a TelemetryRouter instance with name "my-instance" and a description`,
				`$ stackit beta telemetryrouter instance create --display-name "my-instance" --description "Description of the instance"`),
			examples.NewExample(
				`Create a TelemetryRouter instance with name "my-instance" that only routes telemetry data for HTTP GET or HEAD requests, matched at the resource level`,
				`$ stackit beta telemetryrouter instance create --display-name "my-instance" --filter-attribute "key=http.method;level=resource;matcher==;values=GET,HEAD"`),
			examples.NewExample(
				`Create a TelemetryRouter instance with name "my-instance" and many filter attributes sourced from a JSON file (recommended over repeating --filter-attribute for more than a handful of attributes); `+
					`use "stackit beta telemetryrouter instance generate-filter" to generate a starting point for "./filter.json"`,
				`$ stackit beta telemetryrouter instance create --display-name "my-instance" --filter @./filter.json`),
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
			}
			if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a TelemetryRouter instance for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)

			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create TelemetryRouter instance: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("create TelemetryRouter instance: empty response from API")
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating instance", func() error {
					_, err = wait.CreateTelemetryRouterWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, resp.Id).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for TelemetryRouter instance creation: %w", err)
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
	cmd.Flags().StringArray(filterAttributeFlag, nil,
		`Adds a filter attribute that restricts which telemetry data is routed. `+
			`Must be of the form "key=<attribute-key>;level=<resource|scope|logRecord>;matcher=<=|!=>;values=<v1>,<v2>,...", `+
			`e.g. "key=http.method;level=resource;matcher==;values=GET,HEAD" only routes data whose "http.method" resource `+
			`attribute equals "GET" or "HEAD". Can be repeated to add multiple filter attributes, which are combined with AND. `+
			`Impractical for many filter attributes at once; use "--filter" instead. Mutually exclusive with "--filter".`)
	cmd.Flags().Var(flags.ReadFromFileFlag(), filterFlag,
		`Filter configuration as JSON, of the form {"attributes": [{"key": ..., "level": ..., "matcher": ..., "values": [...]}, ...]}. `+
			`Can be a JSON string or a file path, if prefixed with "@" (example: @./filter.json). `+
			`Recommended over "--filter-attribute" when setting many filter attributes at once. `+
			`Run "stackit beta telemetryrouter instance generate-filter" to generate a starting point. Mutually exclusive with "--filter-attribute".`)

	cmd.MarkFlagsMutuallyExclusive(filterAttributeFlag, filterFlag)

	err := flags.MarkFlagsRequired(cmd, displayNameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	filter, err := parseFilter(p, cmd)
	if err != nil {
		return nil, err
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		DisplayName:     flags.FlagToStringPointer(p, cmd, displayNameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Filter:          filter,
	}

	p.DebugInputModel(model)
	return &model, nil
}

// parseFilter builds the Filter field from either --filter (JSON) or --filter-attribute
// (repeatable key=...;level=...;matcher=...;values=... syntax). The two flags are mutually
// exclusive (enforced by cobra), so at most one of them yields a non-nil result here.
func parseFilter(p *print.Printer, cmd *cobra.Command) (*telemetryrouter.ConfigFilter, error) {
	if filterValue := flags.FlagToStringPointer(p, cmd, filterFlag); filterValue != nil {
		return telemetryrouterUtils.ParseFilterJSON(*filterValue)
	}
	return telemetryrouterUtils.BuildConfigFilter(flags.FlagToStringArrayValue(p, cmd, filterAttributeFlag))
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiCreateTelemetryRouterRequest {
	req := apiClient.DefaultAPI.CreateTelemetryRouter(ctx, model.ProjectId, model.Region)
	req = req.CreateTelemetryRouterPayload(telemetryrouter.CreateTelemetryRouterPayload{
		DisplayName: utils.PtrString(model.DisplayName),
		Description: model.Description,
		Filter:      model.Filter,
	})
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *telemetryrouter.TelemetryRouterResponse) error {
	if resp == nil {
		return fmt.Errorf("create TelemetryRouter instance response is empty")
	}
	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("input model is nil")
	}

	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s instance for project %q. Instance ID: %s\n", operationState, projectLabel, resp.Id)
		return nil
	})
}
