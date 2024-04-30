package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logme/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/logme"
)

const (
	instanceIdArg = "INSTANCE_ID"

	aclParameterKey = "sgw_acl"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details  of a LogMe instance",
		Long:  "Shows details  of a LogMe instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a LogMe instance with ID "xxx"`,
				"$ stackit logme instance describe xxx"),
			examples.NewExample(
				`Get details of a LogMe instance with ID "xxx" in a table format`,
				"$ stackit logme instance describe xxx --output-format pretty"),
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
				return fmt.Errorf("read LogMe instance: %w", err)
			}

			return outputResult(p, model.OutputFormat, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *logme.APIClient) logme.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *logme.Instance) error {
	switch outputFormat {
	case print.PrettyOutputFormat:
		table := tables.NewTable()
		table.AddRow("ID", *instance.InstanceId)
		table.AddSeparator()
		table.AddRow("NAME", *instance.Name)
		table.AddSeparator()
		table.AddRow("LAST OPERATION TYPE", *instance.LastOperation.Type)
		table.AddSeparator()
		table.AddRow("LAST OPERATION STATE", *instance.LastOperation.State)
		table.AddSeparator()
		table.AddRow("PLAN ID", *instance.PlanId)
		// Only show ACL if it's present and not empty
		acl := (*instance.Parameters)[aclParameterKey]
		aclStr, ok := acl.(string)
		if ok {
			if aclStr != "" {
				table.AddSeparator()
				table.AddRow("ACL", aclStr)
			}
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	default:
		details, err := json.MarshalIndent(instance, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal LogMe instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	}
}
