package deallocate

import (
	"context"
	"fmt"

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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("deallocate %s", serverIdArg),
		Short: "Deallocates an existing server",
		Long:  "Deallocates an existing server.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Deallocate an existing server with ID "xxx"`,
				"$ stackit beta server deallocate xxx",
			),
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

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, model.ServerId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = model.ServerId
			}
			if serverLabel == "" {
				serverLabel = model.ServerId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to deallocate server %q?", serverLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("server deallocate: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Deallocating server")
				_, err = wait.DeallocateServerWaitHandler(ctx, apiClient, model.ProjectId, model.ServerId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for server deallocating: %w", err)
				}
				s.Stop()
			}

			operationState := "Deallocated"
			if model.Async {
				operationState = "Triggered deallocation of"
			}
			p.Info("%s server %q\n", operationState, serverLabel)

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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeallocateServerRequest {
	return apiClient.DeallocateServer(ctx, model.ProjectId, model.ServerId)
}
