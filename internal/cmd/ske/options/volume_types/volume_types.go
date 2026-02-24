package volume_types

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume-types",
		Short: "Lists SKE provider options for volume-types",
		Long: fmt.Sprintf("%s\n%s",
			"Lists STACKIT Kubernetes Engine (SKE) provider options for volume-types.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List SKE options for volume-types`,
				"$ stackit ske options volume-types"),
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
			req := buildRequest(ctx, apiClient, model)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get SKE provider options: %w", err)
			}

			return outputResult(params.Printer, model, resp)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, apiClient *ske.APIClient, model *inputModel) ske.ApiListProviderOptionsRequest {
	req := apiClient.ListProviderOptions(ctx, model.Region)
	return req
}

func outputResult(p *print.Printer, model *inputModel, options *ske.ProviderOptions) error {
	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("model is nil")
	} else if options == nil {
		return fmt.Errorf("options is nil")
	}

	return p.OutputResult(model.OutputFormat, options, func() error {
		volumeTypes := *options.VolumeTypes

		table := tables.NewTable()
		table.SetHeader("TYPE")
		for i := range volumeTypes {
			z := volumeTypes[i]
			table.AddRow(utils.PtrString(z.Name))
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display output: %w", err)
		}
		return nil
	})
}
