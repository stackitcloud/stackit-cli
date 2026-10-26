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
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	orgIdArg = "ORGANIZATION_ID"

	userGuidFlag = "user-guid"
	userNameFlag = "user-name"
)

var roleFlag = flags.StringEnumFlag(
	"role",
	scf.AllowedOrgRoleTypeEnumValues,
	"The role to assign to the user",
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	Role           scf.OrgRoleType
	UserGuid       *string
	UserName       *string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("set %s", orgIdArg),
		Short: "Assigns a role to a user of a STACKIT Cloud Foundry organization",
		Long:  "Assigns a role to a user of a STACKIT Cloud Foundry organization. The user is identified either by its GUID or by its username.",
		Args:  args.SingleArg(orgIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Assign the role "organization_manager" to the user with GUID "yyy" in the STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization role set xxx --role organization_manager --user-guid yyy`),
			examples.NewExample(
				`Assign the role "organization_auditor" to the user with username "jane.doe" in the STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization role set xxx --role organization_auditor --user-name jane.doe`),
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

			prompt := fmt.Sprintf("Are you sure you want to assign role %q to user %q in STACKIT Cloud Foundry organization %q?", model.Role, userLabel(model), model.OrganizationId)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("set STACKIT Cloud Foundry organization role: %w", err)
			}

			return outputResult(p.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	roleFlag.Register(cmd.Flags())
	cmd.Flags().Var(flags.UUIDFlag(), userGuidFlag, "The GUID of the user to assign the role to")
	cmd.Flags().String(userNameFlag, "", "The username of the user to assign the role to")

	err := flags.MarkFlagsRequired(cmd, roleFlag.Name())
	cobra.CheckErr(err)
	cmd.MarkFlagsOneRequired(userGuidFlag, userNameFlag)
	cmd.MarkFlagsMutuallyExclusive(userGuidFlag, userNameFlag)
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
		Role:            roleFlag.Get(),
		UserGuid:        flags.FlagToStringPointer(p, cmd, userGuidFlag),
		UserName:        flags.FlagToStringPointer(p, cmd, userNameFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiCreateOrgRoleRequest {
	payload := scf.CreateOrgRolePayload{
		Type:     model.Role,
		UserGuid: model.UserGuid,
		UserName: model.UserName,
	}

	req := apiClient.CreateOrgRole(ctx, model.ProjectId, model.Region, model.OrganizationId)
	return req.CreateOrgRolePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, resp *scf.OrgRoleResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		p.Outputf("Assigned role %q to user %q in STACKIT Cloud Foundry organization %s\n", model.Role, userLabel(model), model.OrganizationId)
		return nil
	})
}

// userLabel returns the user identifier given on the command line, either the GUID or the username
func userLabel(model *inputModel) string {
	if model.UserGuid != nil {
		return *model.UserGuid
	}
	return utils.PtrString(model.UserName)
}
