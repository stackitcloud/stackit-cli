package set

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
	userGuidFlag       = "user-guid"
	userNameFlag       = "user-name"
)

var roleFlag = flags.StringEnumFlag(
	"role",
	scf.AllowedSpaceRoleTypeEnumValues,
	"The role to assign to the user",
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	SpaceId        string
	Role           scf.SpaceRoleType
	UserGuid       *string
	UserName       *string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("set %s", spaceIdArg),
		Short: "Assigns a role to a user of a STACKIT Cloud Foundry space",
		Long:  "Assigns a role to a user of a STACKIT Cloud Foundry space. The user is identified either by its GUID or by its username.",
		Args:  args.SingleArg(spaceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Assign the role "space_developer" to the user with GUID "zzz" in the STACKIT Cloud Foundry space with ID "yyy" of the organization with ID "xxx"`,
				`$ stackit beta scf space role set yyy --organization-id xxx --role space_developer --user-guid zzz`),
			examples.NewExample(
				`Assign the role "space_auditor" to the user with username "jane.doe" in the STACKIT Cloud Foundry space with ID "yyy" of the organization with ID "xxx"`,
				`$ stackit beta scf space role set yyy --organization-id xxx --role space_auditor --user-name jane.doe`),
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

			prompt := fmt.Sprintf("Are you sure you want to assign role %q to user %q in STACKIT Cloud Foundry space %q?", model.Role, scfUtils.UserLabel(model.UserGuid, model.UserName), spaceLabel)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("set STACKIT Cloud Foundry space role: %w", err)
			}

			return outputResult(p.Printer, model, spaceLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	roleFlag.Register(cmd.Flags())
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "The UUID of the organization")
	cmd.Flags().Var(flags.UUIDFlag(), userGuidFlag, "The GUID of the user to assign the role to")
	cmd.Flags().String(userNameFlag, "", "The username of the user to assign the role to")

	err := flags.MarkFlagsRequired(cmd, roleFlag.Name(), organizationIdFlag)
	cobra.CheckErr(err)
	cmd.MarkFlagsOneRequired(userGuidFlag, userNameFlag)
	cmd.MarkFlagsMutuallyExclusive(userGuidFlag, userNameFlag)
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
		Role:            roleFlag.Get(),
		UserGuid:        flags.FlagToStringPointer(p, cmd, userGuidFlag),
		UserName:        flags.FlagToStringPointer(p, cmd, userNameFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiCreateSpaceRoleRequest {
	payload := scf.CreateSpaceRolePayload{
		Type:     model.Role,
		UserGuid: model.UserGuid,
		UserName: model.UserName,
	}

	req := apiClient.CreateSpaceRole(ctx, model.ProjectId, model.Region, model.OrganizationId, model.SpaceId)
	return req.CreateSpaceRolePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, spaceLabel string, resp *scf.SpaceRoleCreateResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Assigned role %q to user %q in STACKIT Cloud Foundry space %q\n", model.Role, scfUtils.UserLabel(model.UserGuid, model.UserName), spaceLabel)
		return nil
	})
}
