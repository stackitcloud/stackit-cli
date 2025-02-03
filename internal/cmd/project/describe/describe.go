package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

const (
	includeParentsFlag = "include-parents"

	projectIdArg = "PROJECT_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ArgProjectId   string
	IncludeParents bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Shows details of a STACKIT project",
		Long:  "Shows details of a STACKIT project.",
		Args:  args.SingleOptionalArg(projectIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get the details of the configured STACKIT project`,
				"$ stackit project describe"),
			examples.NewExample(
				`Get the details of a STACKIT project by explicitly providing the project ID`,
				"$ stackit project describe xxx"),
			examples.NewExample(
				`Get the details of the configured STACKIT project, including details of the parent resources`,
				"$ stackit project describe --include-parents"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read project details: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(includeParentsFlag, false, "When true, the details of the parent resources will be included in the output")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	var projectId string
	if len(inputArgs) > 0 {
		projectId = inputArgs[0]
	}

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" && projectId == "" {
		return nil, fmt.Errorf("Project ID needs to be provided either as an argument or as a flag")
	}

	if projectId == "" {
		projectId = globalFlags.ProjectId
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ArgProjectId:    projectId,
		IncludeParents:  flags.FlagToBoolValue(p, cmd, includeParentsFlag),
	}

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

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) resourcemanager.ApiGetProjectRequest {
	req := apiClient.GetProject(ctx, model.ArgProjectId)
	req.IncludeParents(model.IncludeParents)
	return req
}

func outputResult(p *print.Printer, outputFormat string, project *resourcemanager.GetProjectResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(project, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal project details: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(project, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal project details: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(project.ProjectId))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(project.Name))
		table.AddSeparator()
		table.AddRow("CREATION", utils.PtrString(project.CreationTime))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(project.LifecycleState))
		table.AddSeparator()
		table.AddRow("PARENT ID", utils.PtrString(project.Parent.Id))
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
