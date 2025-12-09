package status

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

type statusOutput struct {
	Authenticated bool   `json:"authenticated"`
	Email         string `json:"email,omitempty"`
	AuthFlow      string `json:"auth_flow,omitempty"`
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Shows authentication status for the STACKIT Terraform Provider and SDK",
		Long:  "Shows authentication status for the STACKIT Terraform Provider and SDK, including whether you are authenticated and with which account.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Show authentication status for the STACKIT Terraform Provider and SDK`,
				"$ stackit auth api status"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Check if access token exists (primary credential check)
			accessToken, err := auth.GetAuthFieldWithContext(auth.StorageContextAPI, auth.ACCESS_TOKEN)
			if err != nil || accessToken == "" {
				// Not authenticated
				return outputStatus(params.Printer, model, statusOutput{
					Authenticated: false,
				})
			}

			// Get optional fields for display
			flow, _ := auth.GetAuthFlowWithContext(auth.StorageContextAPI)
			email, err := auth.GetAuthFieldWithContext(auth.StorageContextAPI, auth.USER_EMAIL)
			if err != nil {
				email = ""
			}

			return outputStatus(params.Printer, model, statusOutput{
				Authenticated: true,
				Email:         email,
				AuthFlow:      string(flow),
			})
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

func outputStatus(p *print.Printer, model *inputModel, status statusOutput) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal status: %w", err)
		}
		p.Outputln(string(details))
		return nil
	default:
		if status.Authenticated {
			p.Outputln("API Authentication Status: Authenticated")
			if status.Email != "" {
				p.Outputf("Email: %s\n", status.Email)
			}
			p.Outputf("Auth Flow: %s\n", status.AuthFlow)
		} else {
			p.Outputln("API Authentication Status: Not authenticated")
			p.Outputln("\nTo authenticate, run: stackit auth api login")
		}
		return nil
	}
}
