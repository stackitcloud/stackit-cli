package update

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaasalpha/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	vpcIdArg = "VPC_ID"

	nameFlag        = "name"
	descriptionFlag = "description"
	labelsFlag      = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	VpcId       string
	Name        *string
	Description *string
	Labels      map[string]any
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", vpcIdArg),
		Short: "Updates a VPC",
		Long:  "Updates a VPC.",
		Args:  args.SingleArg(vpcIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete a vpc with ID "xxx"`,
				"$ stackit alpha vpc delete xxx",
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

			vpcLabel, err := iaasUtils.GetVPCName(ctx, apiClient.DefaultAPI, model.ProjectId, model.VpcId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get vpc name: %v", err)
				vpcLabel = model.VpcId
			}

			prompt := fmt.Sprintf("Are you sure you want to update vpc %q?", vpcLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update vpc: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, vpcLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "Name of the VPC")
	cmd.Flags().String(descriptionFlag, "", "Description of the VPC")
	cmd.Flags().StringToString(labelsFlag, nil, "Comma separated list of labels of the VPC")

	cmd.MarkFlagsOneRequired(nameFlag, descriptionFlag, labelsFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	vpcId := inputArgs[0]
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		VpcId:           vpcId,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		Description:     flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:          flags.FlagToStringToAny(p, cmd, labelsFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiPartialUpdateVPCRequest {
	payload := iaas.PartialUpdateVPCPayload{
		Name:        model.Name,
		Description: model.Description,
		Labels:      model.Labels,
	}
	return apiClient.DefaultAPI.PartialUpdateVPC(ctx, model.ProjectId, model.VpcId).PartialUpdateVPCPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, vpcLabel string, resp *iaas.VPC) error {
	if resp == nil {
		return fmt.Errorf("update vpc response is nil")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Updated vpc %q.\n", vpcLabel)
		return nil
	})
}
