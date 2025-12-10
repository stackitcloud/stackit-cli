package log

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	serverIdArg = "SERVER_ID"

	lengthLimitFlag    = "length"
	defaultLengthLimit = 2000 // lines
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
	Length   *int64
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("log %s", serverIdArg),
		Short: "Gets server console log",
		Long:  "Gets server console log.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get server console log for the server with ID "xxx"`,
				"$ stackit server log xxx",
			),
			examples.NewExample(
				`Get server console log for the server with ID "xxx" and limit output lines to 1000`,
				"$ stackit server log xxx --length 1000",
			),
			examples.NewExample(
				`Get server console log for the server with ID "xxx" in JSON format`,
				"$ stackit server log xxx --output-format json",
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("server log: %w", err)
			}

			log := resp.GetOutput()
			lines := strings.Split(log, "\n")

			maxLines := int(*model.Length)
			if len(lines) <= maxLines {
				return outputResult(params.Printer, serverLabel, lines)
			}

			recentLogs := lines[len(lines)-maxLines:]
			return outputResult(params.Printer, serverLabel, recentLogs)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(lengthLimitFlag, defaultLengthLimit, "Maximum number of lines to list")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	serverId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	length := flags.FlagWithDefaultToInt64Value(p, cmd, lengthLimitFlag)
	if length < 0 {
		return nil, &errors.FlagValidationError{
			Flag:    lengthLimitFlag,
			Details: "must not be negative",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        serverId,
		Length:          utils.Ptr(length),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetServerLogRequest {
	return apiClient.GetServerLog(ctx, model.ProjectId, model.Region, model.ServerId)
}

func outputResult(p *print.Printer, serverLabel string, logLines []string) error {
	p.Outputf("Log for server %q\n", serverLabel)
	for _, line := range logLines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}
		p.Outputln(line)
	}

	return nil
}
