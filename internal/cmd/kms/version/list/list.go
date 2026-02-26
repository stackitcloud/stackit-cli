package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
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
	keyIdFlag     = "key-id"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	KeyRingId string
	KeyId     string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all key versions",
		Long:  "List all versions of a given key.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all key versions for the key "my-key-id" inside the key ring "my-keyring-id"`,
				`$ stackit kms version list --key-id "my-key-id" --keyring-id "my-keyring-id"`),
			examples.NewExample(
				`List all key versions in JSON format`,
				`$ stackit kms version list --key-id "my-key-id" --keyring-id "my-keyring-id" -o json`),
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
				return fmt.Errorf("get key version: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, model.ProjectId, model.KeyId, resp)
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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *kms.APIClient) kms.ApiListVersionsRequest {
	return apiClient.ListVersions(ctx, model.ProjectId, model.Region, model.KeyRingId, model.KeyId)
}

func outputResult(p *print.Printer, outputFormat, projectId, keyId string, resp *kms.VersionList) error {
	if resp == nil || resp.Versions == nil {
		return fmt.Errorf("response is nil / empty")
	}
	versions := *resp.Versions

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(versions, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal key versions list: %w", err)
		}
		p.Outputln(string(details))

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(versions, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal key versions list: %w", err)
		}
		p.Outputln(string(details))

	default:
		if len(versions) == 0 {
			p.Outputf("No key versions found for project %q for the key %q\n", projectId, keyId)
			return nil
		}
		table := tables.NewTable()
		table.SetHeader("ID", "NUMBER", "CREATED AT", "DESTROY DATE", "STATUS")

		for _, version := range versions {
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
	}

	return nil
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), keyRingIdFlag, "ID of the KMS key ring")
	cmd.Flags().Var(flags.UUIDFlag(), keyIdFlag, "ID of the key")

	err := flags.MarkFlagsRequired(cmd, keyRingIdFlag, keyIdFlag)
	cobra.CheckErr(err)
}
