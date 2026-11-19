package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	orgIdArg = "ORGANIZATION_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", orgIdArg),
		Short: "Deletes the organization manager of a STACKIT Cloud Foundry organization",
		Long:  "Deletes the organization manager of a STACKIT Cloud Foundry organization.",
		Args:  args.SingleArg(orgIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete the organization manager of a STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization manager delete xxx`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			prompt := fmt.Sprintf("Are you sure you want to delete the organization manager of STACKIT Cloud Foundry organization %q?", model.OrganizationId)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			if _, err = req.Execute(); err != nil {
				return fmt.Errorf("delete STACKIT Cloud Foundry organization manager: %w", err)
			}

			p.Printer.Info("Deleted organization manager of STACKIT Cloud Foundry organization %s\n", model.OrganizationId)

			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	organizationId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  organizationId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiDeleteOrgManagerRequest {
	return apiClient.DeleteOrgManager(ctx, model.ProjectId, model.Region, model.OrganizationId)
}
