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

const (
	keyRingIdArg = "KEYRING_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	keyRingId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("list %s", keyRingIdArg),
		Short: "Lists all KMS Keys",
		Long:  "Lists all KMS Keys inside a key ring.",
		Args:  args.SingleArg(keyRingIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`List all KMS Keys for the key ring "xxx"`,
				"$ stackit beta kms key list xxx"),
			examples.NewExample(
				`List all KMS Keys in JSON format`,
				"$ stackit beta kms key list xxx --output-format json"),
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
				return fmt.Errorf("get KMS Keys: %w", err)
			}
			if resp.Keys == nil || len(*resp.Keys) == 0 {
				params.Printer.Info("No Keys found for project %q in region %q under the key ring %q\n", model.ProjectId, model.Region, model.keyRingId)
				return nil
			}
			keys := *resp.Keys

			return outputResult(params.Printer, model.OutputFormat, keys)
		},
	}

	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	keyRingId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		keyRingId:       keyRingId,
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiListKeysRequest {
	req := apiClient.ListKeys(ctx, model.ProjectId, model.Region, model.keyRingId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, keys []kms.Key) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(keys, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal KMS Keys list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(keys, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS Keys list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "SCOPE", "ALGORITHM", "DELETION DATE", "STATUS")

		for i := range keys {
			key := keys[i]
			table.AddRow(
				utils.PtrString(key.Id),
				utils.PtrString(key.DisplayName),
				utils.PtrString(key.Purpose),
				utils.PtrString(key.Algorithm),
				// utils.PtrString(wrappingKeys.CreatedAt),
				utils.PtrString(key.DeletionDate),
				utils.PtrString(key.State),
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
