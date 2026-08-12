package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"
	"github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	edgeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	commonKubeconfig "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/kubeconfig"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
)

const (
	instanceIdFlag           = "instance-id"
	expirationFlag           = "expiration"
	expirationSecondsDefault = 3600 // 60 * 60 seconds = 1 hour
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Expiration uint64
	InstanceId string
}

// NewCmd https://aip.stackit.cloud/aip/general/0121/
// We have decided to eliminate the usage of display name flag
// To be the AIP compliant, and align with the standard CLI implementation, we will use the InstanceID arg
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a token for an Edge Cloud instance",
		Long: fmt.Sprintf("%s\n\n%s\n%s",
			"Creates a token for a STACKIT Edge Cloud (STEC) instance.",
			fmt.Sprintf("An expiration time can be set for the token. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is %d seconds.", expirationSecondsDefault),
			"Note: the format for the duration is <value><unit>, e.g. 30d for 30 days. You may not combine units."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a token for the Edge Cloud instance with instance ID "xxx".`,
				`$ stackit beta edge-cloud token create --instance-id xxx`),
			examples.NewExample(
				`Create a token for the Edge Cloud instance with instance ID "xxx". The token will be valid for one day.`,
				`$ stackit beta edge-cloud token create --instance-id xxx --expiration 1d`),
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

			if model.Async {
				return fmt.Errorf("async mode is not supported for token create")
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				// If project label can't be determined, fall back to project ID
				projectLabel = model.ProjectId
			}

			instanceLabel, err := edgeUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a new token for Edge Cloud instance %q for project %q?", instanceLabel, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create edge cloud instance token: %w", err)
			}

			var expiration = int64(model.Expiration) // #nosec G115 ValidateExpiration ensures safe bounds, conversion is safe
			err = spinner.Run(params.Printer, "Creating token", func() error {
				_, err = wait.TokenWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, &expiration).WaitWithContext(ctx)
				return err
			})
			if err != nil {
				return fmt.Errorf("wait for token creation: %w", err)
			}

			// Handle output to printer
			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(instanceIdFlag, "", "Edge Cloud instance ID")
	cmd.Flags().StringP(expirationFlag, "e", "", fmt.Sprintf("Expiration time for the token, e.g. 5d. By default, the token is valid for %d seconds.", expirationSecondsDefault))

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag)
	cobra.CheckErr(err)
}

// Parse user input (arguments and/or flags)
func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	// Generate input model based on chosen flags
	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
	}

	// Parse and validate token expiration time
	if expString := flags.FlagToStringPointer(p, cmd, expirationFlag); expString != nil {
		expTime, err := sdkUtils.ConvertToSeconds(*expString)
		if err != nil {
			return nil, &cliErr.FlagValidationError{
				Flag:    expirationFlag,
				Details: err.Error(),
			}
		}
		if err := commonKubeconfig.ValidateExpiration(&expTime); err != nil {
			return nil, &cliErr.FlagValidationError{
				Flag:    expirationFlag,
				Details: err.Error(),
			}
		}
		model.Expiration = expTime
	} else {
		// Default expiration is 1 hour
		defaultExp := uint64(expirationSecondsDefault)
		model.Expiration = defaultExp
	}

	// Make sure to only output if the format is not none
	if globalFlags.OutputFormat == print.NoneOutputFormat {
		return nil, &cliErr.FlagValidationError{
			Flag:    globalflags.OutputFormatFlag.Name(),
			Details: fmt.Sprintf("valid formats for this command are: %s", fmt.Sprintf("%s, %s, %s", print.PrettyOutputFormat, print.JSONOutputFormat, print.YAMLOutputFormat)),
		}
	}

	// Log the parsed model if --verbosity is set to debug
	p.DebugInputModel(model)
	return &model, nil
}

// buildRequest constructs the spec that can be tested.
func buildRequest(ctx context.Context, model *inputModel, apiClient *edge.APIClient) (edge.ApiGetTokenByInstanceIdRequest, error) {
	req := apiClient.DefaultAPI.GetTokenByInstanceId(ctx, model.ProjectId, model.Region, model.InstanceId)
	return req.ExpirationSeconds(int64(model.Expiration)), nil // #nosec G115 ValidateExpiration ensures safe bounds, conversion is safe
}

// Output result based on the configured output format
func outputResult(p *print.Printer, outputFormat string, token *edge.Token) error {
	return p.OutputResult(outputFormat, token, func() error {
		if token == nil {
			// This is only to prevent nil pointer deref.
			// As long as the API behaves as defined by it's spec, instance can not be empty (HTTP 200 with an empty body)
			return fmt.Errorf("no token returned from the API")
		}

		p.Outputln(token.Token)
		return nil
	})
}
