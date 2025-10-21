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
	keyRingIdFlag = "keyring-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all KMS keys",
		Long:  "List all KMS keys inside a key ring.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all KMS keys for the key ring "MY_KEYRING_ID"`,
				`$ stackit beta kms key list --keyring-id "MY_KEYRING_ID"`),
			examples.NewExample(
				`List all KMS keys in JSON format`,
				`$ stackit beta kms key list --keyring-id "MY_KEYRING_ID" --output-format json`),
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
				return fmt.Errorf("get KMS Keys: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.ProjectId, model.KeyRingId, resp)
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiListKeysRequest {
	req := apiClient.ListKeys(ctx, model.ProjectId, model.Region, model.KeyRingId)
	return req
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS Key Ring where the Key is stored")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag)
	cobra.CheckErr(err)
}

func outputResult(p *print.Printer, outputFormat, projectId, keyRingId string, resp *kms.KeyList) error {
	if resp == nil || resp.Keys == nil {
		return fmt.Errorf("response was nil / empty")
	}

	keys := *resp.Keys

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(keys, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal KMS Keys list: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(keys, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal KMS Keys list: %w", err)
		}
		p.Outputln(string(details))

	default:
		if len(keys) == 0 {
			p.Outputf("No keys found for project %q under the key ring %q\n", projectId, keyRingId)
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("ID", "NAME", "SCOPE", "ALGORITHM", "DELETION DATE", "STATUS")

		for _, key := range keys {
			table.AddRow(
				utils.PtrString(key.Id),
				utils.PtrString(key.DisplayName),
				utils.PtrString(key.Purpose),
				utils.PtrString(key.Algorithm),
				utils.PtrString(key.DeletionDate),
				utils.PtrString(key.State),
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}
	}
	return nil
}
