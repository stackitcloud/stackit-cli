package delete

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ImageId string
}

const imageIdArg = "IMAGE_ID"

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", imageIdArg),
		Short: "Deletes an image",
		Long:  "Deletes an image by its internal ID.",
		Args:  args.SingleArg(imageIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Delete an image with ID "xxx"`, `$ stackit image delete xxx`),
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

			imageName, err := iaasUtils.GetImageName(ctx, apiClient, model.ProjectId, model.Region, model.ImageId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get image name: %v", err)
				imageName = model.ImageId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete the image %q for %q?", imageName, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			if err := request.Execute(); err != nil {
				return fmt.Errorf("delete image: %w", err)
			}
			params.Printer.Info("Deleted image %q for %q\n", imageName, projectLabel)

			return nil
		},
	}

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ImageId:         cliArgs[0],
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiDeleteImageRequest {
	request := apiClient.DeleteImage(ctx, model.ProjectId, model.Region, model.ImageId)
	return request
}
