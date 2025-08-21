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
	keyRingIdFlag = "key-ring"
	keyIdFlag     = "key"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string
	KeyId     string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all Key Versions",
		Long:  "Lists all versions of a given key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all key versions for the key "my-key-id" inside the key ring "my-key-ring-id"`,
				`$ stackit beta kms version list --key "my-key-id" --key-ring "my-key-ring-id"`),
			examples.NewExample(
				`List all key versions in JSON format`,
				`$ stackit beta kms version list --key "my-key-id" --key-ring "my-key-ring-id" -o json`),
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
				return fmt.Errorf("get Key Versions: %w", err)
			}
			if resp.Versions == nil || len(*resp.Versions) == 0 {
				params.Printer.Info("No Key Versions found for project %q in region %q for the key %q\n", model.ProjectId, model.Region, model.KeyId)
				return nil
			}
			keys := *resp.Versions

			return outputResult(params.Printer, model.OutputFormat, keys)
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
		KeyId:           flags.FlagToStringValue(p, cmd, keyIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiListVersionsRequest {
	return apiClient.ListVersions(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)
}

func outputResult(p *print.Printer, outputFormat string, versions []kms.Version) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(versions, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Key Versions list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(versions, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Key Versions list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("ID", "NUMBER", "CREATED AT", "DESTROY DATE", "STATUS")

		for i := range versions {
			version := versions[i]
			table.AddRow(
				utils.PtrString(version.KeyId),
				utils.PtrString(version.Number),
				utils.PtrString(version.CreatedAt),
				utils.PtrString(version.DestroyDate),
				utils.PtrString(version.State),
			)
		}

		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS Key Ring")
	cmd.Flags().Var(flags.UUIDFlag(), keyIdFlag, "ID of the Key")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, keyIdFlag)
	cobra.CheckErr(err)
}
