package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	telemetryrouter "github.com/stackitcloud/stackit-sdk-go/services/telemetryrouter/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/telemetryrouter/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	instanceIdFlag   = "instance-id"
	accessTokenIdArg = "ACCESS_TOKEN_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId    string
	AccessTokenId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", accessTokenIdArg),
		Short: "Shows details of a TelemetryRouter access token",
		Long:  "Shows details of a TelemetryRouter access token.",
		Args:  args.SingleArg(accessTokenIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a TelemetryRouter access token with ID "xxx"`,
				"$ stackit beta telemetryrouter access-token describe xxx --instance-id yyy",
			),
			examples.NewExample(
				`Show details of a TelemetryRouter access token with ID "xxx" in JSON format`,
				"$ stackit beta telemetryrouter access-token describe xxx --instance-id yyy --output-format json",
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read TelemetryRouter access token: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the TelemetryRouter instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	accessTokenId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		AccessTokenId:   accessTokenId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *telemetryrouter.APIClient) telemetryrouter.ApiGetAccessTokenRequest {
	return apiClient.DefaultAPI.GetAccessToken(ctx, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)
}

func outputResult(p *print.Printer, outputFormat string, token *telemetryrouter.GetAccessTokenResponse) error {
	if token == nil {
		return fmt.Errorf("access token cannot be nil")
	}
	return p.OutputResult(outputFormat, token, func() error {
		table := tables.NewTable()
		table.AddRow("ID", token.Id)
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", token.DisplayName)
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(token.Description))
		table.AddSeparator()
		table.AddRow("CREATOR ID", token.CreatorId)
		table.AddSeparator()
		table.AddRow("STATUS", token.Status)
		table.AddSeparator()
		table.AddRow("EXPIRATION TIME", utils.PtrString(token.ExpirationTime.Get()))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
