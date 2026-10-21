package describe

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	registry "github.com/stackitcloud/stackit-sdk-go/services/registry/v1api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/registry/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
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
		Use:   fmt.Sprintf("describe %s", artifactoryIdArg),
		Short: "Shows details of a Container Registry artifactory",
		Long:  "Shows details of a Container Registry artifactory.",
		Args:  args.SingleArg(artifactoryIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of an artifactory with ID "xxx"`,
				"$ stackit registry artifactory describe xxx"),
			examples.NewExample(
				`Get details of an artifactory with ID "xxx" in JSON format`,
				"$ stackit registry artifactory describe xxx --output-format json"),
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

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read artifactory: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *registry.APIClient) registry.ApiGetArtifactoryRequest {
	return apiClient.DefaultAPI.GetArtifactory(ctx, model.ProjectId, model.Region, model.ArtifactoryId)
}

func outputResult(p *print.Printer, outputFormat string, artifactory *registry.Artifactory) error {
	return p.OutputResult(outputFormat, artifactory, func() error {
		if artifactory == nil {
			return fmt.Errorf("no artifactory passed")
		}

		table := tables.NewTable()
		table.AddRow("ID", artifactory.Id)
		table.AddSeparator()
		table.AddRow("NAME", artifactory.Name)
		table.AddSeparator()
		table.AddRow("STATE", string(artifactory.State))
		table.AddSeparator()
		table.AddRow("URL", artifactory.Url)
		table.AddSeparator()
		table.AddRow("CONSUMED QUOTA BYTES", artifactory.ConsumedQuotaBytes)
		table.AddSeparator()
		table.AddRow("CREATED AT", artifactory.CreateTime.String())
		if artifactory.Proxy != nil {
			table.AddSeparator()
			table.AddRow("PROXY URL", artifactory.Proxy.Url)
			if artifactory.Proxy.Description != nil && *artifactory.Proxy.Description != "" {
				table.AddSeparator()
				table.AddRow("PROXY DESCRIPTION", *artifactory.Proxy.Description)
			}
			table.AddSeparator()
			table.AddRow("PROXY REGISTRY TYPE", artifactory.Proxy.RegistryType)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
