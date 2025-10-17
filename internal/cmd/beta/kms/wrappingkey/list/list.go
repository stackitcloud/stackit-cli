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
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/kms/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"
)

const (
	keyRingIdFlag = "key-ring-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all KMS wrapping keys",
		Long:  "Lists all KMS wrapping keys inside a key ring.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all KMS wrapping keys for the key ring "my-key-ring-id"`,
				`$ stackit beta kms wrapping-key list --key-ring-id "my-key-ring-id"`),
			examples.NewExample(
				`List all KMS wrapping keys in JSON format`,
				`$ stackit beta kms wrappingkeys list --key-ring-id "my-key-ring-id" --output-format json`),
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
				return fmt.Errorf("get KMS wrapping keys: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.KeyRingId, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		KeyRingId:       flags.FlagToStringValue(p, cmd, keyRingIdFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiListWrappingKeysRequest {
	req := apiClient.ListWrappingKeys(ctx, model.ProjectId, model.Region, model.KeyRingId)
	return req
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS Key Ring where the Key is stored")
	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag)
	cobra.CheckErr(err)
}

func outputResult(p *print.Printer, outputFormat, keyRingId string, resp *kms.WrappingKeyList) error {
	if resp == nil || resp.WrappingKeys == nil {
		return fmt.Errorf("response is nil / empty")
	}

	wrappingKeys := *resp.WrappingKeys

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(wrappingKeys, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal KMS wrapping keys list: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(wrappingKeys, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS wrapping keys list: %w", err)
		}
		p.Outputln(string(details))

	default:
		if len(wrappingKeys) == 0 {
			p.Outputf("No wrapping keys found under the key ring %q\n", keyRingId)
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "SCOPE", "ALGORITHM", "EXPIRES AT", "STATUS")

		for i := range wrappingKeys {
			wrappingKey := wrappingKeys[i]
			table.AddRow(
				utils.PtrString(wrappingKey.Id),
				utils.PtrString(wrappingKey.DisplayName),
				utils.PtrString(wrappingKey.Purpose),
				utils.PtrString(wrappingKey.Algorithm),
				utils.PtrString(wrappingKey.ExpiresAt),
				utils.PtrString(wrappingKey.State),
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
	}

	return nil
}
