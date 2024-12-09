package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SecurityGroupId string
}

const argNameGroupId = "argGroupId"

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "describe security groups",
		Long:  "describe security groups",
		Args:  args.SingleArg(argNameGroupId, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`describe an existing group`, `$ stackit beta security-group describe 9e9c44fe-eb9a-4d45-bf08-365e961845d1`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return executeDescribe(cmd, p, args)
		},
	}

	return cmd
}

func executeDescribe(cmd *cobra.Command, p *print.Printer, args []string) error {
	p.Info("executing describe command")
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
	request := buildRequest(ctx, model, apiClient)

	p.Info("security group %q for %q\n", model.SecurityGroupId, model.ProjectId)

	group, err := request.Execute()
	if err != nil {
		return fmt.Errorf("get security group: %w", err)
	}
	if err := outputResult(p, model, group); err != nil {
		return err
	}

	return nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetSecurityGroupRequest {
	request := apiClient.GetSecurityGroup(ctx, model.ProjectId, model.SecurityGroupId)
	return request

}

func parseInput(p *print.Printer, cmd *cobra.Command, args []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	if err := cmd.ValidateArgs(args); err != nil {
		return nil, &errors.ArgValidationError{
			Arg:     argNameGroupId,
			Details: fmt.Sprintf("argument validation failed: %v", err),
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SecurityGroupId: args[0],
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

func outputResult(p *print.Printer, model *inputModel, resp *iaas.SecurityGroup) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal security group: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal security group: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		if id := resp.Id; id != nil {
			table.AddRow("ID", *id)
		}
		table.AddSeparator()

		if name := resp.Name; name != nil {
			table.AddRow("NAME", *name)
			table.AddSeparator()
		}

		if description := resp.Description; description != nil {
			table.AddRow("DESCRIPTION", *description)
			table.AddSeparator()
		}

		if labels := resp.Labels; labels != nil {
			var builder strings.Builder
			for k, v := range *labels {
				builder.WriteString(fmt.Sprintf("%s=%s ", k, v))
			}
			table.AddRow("LABELS", builder.String())
			table.AddSeparator()
		}

		if stateful := resp.Stateful; stateful != nil {
			table.AddRow("STATEFUL", *stateful)
			table.AddSeparator()
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
