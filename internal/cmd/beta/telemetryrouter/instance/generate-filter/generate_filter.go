package generatefilter

import (
	"context"
	"encoding/json"
	"fmt"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/fileutils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	telemetryrouterUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
)

const (
	instanceIdFlag = "instance-id"
	filePathFlag   = "file-path"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId *string
	FilePath   *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-filter",
		Short: "Generates a filter configuration to use as --filter input for TelemetryRouter instances",
		Long: fmt.Sprintf("%s\n%s\n%s",
			`Generates a JSON filter configuration to be used as "--filter" input for "instance create" or "instance update".`,
			`If "--instance-id" is set, the current filter configuration of that instance is returned, to be adapted with custom values.`,
			`If unset, a default example filter configuration is returned.`,
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Generate a default example filter configuration, and adapt it with custom values`,
				`$ stackit beta telemetryrouter instance generate-filter --file-path ./filter.json`,
				`<Modify filter in file, if needed>`,
				`$ stackit beta telemetryrouter instance create --display-name "my-instance" --filter @./filter.json`),
			examples.NewExample(
				`Generate a filter configuration with the current values of TelemetryRouter instance "xxx", and adapt it with custom values`,
				`$ stackit beta telemetryrouter instance generate-filter --instance-id xxx --file-path ./filter.json`,
				`<Modify filter in file>`,
				`$ stackit beta telemetryrouter instance update xxx --filter @./filter.json`),
			examples.NewExample(
				`Generate a filter configuration with the current values of TelemetryRouter instance "xxx", and preview it in the terminal`,
				`$ stackit beta telemetryrouter instance generate-filter --instance-id xxx`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			if model.InstanceId == nil {
				defaultFilter := telemetryrouterUtils.DefaultConfigFilter
				return outputResult(params.Printer, model.FilePath, &defaultFilter)
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read TelemetryRouter instance: %w", err)
			}
			if resp == nil {
				return fmt.Errorf("read TelemetryRouter instance: empty response from API")
			}

			filter := resp.Filter
			if filter == nil {
				defaultFilter := telemetryrouterUtils.DefaultConfigFilter
				filter = &defaultFilter
			}

			return outputResult(params.Printer, model.FilePath, filter)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag,
		"If set, generates a filter configuration with the current values of the given TelemetryRouter instance. If unset, generates a default example filter configuration")
	cmd.Flags().StringP(filePathFlag, "f", "", "If set, writes the filter configuration to the given file. If unset, writes it to the standard output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	instanceId := flags.FlagToStringPointer(p, cmd, instanceIdFlag)
	if instanceId != nil && globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
		FilePath:        flags.FlagToStringPointer(p, cmd, filePathFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiGetTelemetryRouterRequest {
	return apiClient.DefaultAPI.GetTelemetryRouter(ctx, model.ProjectId, model.Region, *model.InstanceId)
}

func outputResult(p *print.Printer, filePath *string, filter *telemetryrouter.ConfigFilter) error {
	if filter == nil {
		return fmt.Errorf("filter is nil")
	}

	filterBytes, err := json.MarshalIndent(*filter, "", "  ") // nolint:gosec // false positive
	if err != nil {
		return fmt.Errorf("marshal filter: %w", err)
	}

	if filePath != nil {
		err = fileutils.WriteToFile(*filePath, string(filterBytes))
		if err != nil {
			return fmt.Errorf("write filter to the file: %w", err)
		}
	} else {
		p.Outputln(string(filterBytes))
	}

	return nil
}
