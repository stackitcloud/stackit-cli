package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2alpha1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	nameFlag        = "name"
	descriptionFlag = "description"
	labelsFlag      = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name        string
	Description *string
	Labels      map[string]any
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a VPC",
		Long: fmt.Sprintf("%s\n%s",
			"Creates a VPC.",
			"VPC can only be used in projects with enabled VPC."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a vpc with name "my-vpc"`,
				"$ stackit alpha vpc create --name my-vpc",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a vpc for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create vpc: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the VPC")
	cmd.Flags().String(descriptionFlag, "", "Description of the VPC")
	cmd.Flags().StringToString(labelsFlag, nil, "Comma separated list of labels of the VPC")

	err := cmd.MarkFlagRequired(nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringValue(p, cmd, nameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          flags.FlagToStringToAny(p, cmd, labelsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateVPCRequest {
	payload := iaas.CreateVPCPayload{
		Name:        model.Name,
		Description: model.Description,
		Labels:      model.Labels,
	}
	return apiClient.DefaultAPI.CreateVPC(ctx, model.ProjectId).CreateVPCPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *iaas.VPC) error {
	if resp == nil {
		return fmt.Errorf("create vpc response is nil")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Created vpc for %q. VPC ID: %s\n", projectLabel, resp.Id)
		return nil
	})
}
