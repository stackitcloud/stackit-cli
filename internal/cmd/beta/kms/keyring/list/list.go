package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all KMS Keyrings",
		Long:  "Lists all KMS Keyrings.",
		Args:  args.NoArgs,
		Example: examples.Build(
			// Enforce a specific region for the KMS
			examples.NewExample(
				`List all KMS Keyrings`,
				"$ stackit beta kms keyring list"),
			examples.NewExample(
				`List all KMS Keyrings in JSON format`,
				"$ stackit beta kms keyring list --output-format json"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
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
				return fmt.Errorf("get KMS Keyrings: %w", err)
			}
			if resp.KeyRings == nil || len(*resp.KeyRings) == 0 {
				params.Printer.Info("No Keyrings found for project %q in region %q\n", model.ProjectId, model.Region)
				return nil
			}
			keyRings := *resp.KeyRings

			return outputResult(params.Printer, model.OutputFormat, keyRings)
		},
	}

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiListKeyRingsRequest {
	req := apiClient.ListKeyRings(ctx, model.ProjectId, model.Region)
	return req
}

func outputResult(p *print.Printer, outputFormat string, keyRings []kms.KeyRing) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(keyRings, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal KMS Keyrings list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(keyRings, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS Keyrings list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "STATUS")

		for i := range keyRings {
			keyRing := keyRings[i]
			table.AddRow(
				utils.PtrString(keyRing.Id),
				utils.PtrString(keyRing.DisplayName),
				utils.PtrString(keyRing.State),
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
