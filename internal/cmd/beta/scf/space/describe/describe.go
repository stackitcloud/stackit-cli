package describe

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	scf "github.com/stackitcloud/stackit-sdk-go/services/scf/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	spaceIdArg = "SPACE_ID"

	organizationIdFlag = "organization-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	SpaceId        string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", spaceIdArg),
		Short: "Shows details of a STACKIT Cloud Foundry space",
		Long:  "Shows details of a STACKIT Cloud Foundry space.",
		Args:  args.SingleArg(spaceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of the STACKIT Cloud Foundry space with ID "yyy" in the organization with ID "xxx"`,
				`$ stackit beta scf space describe yyy --organization-id xxx`),
			examples.NewExample(
				`Get details of the STACKIT Cloud Foundry space with ID "yyy" in the organization with ID "xxx" in JSON format`,
				`$ stackit beta scf space describe yyy --organization-id xxx --output-format json`),
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

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get STACKIT Cloud Foundry space: %w", err)
			}

			return outputResult(p.Printer, model.OutputFormat, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "The UUID of the organization")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
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
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiGetSpaceRequest {
	return apiClient.GetSpace(ctx, model.ProjectId, model.Region, model.OrganizationId, model.SpaceId)
}

func outputResult(p *print.Printer, outputFormat string, space *scf.Space) error {
	return p.OutputResult(outputFormat, space, func() error {
		if space == nil {
			return fmt.Errorf("received nil response, could not display details")
		}

		table := tables.NewTable()
		table.SetHeader("Attribute", "Value")

		table.AddRow("ID", space.GetGuid())
		table.AddRow("Name", space.GetName())
		table.AddRow("Organization ID", space.GetOrgId())
		table.AddRow("Platform ID", space.GetPlatformId())
		table.AddRow("Project ID", space.GetProjectId())
		table.AddRow("Region", space.GetRegion())
		table.AddRow("Created", space.CreatedAt.Format(time.DateTime))
		table.AddRow("Updated", space.UpdatedAt.Format(time.DateTime))

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
