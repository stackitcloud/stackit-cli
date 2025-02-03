package list

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	limitFlag         = "limit"
	labelSelectorFlag = "label-selector"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit         *int64
	LabelSelector *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all key pairs",
		Long:  "Lists all key pairs.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Lists all key pairs`,
				"$ stackit beta key-pair list",
			),
			examples.NewExample(
				`Lists all key pairs which contains the label xxx`,
				"$ stackit beta key-pair list --label-selector xxx",
			),
			examples.NewExample(
				`Lists all key pairs in JSON format`,
				"$ stackit beta key-pair list --output-format json",
			),
			examples.NewExample(
				`Lists up to 10 key pairs`,
				"$ stackit beta key-pair list --limit 10",
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list key pairs: %w", err)
			}

			if resp.Items == nil || len(*resp.Items) == 0 {
				p.Info("No key pairs found\n")
				return nil
			}

			items := *resp.Items
			if model.Limit != nil && len(items) > int(*model.Limit) {
				items = items[:*model.Limit]
			}

			return outputResult(p, model.OutputFormat, items)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Number of key pairs to list")
	cmd.Flags().String(labelSelectorFlag, "", "Filter by label")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		LabelSelector:   flags.FlagToStringPointer(p, cmd, labelSelectorFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.InfoLevel, modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiListKeyPairsRequest {
	req := apiClient.ListKeyPairs(ctx)
	if model.LabelSelector != nil {
		req = req.LabelSelector(*model.LabelSelector)
	}
	return req
}

func outputResult(p *print.Printer, outputFormat string, keyPairs []iaas.Keypair) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(keyPairs, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal key pairs: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(keyPairs, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal key pairs: %w", err)
		}
		p.Outputln(string(details))

	default:
		table := tables.NewTable()
		table.SetHeader("KEY PAIR NAME", "LABELS", "FINGERPRINT", "CREATED AT", "UPDATED AT")

		for idx := range keyPairs {
			keyPair := keyPairs[idx]

			var labels []string
			for key, value := range *keyPair.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}

			table.AddRow(
				utils.PtrString(keyPair.Name),
				strings.Join(labels, ", "),
				utils.PtrString(keyPair.Fingerprint),
				utils.PtrString(keyPair.CreatedAt),
				utils.PtrString(keyPair.UpdatedAt),
			)
		}

		p.Outputln(table.Render())
	}
	return nil
}
