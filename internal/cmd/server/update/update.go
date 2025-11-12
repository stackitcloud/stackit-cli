package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serverIdArg = "SERVER_ID"

	nameFlag  = "name"
	labelFlag = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	Name     *string
	Labels   *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", serverIdArg),
		Short: "Updates a server",
		Long:  "Updates a server.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update server with ID "xxx" with new name "server-1-new"`,
				`$ stackit server update xxx --name server-1-new`,
			),
			examples.NewExample(
				`Update server with ID "xxx" with new name "server-1-new" and label(s)`,
				`$ stackit server update xxx --name server-1-new --labels key=value,foo=bar`,
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

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, model.Region, model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = model.ServerId
			} else if serverLabel == "" {
				serverLabel = model.ServerId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update server %q?", serverLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update server: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, serverLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Server name")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a server. E.g. '--labels key1=value1,key2=value2,...'")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	serverId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		ServerId:        serverId,
		Labels:          flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateServerRequest {
	req := apiClient.UpdateServer(ctx, model.ProjectId, model.Region, model.ServerId)

	payload := iaas.UpdateServerPayload{
		Name:   model.Name,
		Labels: utils.ConvertStringMapToInterfaceMap(model.Labels),
	}

	return req.UpdateServerPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, serverLabel string, server *iaas.Server) error {
	return p.OutputResult(outputFormat, server, func() error {
		p.Outputf("Updated server %q.\n", serverLabel)
		return nil
	})
}
