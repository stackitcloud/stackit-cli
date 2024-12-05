package detach

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serviceAccMailArg = "SERVICE_ACCOUNT_MAIL"

	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId       *string
	ServiceAccMail string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detach",
		Short: "Detach a service account from a server",
		Long:  "Detach a service account from a server",
		Args:  args.SingleArg(serviceAccMailArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Detach a service account with mail "xxx@sa.stackit.cloud" from a server "yyy"`,
				"$ stackit beta server service-account detach xxx@sa.stackit.cloud --server-id yyy",
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
			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, *model.ServerId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = *model.ServerId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are your sure you want to detach service account %q from a server %q?", model.ServiceAccMail, serverLabel)
				err = p.PromptForConfirmation(prompt)
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

			return outputResult(p, model.OutputFormat, model.ServiceAccMail, serverLabel, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiRemoveServiceAccountFromServerRequest {
	req := apiClient.RemoveServiceAccountFromServer(ctx, model.ProjectId, *model.ServerId, model.ServiceAccMail)
	return req
}

func outputResult(p *print.Printer, outputFormat, serviceAccMail, serverLabel string, service *iaas.ServiceAccountMailListResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(service, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal service account: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(service, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal service account: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Detached service account %q from server %q\n", serviceAccMail, serverLabel)
		return nil
	}
}
