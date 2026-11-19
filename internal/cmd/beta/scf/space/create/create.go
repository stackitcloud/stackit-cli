package create

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
)

const (
	organizationIdFlag = "organization-id"
	nameFlag           = "name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	Name           string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a STACKIT Cloud Foundry space",
		Long:  "Creates a STACKIT Cloud Foundry space in an organization.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a STACKIT Cloud Foundry space with name "my-space" in the organization with ID "xxx"`,
				`$ stackit beta scf space create --organization-id xxx --name my-space`),
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

			organizationLabel := scfUtils.GetOrganizationLabel(ctx, p.Printer, apiClient.DefaultAPI, model.ProjectId, model.Region, model.OrganizationId)

			prompt := fmt.Sprintf("Are you sure you want to create STACKIT Cloud Foundry space %q in organization %q?", model.Name, organizationLabel)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create STACKIT Cloud Foundry space: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, organizationLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "The UUID of the organization")
	cmd.Flags().String(nameFlag, "", "Name of the space")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiCreateSpaceRequest {
	req := apiClient.CreateSpace(ctx, model.ProjectId, model.Region, model.OrganizationId)
	return req.CreateSpacePayload(scf.CreateSpacePayload{Name: model.Name})
}

func outputResult(p *print.Printer, outputFormat, organizationLabel string, resp *scf.Space) error {
	return p.OutputResult(outputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Created STACKIT Cloud Foundry space in organization %q, but no space ID was returned.\n", organizationLabel)
			return nil
		}

		p.Outputf("Created STACKIT Cloud Foundry space in organization %q. Space ID: %s\n", organizationLabel, resp.Guid)
		return nil
	})
}
