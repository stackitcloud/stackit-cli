package disable

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"

	"github.com/spf13/cobra"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("disable %s", instanceIdArg),
		Short: "Disables public read access for Grafana on Observability instances",
		Long: fmt.Sprintf("%s\n%s",
			"Disables public read access for Grafana on Observability instances.",
			"When disabled, a login is required to access the Grafana dashboards of the instance. Otherwise, anyone can access the dashboards.",
		),
		Args: args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Disable public read access for Grafana on an Observability instance with ID "xxx"`,
				"$ stackit observability grafana public-read-access disable xxx"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
			if err != nil || instanceLabel == "" {
				instanceLabel = model.InstanceId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to disable Grafana public read access for instance %q?", instanceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("build request: %w", err)
			}
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("disable grafana public read access: %w", err)
			}

			params.Printer.Info("Disabled Grafana public read access for instance %q\n", instanceLabel)
			return nil
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

func buildRequest(ctx context.Context, model *inputModel, apiClient observabilityUtils.ObservabilityClient) (observability.ApiUpdateGrafanaConfigsRequest, error) {
	req := apiClient.UpdateGrafanaConfigs(ctx, model.InstanceId, model.ProjectId)
	payload, err := observabilityUtils.GetPartialUpdateGrafanaConfigsPayload(ctx, apiClient, model.InstanceId, model.ProjectId, nil, utils.Ptr(false))
	if err != nil {
		return req, fmt.Errorf("build request payload: %w", err)
	}
	req = req.UpdateGrafanaConfigsPayload(*payload)
	return req, nil
}
