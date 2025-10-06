package detach

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serviceAccMailArg = "SERVICE_ACCOUNT_EMAIL"

	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId       *string
	ServiceAccMail string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("detach %s", serviceAccMailArg),
		Short: "Detach a service account from a server",
		Long:  "Detach a service account from a server",
		Args:  args.SingleArg(serviceAccMailArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Detach a service account with mail "xxx@sa.stackit.cloud" from a server "yyy"`,
				"$ stackit server service-account detach xxx@sa.stackit.cloud --server-id yyy",
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
			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, *model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = *model.ServerId
			} else if serverLabel == "" {
				serverLabel = *model.ServerId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are your sure you want to detach service account %q from a server %q?", model.ServiceAccMail, serverLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("detach service account request: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.ServiceAccMail, serverLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server id")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	serviceAccMail := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringPointer(p, cmd, serverIdFlag),
		ServiceAccMail:  serviceAccMail,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiRemoveServiceAccountFromServerRequest {
	req := apiClient.RemoveServiceAccountFromServer(ctx, model.ProjectId, *model.ServerId, model.ServiceAccMail)
	return req
}

func outputResult(p *print.Printer, outputFormat, serviceAccMail, serverLabel string, service iaas.ServiceAccountMailListResponse) error {
	return p.OutputResult(outputFormat, service, func() error {
		p.Outputf("Detached service account %q from server %q\n", serviceAccMail, serverLabel)
		return nil
	})
}
