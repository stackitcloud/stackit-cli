package getaccesstoken

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-access-token",
		Short: "Prints a short-lived access token.",
		Long:  "Prints a short-lived access token which can be used e.g. for API calls.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Print a short-lived access token`,
				"$ stackit auth get-access-token"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			userSessionExpired, err := auth.UserSessionExpired()
			if err != nil {
				return err
			}
			if userSessionExpired {
				return &cliErr.SessionExpiredError{}
			}

			accessToken, err := auth.GetValidAccessToken(params.Printer)
			if err != nil {
				return err
			}

			switch model.OutputFormat {
			case print.JSONOutputFormat:
				details, err := json.MarshalIndent(map[string]string{
					"access_token": accessToken,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("marshal image list: %w", err)
				}
				params.Printer.Outputln(string(details))

				return nil
			case print.YAMLOutputFormat:
				params.Printer.Outputln(accessToken)

				return nil
			default:
				params.Printer.Outputln(accessToken)

				return nil
			}
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
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
