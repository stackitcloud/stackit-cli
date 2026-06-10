package attach

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"
)

const (
	serviceAccMailArg = "SERVICE_ACCOUNT_EMAIL" // Deprecated: positional argument is not used anymore, use the flag instead, will be removed after 2026-12

	serverIdFlag = "server-id"

	serviceAccFlag = "service-account-email"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId       string
	ServiceAccMail string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "attach",
		Short: "Attach a service account to a server",
		Long:  "Attach a service account to a server",
		Args:  args.SingleOptionalArg(serviceAccMailArg, nil), // Deprecated: positional argument is not used anymore, use the flag instead, will be removed after 2026-12
		Example: examples.Build(
			examples.NewExample(
				`Attach a service account with mail "xxx@sa.stackit.cloud" to a server with ID "yyy"`,
				"$ stackit server service-account attach --service-account-email xxx@sa.stackit.cloud --server-id yyy",
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
			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = model.ServerId
			} else if serverLabel == "" {
				serverLabel = model.ServerId
			}

			prompt := fmt.Sprintf("Are you sure you want to attach service account %q to server %q?", model.ServiceAccMail, serverLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("attach service account to server: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.ServiceAccMail, serverLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().VarP(flags.EmailFlag(), serviceAccFlag, "a", "Service Account Email")
	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	var serviceAccMail string
	if cmd.Flags().Changed(serviceAccFlag) {
		serviceAccMail = flags.FlagToStringValue(p, cmd, serviceAccFlag)
	} else if len(inputArgs) > 0 {
		serviceAccMail = inputArgs[0]
		p.Warn("using a positional argument for the service account email is deprecated and will be removed after 2026-12. Please use '--%s' instead.\n", serviceAccFlag)
	} else {
		return nil, fmt.Errorf(`service account must be specified by using either the --%s flag or (deprecated) as a positional argument`, serviceAccFlag)
	}

	if serviceAccMail == "" || !strings.Contains(serviceAccMail, "@") {
		return nil, fmt.Errorf("invalid service account email format: %q", serviceAccMail)
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		ServiceAccMail:  serviceAccMail,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiAddServiceAccountToServerRequest {
	req := apiClient.DefaultAPI.AddServiceAccountToServer(ctx, model.ProjectId, model.Region, model.ServerId, model.ServiceAccMail)
	return req
}

func outputResult(p *print.Printer, outputFormat, serviceAccMail, serverLabel string, serviceAccounts iaas.ServiceAccountMailListResponse) error {
	return p.OutputResult(outputFormat, serviceAccounts, func() error {
		p.Outputf("Attached service account %q to server %q\n", serviceAccMail, serverLabel)
		return nil
	})
}
