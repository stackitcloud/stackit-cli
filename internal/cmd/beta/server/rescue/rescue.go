package rescue

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
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

	imageIdFlag = "image-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	ImageId  *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rescue",
		Short: "Rescues an existing server",
		Long:  "Rescues an existing server.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Rescue an existing server with ID "xxx" using image with ID "yyy" as boot volume`,
				"$ stackit beta server rescue xxx --image-id yyy",
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to rescue server %q?", serverLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("server rescue: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p)
				s.Start("Rescuing server")
				_, err = wait.RescueServerWaitHandler(ctx, apiClient, model.ProjectId, model.ServerId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for server rescuing: %w", err)
				}
				s.Stop()
			}

			operationState := "Rescued"
			if model.Async {
				operationState = "Triggered rescue of"
			}
			p.Info("%s server %q. Image %q is used as temporary boot image\n", operationState, serverLabel, *model.ImageId)

			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), imageIdFlag, "The image ID to be used for a temporary boot volume.")

	err := flags.MarkFlagsRequired(cmd, imageIdFlag)
	cobra.CheckErr(err)
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
		ImageId:         flags.FlagToStringPointer(p, cmd, imageIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiRescueServerRequest {
	req := apiClient.RescueServer(ctx, model.ProjectId, model.ServerId)
	payload := iaas.RescueServerPayload{
		Image: model.ImageId,
	}
	return req.RescueServerPayload(payload)
}
