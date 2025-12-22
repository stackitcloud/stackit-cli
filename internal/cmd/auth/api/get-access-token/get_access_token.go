package getaccesstoken

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(p *types.CmdParams) *cobra.Command {
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
			model, err := parseInput(p.Printer, cmd, args)
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

			accessToken, err := auth.GetValidAccessTokenWithContext(p.Printer, auth.StorageContextAPI)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get valid access token: %v", err)
				return &cliErr.SessionExpiredError{}
			}

			result := map[string]string{
				"access_token": accessToken,
			}
			return p.Printer.OutputResult(model.OutputFormat, result, func() error {
				p.Printer.Outputln(accessToken)
				return nil
			})
		},
	}

	// hide project id flag from help command because it could mislead users
	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {
		cobra.CheckErr(command.Flags().MarkHidden(globalflags.ProjectIdFlag))
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
