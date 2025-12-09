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
		Short: "Prints a short-lived access token for the STACKIT Terraform Provider and SDK",
		Long:  "Prints a short-lived access token for the STACKIT Terraform Provider and SDK which can be used e.g. for API calls.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Print a short-lived access token for the STACKIT Terraform Provider and SDK`,
				"$ stackit auth api get-access-token"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			userSessionExpired, err := auth.UserSessionExpiredWithContext(auth.StorageContextAPI)
			if err != nil {
				return err
			}
			if userSessionExpired {
				return &cliErr.SessionExpiredError{}
			}

			accessToken, err := auth.GetValidAccessTokenWithContext(params.Printer, auth.StorageContextAPI)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get valid access token: %v", err)
				return &cliErr.SessionExpiredError{}
			}

			switch model.OutputFormat {
			case print.JSONOutputFormat:
				details, err := json.MarshalIndent(map[string]string{
					"access_token": accessToken,
				}, "", "  ")
				if err != nil {
					return fmt.Errorf("marshal access token: %w", err)
				}
				params.Printer.Outputln(string(details))

				return nil
			default:
				params.Printer.Outputln(accessToken)

				return nil
			}
		},
	}

	// hide project id flag from help command because it could mislead users
	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		_ = command.Flags().MarkHidden(globalflags.ProjectIdFlag) // nolint:errcheck // there's no chance to handle the error here
		command.Parent().HelpFunc()(command, strings)
	})

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
	}

	p.DebugInputModel(model)
	return &model, nil
}
