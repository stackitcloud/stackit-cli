package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"
	"github.com/stackitcloud/stackit-sdk-go/services/registry/v1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
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

	proxyUrlFlag          = "proxy-url"
	proxyDescriptionFlag  = "proxy-description"
	proxyAccessKeyFlag    = "proxy-access-key"
	proxyAccessSecretFlag = "proxy-access-secret"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ArtifactoryId     string
	ProxyUrl          *string
	ProxyDescription  *string
	ProxyAccessKey    *string
	ProxyAccessSecret *string `json:"-"`
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", artifactoryIdArg),
		Short: "Updates a Container Registry artifactory",
		Long:  "Updates a Container Registry artifactory.",
		Args:  args.SingleArg(artifactoryIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the proxy configuration for an artifactory with ID "xxx"`,
				"$ stackit registry artifactory update xxx --proxy-url https://registry-1.docker.io"),
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

			prompt := fmt.Sprintf("Are you sure you want to update artifactory %q?", artifactoryLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update artifactory: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Updating artifactory", func() error {
					waitResp, waitErr := wait.UpdateArtifactoryWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.ArtifactoryId).WaitWithContext(ctx)
					if waitErr != nil {
						return waitErr
					}
					resp = waitResp
					return nil
				})
				if err != nil {
					return fmt.Errorf("wait for artifactory update: %w", err)
				}
			}

			return outputResult(params.Printer, model, artifactoryLabel, resp)
		},
	}
	configureFlags(cmd, params)
	return cmd
}

func configureFlags(cmd *cobra.Command, params *types.CmdParams) {
	cmd.Flags().String(proxyUrlFlag, "", "Upstream registry URL to proxy")
	cmd.Flags().String(proxyDescriptionFlag, "", "Description of the proxy registry")
	cmd.Flags().String(proxyAccessKeyFlag, "", "Access key for upstream proxy registry authentication")

	secret := flags.SecretFlag(proxyAccessSecretFlag, params)
	cmd.Flags().Var(secret, proxyAccessSecretFlag, secret.Usage())

	cmd.MarkFlagsRequiredTogether(proxyAccessKeyFlag, proxyAccessSecretFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	artifactoryId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	proxyUrl := flags.FlagToStringPointer(p, cmd, proxyUrlFlag)
	proxyDesc := flags.FlagToStringPointer(p, cmd, proxyDescriptionFlag)
	proxyKey := flags.FlagToStringPointer(p, cmd, proxyAccessKeyFlag)
	proxySecret := flags.SecretFlagToStringPointer(p, cmd, proxyAccessSecretFlag)

	if proxyUrl == nil && proxyDesc == nil && proxyKey == nil && proxySecret == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		ArtifactoryId:     artifactoryId,
		ProxyUrl:          proxyUrl,
		ProxyDescription:  proxyDesc,
		ProxyAccessKey:    proxyKey,
		ProxyAccessSecret: proxySecret,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *registry.APIClient) registry.ApiPatchArtifactoryRequest {
	req := apiClient.DefaultAPI.PatchArtifactory(ctx, model.ProjectId, model.Region, model.ArtifactoryId)

	proxyReq := &registry.UpdateProxyRequest{
		Url:         model.ProxyUrl,
		Description: model.ProxyDescription,
	}

	if model.ProxyAccessKey != nil && model.ProxyAccessSecret != nil {
		proxyReq.Credentials = &registry.Credentials{
			AccessKey:    *model.ProxyAccessKey,
			AccessSecret: *model.ProxyAccessSecret,
		}
	}

	return req.PatchArtifactoryPayload(registry.PatchArtifactoryPayload{
		Proxy: proxyReq,
	})
}

func outputResult(p *print.Printer, model *inputModel, artifactoryLabel string, resp *registry.Artifactory) error {
	if model == nil {
		return fmt.Errorf("no model passed")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s artifactory %q\n", operationState, artifactoryLabel)
		return nil
	})
}
