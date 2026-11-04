package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	scfUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	orgIdArg = "ORGANIZATION_ID"

	nameFlag      = "name"
	suspendedFlag = "suspended"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	Name           *string
	Suspended      *bool
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", orgIdArg),
		Short: "Updates a STACKIT Cloud Foundry organization",
		Long:  "Updates a STACKIT Cloud Foundry organization. Only the specified fields are updated.",
		Args:  args.SingleArg(orgIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the name of a STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization update xxx --name new-name`),
			examples.NewExample(
				`Suspend a STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization update xxx --suspended`),
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

			organizationLabel, err := scfUtils.GetOrganizationName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.OrganizationId)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get organization name: %v", err)
				organizationLabel = model.OrganizationId
			} else if organizationLabel == "" {
				organizationLabel = model.OrganizationId
			}

			prompt := fmt.Sprintf("Are you sure you want to update STACKIT Cloud Foundry organization %q?", organizationLabel)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}
			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update STACKIT Cloud Foundry organization: %w", err)
			}

			return outputResult(p.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the organization")
	cmd.Flags().Bool(suspendedFlag, false, "Whether the organization is suspended")
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	organizationId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := &inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  organizationId,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		Suspended:       flags.FlagToBoolPointer(p, cmd, suspendedFlag),
	}

	if model.Name == nil && model.Suspended == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiUpdateOrganizationRequest {
	payload := scf.UpdateOrganizationPayload{
		Name:      model.Name,
		Suspended: model.Suspended,
	}

	req := apiClient.UpdateOrganization(ctx, model.ProjectId, model.Region, model.OrganizationId)
	return req.UpdateOrganizationPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, resp *scf.Organization) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Updated STACKIT Cloud Foundry organization %q\n", model.OrganizationId)
		return nil
	})
}
