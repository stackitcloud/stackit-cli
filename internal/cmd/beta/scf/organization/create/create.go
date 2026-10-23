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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/scf/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	nameFlag       = "name"
	platformIdFlag = "platform-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name       *string
	PlatformId *string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a STACKIT Cloud Foundry organization",
		Long:  "Creates a STACKIT Cloud Foundry organization.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a STACKIT Cloud Foundry organization with name "my-org"`,
				`$ stackit beta scf organization create --name my-org`),
			examples.NewExample(
				`Create a STACKIT Cloud Foundry organization with name "my-org" on platform "xxx"`,
				`$ stackit beta scf organization create --name my-org --platform-id xxx`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a STACKIT Cloud Foundry organization for project %q?", projectLabel)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create STACKIT Cloud Foundry organization: %w", err)
			}

			return outputResult(p.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the organization")
	cmd.Flags().Var(flags.UUIDFlag(), platformIdFlag, "The UUID of the platform to create the organization on")

	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		PlatformId:      flags.FlagToStringPointer(p, cmd, platformIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *scf.APIClient) scf.ApiCreateOrganizationRequest {
	payload := scf.CreateOrganizationPayload{
		Name:       *model.Name,
		PlatformId: model.PlatformId,
	}

	req := apiClient.DefaultAPI.CreateOrganization(ctx, model.ProjectId, model.Region)
	return req.CreateOrganizationPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *scf.OrganizationCreateResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Created STACKIT Cloud Foundry organization for project %q, but no organization ID was returned.\n", projectLabel)
			return nil
		}

		p.Outputf("Created STACKIT Cloud Foundry organization for project %q. Organization ID: %s\n", projectLabel, resp.Guid)
		return nil
	})
}
