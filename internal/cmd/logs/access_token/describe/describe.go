package describe

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/logs/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/logs"

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
		Short: "Shows details of a Logs access token",
		Long:  "Shows details of a Logs access token.",
		Args:  args.SingleArg(accessTokenIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Show details of a Logs access token with ID "xxx"`,
				"$ stackit logs access-token describe xxx",
			),
			examples.NewExample(
				`Show details of a Logs access token with ID "xxx" in JSON format`,
				"$ stackit logs access-token describe xxx --output-format json",
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
				return fmt.Errorf("read access token: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "ID of the Logs instance")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	accessTokenId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		AccessTokenId:   accessTokenId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *logs.APIClient) logs.ApiGetAccessTokenRequest {
	return apiClient.GetAccessToken(ctx, model.ProjectId, model.Region, model.InstanceId, model.AccessTokenId)
}

func outputResult(p *print.Printer, outputFormat string, token *logs.AccessToken) error {
	if token == nil {
		return fmt.Errorf("access token cannot be nil")
	}
	return p.OutputResult(outputFormat, token, func() error {
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(token.Id))
		table.AddSeparator()
		table.AddRow("DISPLAY NAME", utils.PtrString(token.DisplayName))
		table.AddSeparator()
		table.AddRow("DESCRIPTION", utils.PtrString(token.Description))
		table.AddSeparator()
		table.AddRow("PERMISSIONS", utils.PtrString(token.Permissions))
		table.AddSeparator()
		table.AddRow("CREATOR", utils.PtrString(token.Creator))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(token.Status))
		table.AddSeparator()
		table.AddRow("EXPIRES", utils.PtrString(token.Expires))
		table.AddSeparator()
		table.AddRow("VALID UNTIL", utils.PtrString(token.ValidUntil))
		table.AddSeparator()

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
		return nil
	})
}
