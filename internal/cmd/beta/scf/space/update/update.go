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
	spaceIdArg = "SPACE_ID"

	organizationIdFlag = "organization-id"
	nameFlag           = "name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	SpaceId        string
	Name           string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", spaceIdArg),
		Short: "Updates a STACKIT Cloud Foundry space",
		Long:  "Updates a STACKIT Cloud Foundry space.",
		Args:  args.SingleArg(spaceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Rename the STACKIT Cloud Foundry space with ID "yyy" in the organization with ID "xxx"`,
				`$ stackit beta scf space update yyy --organization-id xxx --name new-name`),
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

			spaceLabel := scfUtils.GetSpaceLabel(ctx, p.Printer, apiClient.DefaultAPI, model.ProjectId, model.Region, model.OrganizationId, model.SpaceId)

			prompt := fmt.Sprintf("Are you sure you want to update STACKIT Cloud Foundry space %q?", spaceLabel)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update STACKIT Cloud Foundry space: %w", err)
			}

			return outputResult(p.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "The UUID of the organization")
	cmd.Flags().String(nameFlag, "", "New name of the space")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	spaceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringValue(p, cmd, organizationIdFlag),
		SpaceId:         spaceId,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiUpdateSpaceRequest {
	req := apiClient.UpdateSpace(ctx, model.ProjectId, model.Region, model.OrganizationId, model.SpaceId)
	return req.UpdateSpacePayload(scf.UpdateSpacePayload{Name: model.Name})
}

func outputResult(p *print.Printer, model *inputModel, resp *scf.Space) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		spaceLabel := model.SpaceId
		if resp != nil && resp.Name != "" {
			spaceLabel = resp.Name
		}
		p.Outputf("Updated STACKIT Cloud Foundry space %q\n", spaceLabel)
		return nil
	})
}
