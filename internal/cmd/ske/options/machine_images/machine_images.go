package machine_images

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

type inputModel struct {
	globalflags.GlobalFlagModel
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "machine-images",
		Short: "Lists SKE provider options for machine-images",
		Long:  "Lists STACKIT Kubernetes Engine (SKE) provider options for machine-images.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List SKE options for machine-images`,
				"$ stackit ske options machine-images"),
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
		GlobalFlagModel: utils.PtrValue(globalFlags),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, apiClient *ske.APIClient, model *inputModel) ske.ApiListProviderOptionsRequest {
	req := apiClient.ListProviderOptions(ctx, model.Region)
	return req
}

func outputResult(p *print.Printer, model *inputModel, options *ske.ProviderOptions) error {
	if options == nil {
		return fmt.Errorf("options is nil")
	}

	options.AvailabilityZones = nil
	options.KubernetesVersions = nil
	options.MachineTypes = nil
	options.VolumeTypes = nil

	return p.OutputResult(model.OutputFormat, options, func() error {
		images := utils.PtrValue(options.MachineImages)

		table := tables.NewTable()
		table.SetHeader("NAME", "VERSION", "STATE", "EXPIRATION DATE", "SUPPORTED CRI")
		for i := range images {
			image := images[i]
			versions := utils.PtrValue(image.Versions)
			for j := range versions {
				version := versions[j]
				criNames := make([]string, 0)
				for i := range utils.PtrValue(version.Cri) {
					cri := utils.PtrValue(version.Cri)[i]
					criNames = append(criNames, utils.PtrString(cri.Name))
				}
				criNamesString := strings.Join(criNames, ", ")

				expirationDate := "-"
				if version.ExpirationDate != nil {
					expirationDate = version.ExpirationDate.Format(time.RFC3339)
				}
				table.AddRow(
					utils.PtrString(image.Name),
					utils.PtrString(version.Version),
					utils.PtrString(version.State),
					expirationDate,
					criNamesString,
				)
			}
		}
		table.EnableAutoMergeOnColumns(1)

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("display output: %w", err)
		}
		return nil
	})
}
