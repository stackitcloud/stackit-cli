package list

import (
	"context"
	"encoding/json"
	"fmt"

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
		Short: "Lists all observability credentials for Load Balancer",
		Long:  "Lists all observability credentials for Load Balancer.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List all observability credentials for Load Balancer`,
				"$ stackit load-balancer observability-credentials list"),
			examples.NewExample(
				`List all observability credentials for Load Balancer in JSON format`,
				"$ stackit load-balancer observability-credentials list --output-format json"),
			examples.NewExample(
				`List up to 10 observability credentials for Load Balancer`,
				"$ stackit load-balancer observability-credentials list --limit 10"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
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

			var usedCredentials []loadbalancer.CredentialsResponse
			if model.Used {
				usedCredentials, err = utils.GetUsedObsCredentials(ctx, apiClient, model.ProjectId)
				if err != nil {
					return fmt.Errorf("get used observability credentials: %w", err)
				}
				if len(usedCredentials) == 0 {
					p.Info("No used observability credentials found for Load Balancer on project %q\n", projectLabel)
					return nil
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("list Load Balancer observability credentials: %w", err)
			}
			credentialsPtr := resp.Credentials
			if credentialsPtr == nil || (credentialsPtr != nil && len(*credentialsPtr) == 0) {
				p.Info("No observability credentials found for Load Balancer on project %q\n", projectLabel)
				return nil
			}

			credentials := *credentialsPtr

			// Truncate output
			if model.Limit != nil && len(credentials) > int(*model.Limit) {
				credentials = credentials[:*model.Limit]
			}
			return outputResult(p, model, credentials, usedCredentials)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Int64(limitFlag, 0, "Maximum number of entries to list")
	cmd.Flags().Bool(usedFlag, false, "List only used credentials")
	cmd.Flags().Bool(unusedFlag, false, "List only unused credentials")
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

func outputResult(p *print.Printer, model *inputModel, credentials, usedCredentials []loadbalancer.CredentialsResponse) error {
	creds := credentials
	if model.Used {
		creds = usedCredentials
	}

	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(creds, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Load Balancer observability credentials list: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.SetHeader("REFERENCE", "DISPLAY NAME", "USERNAME")
		for i := range creds {
			c := creds[i]
			table.AddRow(*c.CredentialsRef, *c.DisplayName, *c.Username)
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
