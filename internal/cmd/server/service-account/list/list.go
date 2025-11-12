package list

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serverIdFlag = "server-id"
	limitFlag    = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit    *int64
	ServerId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all attached service accounts for a server",
		Long:  "List all attached service accounts for a server",
		Args:  cobra.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all attached service accounts for a server with ID "xxx"`,
				"$ stackit server service-account list --server-id xxx",
			),
			examples.NewExample(
				`List up to 10 attached service accounts for a server with ID "xxx"`,
				"$ stackit server service-account list --server-id xxx --limit 10",
			),
			examples.NewExample(
				`List all attached service accounts for a server with ID "xxx" in JSON format`,
				"$ stackit server service-account list --server-id xxx --output-format json",
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

			serverName, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, model.Region, model.ServerId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get server name: %v", err)
				serverName = model.ServerId
			} else if serverName == "" {
				serverName = model.ServerId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list service accounts: %w", err)
			}
			serviceAccounts := *resp.Items
			if len(serviceAccounts) == 0 {
				params.Printer.Info("No service accounts found for server %s\n", serverName)
				return nil
			}

			if model.Limit != nil && len(serviceAccounts) > int(*model.Limit) {
				serviceAccounts = serviceAccounts[:int(*model.Limit)]
			}

			return outputResult(params.Printer, model.OutputFormat, model.ServerId, serverName, serviceAccounts)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListServerServiceAccountsRequest {
	req := apiClient.ListServerServiceAccounts(ctx, model.ProjectId, model.Region, model.ServerId)
	return req
}

func outputResult(p *print.Printer, outputFormat, serverId, serverName string, serviceAccounts []string) error {
	return p.OutputResult(outputFormat, serviceAccounts, func() error {
		table := tables.NewTable()
		table.SetHeader("SERVER ID", "SERVER NAME", "SERVICE ACCOUNT")
		for i := range serviceAccounts {
			table.AddRow(serverId, serverName, serviceAccounts[i])
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("rednder table: %w", err)
		}
		return nil
	})
}
