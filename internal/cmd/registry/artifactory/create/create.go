package create

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/registry/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	nameFlag = "name"

	proxyUrlFlag          = "proxy-url"
	proxyRegistryTypeFlag = "proxy-registry-type"
	proxyDescriptionFlag  = "proxy-description"
	proxyAccessKeyFlag    = "proxy-access-key"
	proxyAccessSecretFlag = "proxy-access-secret"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name              string
	ProxyUrl          *string
	ProxyRegistryType *string
	ProxyDescription  *string
	ProxyAccessKey    *string
	ProxyAccessSecret *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Container Registry artifactory",
		Long:  "Creates a Container Registry artifactory.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an artifactory with name "my-registry"`,
				"$ stackit registry artifactory create --name my-registry",
			),
			examples.NewExample(
				`Create a proxy-cache artifactory for Docker Hub`,
				"$ stackit registry artifactory create --name docker-proxy --proxy-url https://registry-1.docker.io --proxy-registry-type docker-hub",
			),
			examples.NewExample(
				`Create a proxy-cache artifactory with credentials`,
				"$ stackit registry artifactory create --name my-proxy --proxy-url https://custom-registry.com --proxy-registry-type harbor --proxy-access-key my-user --proxy-access-secret my-secret",
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

			prompt := fmt.Sprintf("Are you sure you want to create an artifactory for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create artifactory: %w", err)
			}
			artifactoryId := resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating artifactory", func() error {
					_, err = wait.CreateArtifactoryWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, artifactoryId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for artifactory creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, projectLabel, artifactoryId, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Artifactory name")
	cmd.Flags().String(proxyUrlFlag, "", "Upstream registry URL to proxy")
	cmd.Flags().String(proxyRegistryTypeFlag, "", "Registry type of the upstream proxy (e.g., docker-hub, harbor)")
	cmd.Flags().String(proxyDescriptionFlag, "", "Description of the proxy registry")
	cmd.Flags().String(proxyAccessKeyFlag, "", "Access key for upstream proxy registry authentication")
	cmd.Flags().String(proxyAccessSecretFlag, "", "Access secret for upstream proxy registry authentication")

	err := flags.MarkFlagsRequired(cmd, nameFlag)
	cobra.CheckErr(err)

	cmd.MarkFlagsRequiredTogether(proxyUrlFlag, proxyRegistryTypeFlag)
	cmd.MarkFlagsRequiredTogether(proxyAccessKeyFlag, proxyAccessSecretFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	proxyUrl := flags.FlagToStringPointer(p, cmd, proxyUrlFlag)
	proxyRegistryType := flags.FlagToStringPointer(p, cmd, proxyRegistryTypeFlag)
	proxyDesc := flags.FlagToStringPointer(p, cmd, proxyDescriptionFlag)
	proxyKey := flags.FlagToStringPointer(p, cmd, proxyAccessKeyFlag)
	proxySecret := flags.FlagToStringPointer(p, cmd, proxyAccessSecretFlag)

	hasProxyDetails := proxyDesc != nil || proxyKey != nil || proxySecret != nil
	if hasProxyDetails && (proxyUrl == nil || proxyRegistryType == nil) {
		return nil, fmt.Errorf("both --%s and --%s are required when configuring a proxy", proxyUrlFlag, proxyRegistryTypeFlag)
	}

	model := inputModel{
		GlobalFlagModel:   globalFlags,
		Name:              flags.FlagToStringValue(p, cmd, nameFlag),
		ProxyUrl:          proxyUrl,
		ProxyRegistryType: proxyRegistryType,
		ProxyDescription:  proxyDesc,
		ProxyAccessKey:    proxyKey,
		ProxyAccessSecret: proxySecret,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *registry.APIClient) registry.ApiCreateArtifactoryRequest {
	req := apiClient.DefaultAPI.CreateArtifactory(ctx, model.ProjectId, model.Region)
	payload := registry.CreateArtifactoryPayload{
		Name: model.Name,
	}

	if model.ProxyUrl != nil && model.ProxyRegistryType != nil {
		proxyReq := registry.NewCreateProxyRequest(*model.ProxyRegistryType, *model.ProxyUrl)
		if model.ProxyDescription != nil {
			proxyReq.Description = model.ProxyDescription
		}
		if model.ProxyAccessKey != nil && model.ProxyAccessSecret != nil {
			proxyReq.Credentials = &registry.Credentials{
				AccessKey:    *model.ProxyAccessKey,
				AccessSecret: *model.ProxyAccessSecret,
			}
		}
		payload.Proxy = proxyReq
	}

	return req.CreateArtifactoryPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel, artifactoryId string, resp *registry.Artifactory) error {
	if model == nil {
		return fmt.Errorf("no model passed")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			return fmt.Errorf("no response defined")
		}

		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s artifactory for project %q. Artifactory ID: %s\n", operationState, projectLabel, artifactoryId)
		return nil
	})
}
