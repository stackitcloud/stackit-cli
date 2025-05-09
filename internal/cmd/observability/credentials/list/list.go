package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/observability"
)

const (
	instanceIdFlag = "instance-id"
	limitFlag      = "limit"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Limit      *int64
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists the usernames of all credentials for an Observability instance",
		Long:  "Lists the usernames of all credentials for an Observability instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List the usernames of all credentials for an Observability instance with ID "xxx"`,
				"$ stackit observability credentials list --instance-id xxx"),
			examples.NewExample(
				`List the usernames of all credentials for an Observability instance in JSON format`,
				"$ stackit observability credentials list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List the usernames of up to 10 credentials for an Observability instance`,
				"$ stackit observability credentials list --instance-id xxx --limit 10"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list Observability credentials: %w", err)
			}
			credentials := *resp.Credentials
			if len(credentials) == 0 {
				instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient, model.InstanceId, model.ProjectId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
					instanceLabel = model.InstanceId
				}
				params.Printer.Info("No credentials found for instance %q\n", instanceLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(credentials) > int(*model.Limit) {
				credentials = credentials[:*model.Limit]
			}
			return outputResult(params.Printer, model.OutputFormat, credentials)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	limit := flags.FlagToInt64Pointer(p, cmd, limitFlag)
	if limit != nil && *limit < 1 {
		return nil, &errors.FlagValidationError{
			Flag:    limitFlag,
			Details: "must be greater than 0",
		}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Limit:           limit,
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *observability.APIClient) observability.ApiListCredentialsRequest {
	req := apiClient.ListCredentials(ctx, model.InstanceId, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, credentials []observability.ServiceKeysList) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(credentials, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Observability credentials list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(credentials, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Observability credentials list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("USERNAME")
		for i := range credentials {
			c := credentials[i]
			table.AddRow(utils.PtrString(c.Name))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
