package list

import (
	"context"
	"fmt"

	iaasClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/serverbackup/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/serverbackup"
)

const (
	limitFlag    = "limit"
	serverIdFlag = "server-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	Limit    *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all server backups",
		Long:  "Lists all server backups.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all backups for a server with ID "xxx"`,
				"$ stackit server backup list --server-id xxx"),
			examples.NewExample(
				`List all backups for a server with ID "xxx" in JSON format`,
				"$ stackit server backup list --server-id xxx --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list server backups: %w", err)
			}
			backups := *resp.Items
			if len(backups) == 0 {
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
				params.Printer.Info("No backups found for server %s\n", serverLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(backups) > int(*model.Limit) {
				backups = backups[:*model.Limit]
			}
			return outputResult(params.Printer, model.OutputFormat, backups)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().VarP(flags.UUIDFlag(), serverIdFlag, "s", "Server ID")

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
		ServerId:        flags.FlagToStringValue(p, cmd, serverIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *serverbackup.APIClient) serverbackup.ApiListBackupsRequest {
	req := apiClient.ListBackups(ctx, model.ProjectId, model.ServerId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, backups []serverbackup.Backup) error {
	return p.OutputResult(outputFormat, backups, func() error {
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "SIZE (GB)", "STATUS", "CREATED AT", "EXPIRES AT", "LAST RESTORED AT", "VOLUME BACKUPS")
		for i := range backups {
			s := backups[i]

			lastRestored := utils.PtrStringDefault(s.LastRestoredAt, "")
			var volBackups int
			if s.VolumeBackups != nil {
				volBackups = len(*s.VolumeBackups)
			}
			table.AddRow(
				utils.PtrString(s.Id),
				utils.PtrString(s.Name),
				utils.PtrString(s.Size),
				utils.PtrString(s.Status),
				utils.PtrString(s.CreatedAt),
				utils.PtrString(s.ExpireAt),
				lastRestored,
				volBackups,
			)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
