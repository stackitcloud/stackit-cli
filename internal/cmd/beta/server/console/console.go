package console

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"

	"github.com/spf13/cobra"
)

const (
	serverIdArg = "SERVER_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ServerId string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("console %s", serverIdArg),
		Short: "Gets a URL for server remote console",
		Long:  "Gets a URL for server remote console.",
		Args:  args.SingleArg(serverIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get a URL for the server remote console with server ID "xxx"`,
				"$ stackit beta server console xxx",
			),
			examples.NewExample(
				`Get a URL for the server remote console with server ID "xxx" in JSON format`,
				"$ stackit beta server console xxx --output-format json",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			serverLabel, err := iaasUtils.GetServerName(ctx, apiClient, model.ProjectId, model.ServerId)
			if err != nil {
				p.Debug(print.ErrorLevel, "get server name: %v", err)
				serverLabel = model.ServerId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("server console: %w", err)
			}

			return outputResult(p, model, serverLabel, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	serverId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ServerId:        serverId,
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetServerConsoleRequest {
	return apiClient.GetServerConsole(ctx, model.ProjectId, model.ServerId)
}

func outputResult(p *print.Printer, model *inputModel, serverLabel string, serverUrl *iaas.ServerConsoleUrl) error {
	outputFormat := model.OutputFormat

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(serverUrl, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal url: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(serverUrl, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal url: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		if serverUrl.GetUrl() == nil {
			return fmt.Errorf("server url is nil")
		}
		// unescape url in order to get rid of e.g. %40
		unescapedURL, err := url.PathUnescape(*serverUrl.GetUrl())
		if err != nil {
			return fmt.Errorf("unescape url: %w", err)
		}

		p.Outputf("Remote console URL %q for server %q\n", unescapedURL, serverLabel)

		return nil
	}
}
