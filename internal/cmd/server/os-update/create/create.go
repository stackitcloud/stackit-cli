package create

import (
	"context"
	"fmt"

	iaasClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverosupdate/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/serverupdate"
)

const (
	serverIdFlag             = "server-id"
	maintenanceWindowFlag    = "maintenance-window"
	defaultMaintenanceWindow = 23
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	ServerId          string
	MaintenanceWindow int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Server os-update.",
		Long:  "Creates a Server os-update. Operation always is async.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Server os-update with name "myupdate"`,
				`$ stackit server os-update create --server-id xxx`),
			examples.NewExample(
				`Create a Server os-update with name "myupdate" and maintenance window for 13 o'clock.`,
				`$ stackit server os-update create --server-id xxx --maintenance-window=13`),
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

			serverLabel := model.ServerId
			// Get server name
			if iaasApiClient, err := iaasClient.ConfigureClient(params.Printer, params.CliVersion); err == nil {
				serverName, err := iaasUtils.GetServerName(ctx, iaasApiClient, model.ProjectId, model.Region, model.ServerId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				} else if serverName != "" {
					serverLabel = serverName
				}
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a os-update for server %s?", serverLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Server os-update: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, serverLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().Int64P(maintenanceWindowFlag, "m", defaultMaintenanceWindow, "Maintenance window (in hours, 1-24)")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		ServerId:          flags.FlagToStringValue(p, cmd, serverIdFlag),
		MaintenanceWindow: flags.FlagWithDefaultToInt64Value(p, cmd, maintenanceWindowFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverupdate.APIClient) (serverupdate.ApiCreateUpdateRequest, error) {
	req := apiClient.CreateUpdate(ctx, model.ProjectId, model.ServerId, model.Region)
	payload := serverupdate.CreateUpdatePayload{
		MaintenanceWindow: &model.MaintenanceWindow,
	}
	req = req.CreateUpdatePayload(payload)
	return req, nil
}

func outputResult(p *print.Printer, outputFormat, serverLabel string, resp serverupdate.Update) error {
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Triggered creation of server os-update for server %s. Update ID: %s\n", serverLabel, utils.PtrString(resp.Id))
		return nil
	})
}
