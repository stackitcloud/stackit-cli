package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"

	"github.com/spf13/cobra"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

const (
	nameArg = "NAME"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", nameArg),
		Short: "Deletes an ALB WAF configuration",
		Long:  "Deletes a STACKIT Application Load Balancer (ALB) Web Application Firewall (WAF) configuration.",
		Args:  args.SingleArg(nameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Delete an ALB WAF configuration with name "my-waf-config"`,
				"$ stackit beta alb-waf config delete my-waf-config",
			),
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

			prompt := fmt.Sprintf("Are you sure you want to delete the ALB WAF configuration %q for project %q?", model.Name, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete ALB WAF configuration: %w", err)
			}

			params.Printer.Info("WAF configuration %q deleted.\n", model.Name)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	name := inputArgs[0]
	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            name,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) albwaf.ApiDeleteWAFRequest {
	return apiClient.DefaultAPI.DeleteWAF(ctx, model.ProjectId, model.Region, model.Name)
}
