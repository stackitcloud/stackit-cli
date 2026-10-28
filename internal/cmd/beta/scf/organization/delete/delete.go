package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/scf/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
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
		Short: "Deletes a STACKIT Cloud Foundry organization",
		Long:  "Deletes a STACKIT Cloud Foundry organization.",
		Args:  args.SingleArg(orgIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization delete xxx`),
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

			prompt := fmt.Sprintf("Are you sure you want to delete STACKIT Cloud Foundry organization %q?", model.OrganizationId)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			if _, err = req.Execute(); err != nil {
				return fmt.Errorf("delete STACKIT Cloud Foundry organization: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(p.Printer, "Deleting STACKIT Cloud Foundry organization", func() error {
					_, err = wait.DeleteOrganizationWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.OrganizationId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for STACKIT Cloud Foundry organization deletion: %w", err)
				}
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			p.Printer.Info("%s STACKIT Cloud Foundry organization %s\n", operationState, model.OrganizationId)

			return nil
		},
	}
	return cmd
}

// parseInput parses the command arguments and flags into a standardized model
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

// buildRequest creates the API request to delete a STACKIT Cloud Foundry organization
func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiDeleteOrganizationRequest {
	return apiClient.DeleteOrganization(ctx, model.ProjectId, model.Region, model.OrganizationId)
}
