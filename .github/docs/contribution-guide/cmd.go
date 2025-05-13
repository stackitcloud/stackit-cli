package bar

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"gopkg.in/yaml.v2"
	// (...)
)

// Define consts for command flags
const (
	someArg  = "MY_ARG"
	someFlag = "my-flag"
)

// Struct to model user input (arguments and/or flags)
type inputModel struct {
	*globalflags.GlobalFlagModel
	MyArg  string
	MyFlag *string
}

// "bar" command constructor
func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bar",
		Short: "Short description of the command (is shown in the help of parent command)",
		Long:  "Long description of the command. Can contain some more information about the command usage. It is shown in the help of the current command.",
		Args:  args.SingleArg(someArg, utils.ValidateUUID), // Validate argument, with an optional validation function
		Example: examples.Build(
			examples.NewExample(
				`Do something with command "bar"`,
				"$ stackit foo bar arg-value --my-flag flag-value"),
			//...
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("(...): %w", err)
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				projectLabel = model.ProjectId
			}

			// Check API response "resp" and output accordingly
			if resp.Item == nil {
				params.Printer.Info("(...)", projectLabel)
				return nil
			}
			return outputResult(params.Printer, cmd, model.OutputFormat, instances)
		},
	}

	configureFlags(cmd)
	return cmd
}

// Configure command flags (type, default value, and description)
func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(someFlag, "shorthand", "defaultValue", "My flag description")
}

// Parse user input (arguments and/or flags)
func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	myArg := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		MyArg:           myArg,
		MyFlag:          flags.FlagToStringPointer(p, cmd, someFlag),
	}

	// Write the input model to the debug logs
	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

// Build request to the API
func buildRequest(ctx context.Context, model *inputModel, apiClient *foo.APIClient) foo.ApiListInstancesRequest {
	req := apiClient.GetBar(ctx, model.ProjectId, model.MyArg, someArg)
	return req
}

// Output result based on the configured output format
func outputResult(p *print.Printer, cmd *cobra.Command, outputFormat string, resources []foo.Resource) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resources, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal resource list: %w", err)
		}
		p.Outputln(string(details))
		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.Marshal(resources)
		if err != nil {
			return fmt.Errorf("marshal resource list: %w", err)
		}
		p.Outputln(string(details))
		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATE")
		for i := range resources {
			resource := resources[i]
			table.AddRow(*resource.ResourceId, *resource.Name, *resource.State)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	}
}
