package start

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/spf13/cobra"
)

const (
	serverIdArg = "SERVER_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("start %s", serverIdArg),
		Short: "Starts an existing server or allocates the server if deallocated",
		Long:  "Starts an existing server or allocates the server if deallocated.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Start an existing server with ID "xxx"`,
				"$ stackit server start xxx",
			),
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

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = model.ServerId
			} else if serverLabel == "" {
				serverLabel = model.ServerId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("server start: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Starting server")
				_, err = wait.StartServerWaitHandler(ctx, apiClient, model.ProjectId, model.ServerId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for server starting: %w", err)
				}
				s.Stop()
			}

			operationState := "Started"
			if model.Async {
				operationState = "Triggered start of"
			}
			params.Printer.Info("%s server %q\n", operationState, serverLabel)

			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	serverId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        serverId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiStartServerRequest {
	return apiClient.StartServer(ctx, model.ProjectId, model.ServerId)
}
