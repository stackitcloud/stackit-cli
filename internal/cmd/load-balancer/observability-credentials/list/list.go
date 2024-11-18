package list

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/load-balancer/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

const (
	instanceIdFlag = "instance-id"
	limitFlag      = "limit"
	usedFlag       = "used"
	unusedFlag     = "unused"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Limit  *int64
	Used   bool
	Unused bool
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists observability credentials for Load Balancer",
		Long:  "Lists observability credentials for Load Balancer.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all Load Balancer observability credentials`,
				"$ stackit load-balancer observability-credentials list"),
			examples.NewExample(
				`List all observability credentials being used by Load Balancer`,
				"$ stackit load-balancer observability-credentials list --used"),
			examples.NewExample(
				`List all observability credentials not being used by Load Balancer`,
				"$ stackit load-balancer observability-credentials list --unused"),
			examples.NewExample(
				`List all Load Balancer observability credentials in JSON format`,
				"$ stackit load-balancer observability-credentials list --output-format json"),
			examples.NewExample(
				`List up to 10 Load Balancer observability credentials`,
				"$ stackit load-balancer observability-credentials list --limit 10"),
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list Load Balancer observability credentials: %w", err)
			}
			credentialsPtr := resp.Credentials

			var credentials []loadbalancer.CredentialsResponse
			if credentialsPtr != nil && len(*credentialsPtr) > 0 {
				credentials = *credentialsPtr
				filterOp, err := getFilterOp(model.Used, model.Unused)
				if err != nil {
					return err
				}
				credentials, err = utils.FilterCredentials(ctx, apiClient, credentials, model.ProjectId, filterOp)
				if err != nil {
					return fmt.Errorf("filter credentials: %w", err)
				}
			}

			if len(credentials) == 0 {
				opLabel := "No "
				if model.Used {
					opLabel += "used"
				} else if model.Unused {
					opLabel += "unused"
				}
				p.Info("%s observability credentials found for Load Balancer on project %q\n", opLabel, projectLabel)
				return nil
			}

			// Truncate output
			if model.Limit != nil && len(credentials) > int(*model.Limit) {
				credentials = credentials[:*model.Limit]
			}
			return outputResult(p, model.OutputFormat, credentials)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Bool(usedFlag, false, "List only credentials being used by a Load Balancer")
	cmd.Flags().Bool(unusedFlag, false, "List only credentials not being used by a Load Balancer")

	cmd.MarkFlagsMutuallyExclusive(usedFlag, unusedFlag)
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

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Limit:           limit,
		Used:            flags.FlagToBoolValue(p, cmd, usedFlag),
		Unused:          flags.FlagToBoolValue(p, cmd, unusedFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *loadbalancer.APIClient) loadbalancer.ApiListCredentialsRequest {
	req := apiClient.ListCredentials(ctx, model.ProjectId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, credentials []loadbalancer.CredentialsResponse) error {
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(credentials, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Load Balancer observability credentials list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(credentials, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal Load Balancer observability credentials list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("REFERENCE", "DISPLAY NAME", "USERNAME")
		for i := range credentials {
			c := credentials[i]
			table.AddRow(*c.CredentialsRef, *c.DisplayName, *c.Username)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}

func getFilterOp(used, unused bool) (int, error) {
	// should not happen, cobra handles this
	if used && unused {
		return 0, fmt.Errorf("used and unused flags are mutually exclusive")
	}

	if !used && !unused {
		return utils.OP_FILTER_NOP, nil
	}

	if used {
		return utils.OP_FILTER_USED, nil
	}

	return utils.OP_FILTER_UNUSED, nil
}
