package list

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/redis/client"
	redisUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/redis/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/redis"
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
		Short: "Lists all credentials' IDs for a Redis instance",
		Long:  "Lists all credentials' IDs for a Redis instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all credentials' IDs for a Redis instance`,
				"$ stackit redis credentials list --instance-id xxx"),
			examples.NewExample(
				`List all credentials' IDs for a Redis instance in JSON format`,
				"$ stackit redis credentials list --instance-id xxx --output-format json"),
			examples.NewExample(
				`List up to 10 credentials' IDs for a Redis instance`,
				"$ stackit redis credentials list --instance-id xxx --limit 10"),
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
				return fmt.Errorf("list Redis credentials: %w", err)
			}
			credentials := *resp.CredentialsList

			// Truncate output
			if model.Limit != nil && len(credentials) > int(*model.Limit) {
				credentials = credentials[:*model.Limit]
			}

			instanceLabel := model.InstanceId
			if len(credentials) == 0 {
				instanceLabel, err = redisUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				}
			}

			return outputResult(params.Printer, model.OutputFormat, instanceLabel, credentials)
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

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Limit:           limit,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *redis.APIClient) redis.ApiListCredentialsRequest {
	req := apiClient.ListCredentials(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat, instanceLabel string, credentials []redis.CredentialsListItem) error {
	return p.OutputResult(outputFormat, credentials, func() error {
		if len(credentials) == 0 {
			p.Outputf("No credentials found for instance %q\n", instanceLabel)
			return nil
		}

		table := tables.NewTable()
		table.SetHeader("ID")
		for i := range credentials {
			c := credentials[i]
			table.AddRow(utils.PtrString(c.Id))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
