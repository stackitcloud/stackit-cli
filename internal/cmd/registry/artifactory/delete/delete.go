package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/registry/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/registry/client"
	registryUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/registry/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	artifactoryIdArg = "ARTIFACTORY_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ArtifactoryId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("delete %s", artifactoryIdArg),
		Short: "Deletes a Container Registry artifactory",
		Long:  "Deletes a Container Registry artifactory.",
		Args:  args.SingleArg(artifactoryIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Delete an artifactory with ID "xxx"`,
				"$ stackit registry artifactory delete xxx"),
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

			artifactoryLabel, err := registryUtils.GetArtifactoryName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.ArtifactoryId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get artifactory name: %v", err)
				artifactoryLabel = model.ArtifactoryId
			}

			prompt := fmt.Sprintf("Are you sure you want to delete artifactory %q? (This cannot be undone)", artifactoryLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			err = req.Execute()
			if err != nil {
				return fmt.Errorf("delete artifactory: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Deleting artifactory", func() error {
					_, err = wait.DeleteArtifactoryWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.ArtifactoryId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for artifactory deletion: %w", err)
				}
			}

			operationState := "Deleted"
			if model.Async {
				operationState = "Triggered deletion of"
			}
			params.Printer.Outputf("%s artifactory %q\n", operationState, artifactoryLabel)
			return nil
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	artifactoryId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ArtifactoryId:   artifactoryId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *registry.APIClient) registry.ApiDeleteArtifactoryRequest {
	return apiClient.DefaultAPI.DeleteArtifactory(ctx, model.ProjectId, model.Region, model.ArtifactoryId)
}
