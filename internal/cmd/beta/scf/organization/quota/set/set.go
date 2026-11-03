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

	quotaIdFlag = "quota-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	OrganizationId string
	QuotaId        *string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("set %s", orgIdArg),
		Short: "Applies a quota to a STACKIT Cloud Foundry organization",
		Long:  "Applies an existing quota to a STACKIT Cloud Foundry organization. The quota limits are not modified by this command.",
		Args:  args.SingleArg(orgIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Apply the quota with ID "yyy" to the STACKIT Cloud Foundry organization with ID "xxx"`,
				`$ stackit beta scf organization quota set xxx --quota-id yyy`),
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

			prompt := fmt.Sprintf("Are you sure you want to apply quota %q to STACKIT Cloud Foundry organization %q?", *model.QuotaId, model.OrganizationId)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			if _, err = req.Execute(); err != nil {
				return fmt.Errorf("apply quota to STACKIT Cloud Foundry organization: %w", err)
			}

			return outputResult(p.Printer, model)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), quotaIdFlag, "The UUID of the quota to apply")

	err := flags.MarkFlagsRequired(cmd, quotaIdFlag)
	cobra.CheckErr(err)
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
		QuotaId:         flags.FlagToStringPointer(p, cmd, quotaIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient scf.DefaultAPI) scf.ApiApplyOrganizationQuotaRequest {
	payload := scf.ApplyOrganizationQuotaPayload{
		QuotaId: *model.QuotaId,
	}

	req := apiClient.ApplyOrganizationQuota(ctx, model.ProjectId, model.Region, model.OrganizationId)
	return req.ApplyOrganizationQuotaPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel) error {
	return p.OutputResult(model.OutputFormat, nil, func() error {
		p.Outputf("Applied quota %s to STACKIT Cloud Foundry organization %s\n", *model.QuotaId, model.OrganizationId)
		return nil
	})
}
