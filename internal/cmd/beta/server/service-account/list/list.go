package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all attached service accounts from a server",
		Long:  "List all attached service accounts from a server",
		Args:  cobra.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all attached service accounts for a server with ID "xxx"`,
				"$ stackit beta server service-account list --server-id xxx",
			),
			examples.NewExample(
				`List all attached service accounts for a server with ID "xxx" in JSON format`,
				"$ stackit beta server service-account list --server-id xxx --output-format json",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			serverName, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, *model.ServerId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get server name: %v", err)
				serverName = *model.ServerId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list service accounts: %w", err)
			}
			serviceAccounts := *resp.Items
			if len(serviceAccounts) == 0 {
				p.Info("No service accounts found for server %s\n", *model.ServerId)
				return nil
			}

			return outputResult(p, model.OutputFormat, *model.ServerId, serverName, serviceAccounts)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

	err := flags.MarkFlagsRequired(cmd, serverIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        flags.FlagToStringPointer(p, cmd, serverIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListServerServiceAccountsRequest {
	req := apiClient.ListServerServiceAccounts(ctx, model.ProjectId, *model.ServerId)
	return req
}

func outputResult(p *print.Printer, outputFormat, serverId, serverName string, serviceAccounts []string) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(serviceAccounts, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal service accounts list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(serviceAccounts, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal service accounts list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
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
	}
}
