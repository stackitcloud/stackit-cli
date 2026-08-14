package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	edge "github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api"
	"github.com/stackitcloud/stackit-sdk-go/services/edge/v1beta1api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	edgeUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	sdkUtils "github.com/stackitcloud/stackit-sdk-go/core/utils"

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
	instanceIdFlag     = "instance-id"
	expirationFlag     = "expiration"
	disableWritingFlag = "disable-writing"
	filepathFlag       = "filepath"
	overwriteFlag      = "overwrite"
	switchContextFlag  = "switch-context"

	expirationSecondsDefault = 3600 // 60 * 60 seconds = 1 hour
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId     string
	DisableWriting bool
	Filepath       *string
	Overwrite      bool
	Expiration     uint64
	SwitchContext  bool
}

// NewCmd https://aip.stackit.cloud/aip/general/0121/
// We have decided to eliminate the usage of display name flag
// To be the AIP compliant, and align with the standard CLI implementation, we will use the InstanceID arg
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates or updates a local kubeconfig file of an Edge Cloud instance",
		Long: fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
			"Creates or updates a local kubeconfig file of a STACKIT Edge Cloud (STEC) instance. If the config exists in the kubeconfig file, the information will be updated.",
			"By default, the kubeconfig information of the edge instance is merged into the current kubeconfig file which is determined by Kubernetes client logic. If the kubeconfig file doesn't exist, a new one will be created.",
			fmt.Sprintf("You can override this behavior by specifying a custom filepath with the --%s flag or disable writing with the --%s flag.", filepathFlag, disableWritingFlag),
			fmt.Sprintf("An expiration time can be set for the kubeconfig. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is %d seconds.", expirationSecondsDefault),
			"Note: the format for the duration is <value><unit>, e.g. 30d for 30 days. You may not combine units."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create or update a kubeconfig for the Edge Cloud instance with instance ID "xxx". If the config exists in the kubeconfig file, the information will be updated.`,
				`$ stackit beta edge-cloud kubeconfig create --instance-id xxx`),
			examples.NewExample(
				`Create or update a kubeconfig for the Edge Cloud instance with instance ID "xxx" in a custom filepath.`,
				`$ stackit beta edge-cloud kubeconfig create --instance-id xxx --filepath yyy`),
			examples.NewExample(
				`Get a kubeconfig for the Edge Cloud instance with instance ID "xxx" without writing it to a file and format the output as json.`,
				`$ stackit beta edge-cloud kubeconfig create --instance-id xxx --disable-writing --output-format json`),
			examples.NewExample(
				`Create a kubeconfig for the Edge Cloud instance with instance ID "xxx". This will replace your current kubeconfig file.`,
				`$ stackit beta edge-cloud kubeconfig create --instance-id xxx --overwrite`),
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
				return fmt.Errorf("async mode is not supported for kubeconfig create")
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

			if !model.DisableWriting {
				var prompt string
				if model.Overwrite {
					prompt = fmt.Sprintf("Are you sure you want to create a kubeconfig for instance %q of project %q? This will OVERWRITE your current kubeconfig file, if it exists.", instanceLabel, projectLabel)
				} else {
					prompt = fmt.Sprintf("Are you sure you want to update your kubeconfig for instance %q of project %q? This will update your kubeconfig file. \nIf the kubeconfig file does not exist, it will create a new one.", instanceLabel, projectLabel)
				}
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("build kubeconfig create request: %w", err)
			}
			respKubeconfig, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create kubeconfig for Edge Cloud instance: %w", err)
			}
			if respKubeconfig == nil {
				return fmt.Errorf("no kubeconfig returned from the API")
			}

			var expiration = int64(model.Expiration) // #nosec G115 ValidateExpiration ensures safe bounds, conversion is safe
			err = spinner.Run(params.Printer, "Creating kubeconfig", func() error {
				_, err = wait.KubeconfigWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, model.InstanceId, &expiration).WaitWithContext(ctx)
				return err
			})
			if err != nil {
				return fmt.Errorf("wait for kubeconfig creation: %w", err)
			}

			// Handle file operations or output to printer
			return outputResult(params.Printer, model.OutputFormat, model, respKubeconfig)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(instanceIdFlag, "", "Edge Cloud instance ID")
	cmd.Flags().Bool(disableWritingFlag, false, "Disable writing the kubeconfig to a file.")
	cmd.Flags().StringP(filepathFlag, "f", "", "Path to the kubeconfig file. A default is chosen by Kubernetes if not set.")
	cmd.Flags().StringP(expirationFlag, "e", "", "Expiration time for the kubeconfig, e.g. 5d. By default, the token is valid for 1h.")
	cmd.Flags().Bool(overwriteFlag, false, "Force overwrite the kubeconfig file if it exists.")
	cmd.Flags().Bool(switchContextFlag, false, "Switch to the context in the kubeconfig file to the new context.")

	cmd.MarkFlagsMutuallyExclusive(disableWritingFlag, filepathFlag)  // DisableWriting xor Filepath
	cmd.MarkFlagsMutuallyExclusive(disableWritingFlag, overwriteFlag) // DisableWriting xor Overwrite

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
		Filepath:        flags.FlagToStringPointer(p, cmd, filepathFlag),
		Overwrite:       flags.FlagToBoolValue(p, cmd, overwriteFlag),
		SwitchContext:   flags.FlagToBoolValue(p, cmd, switchContextFlag),
	}

	// Parse and validate kubeconfig expiration time
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

	disableWriting := flags.FlagToBoolValue(p, cmd, disableWritingFlag)
	model.DisableWriting = disableWriting
	// Make sure to only output if the format is explicitly set
	if disableWriting {
		if globalFlags.OutputFormat == "" || globalFlags.OutputFormat == print.NoneOutputFormat {
			return nil, &cliErr.FlagValidationError{
				Flag:    disableWritingFlag,
				Details: fmt.Sprintf("must be used with --%s", globalflags.OutputFormatFlag.Name()),
			}
		}
		if globalFlags.OutputFormat != print.JSONOutputFormat && globalFlags.OutputFormat != print.YAMLOutputFormat {
			return nil, &cliErr.FlagValidationError{
				Flag:    globalflags.OutputFormatFlag.Name(),
				Details: fmt.Sprintf("valid output formats for this command are: %s", fmt.Sprintf("%s, %s", print.JSONOutputFormat, print.YAMLOutputFormat)),
			}
		}
	}

	// Log the parsed model if --verbosity is set to debug
	p.DebugInputModel(model)
	return &model, nil
}

// buildRequest constructs the spec that can be tested.
func buildRequest(ctx context.Context, model *inputModel, apiClient *edge.APIClient) (edge.ApiGetKubeconfigByInstanceIdRequest, error) {
	req := apiClient.DefaultAPI.GetKubeconfigByInstanceId(ctx, model.ProjectId, model.Region, model.InstanceId)
	return req.ExpirationSeconds(int64(model.Expiration)), nil // #nosec G115 ValidateExpiration ensures safe bounds, conversion is safe
}

func outputResult(p *print.Printer, outputFormat string, model *inputModel, kubeconfig *edge.Kubeconfig) error {
	if kubeconfig == nil || kubeconfig.Kubeconfig == nil {
		return fmt.Errorf("no kubeconfig returned from the API")
	}

	// Determine output format for terminal or file output
	var format string
	switch outputFormat {
	case print.JSONOutputFormat:
		// JSON if explicitly requested
		format = print.JSONOutputFormat
	case print.YAMLOutputFormat:
		// YAML if explicitly requested
		format = print.YAMLOutputFormat
	default:
		if model.DisableWriting {
			// If not explicitly requested, use JSON as default for terminal output
			format = print.JSONOutputFormat
		} else {
			// If not explicitly requested, use YAML as default for file output
			format = print.YAMLOutputFormat
		}
	}

	// Marshal kubeconfig data based on the determined format
	kubeconfigData, err := marshalKubeconfig(kubeconfig.Kubeconfig, format)
	if err != nil {
		return err
	}

	// Handle file writing and output
	if !model.DisableWriting {
		// Build options for writing kubeconfig
		opts := commonKubeconfig.NewWriteOptions().
			WithOverwrite(model.Overwrite).
			WithSwitchContext(model.SwitchContext)

		// Add confirmation callback if not assumeYes
		if !model.AssumeYes {
			confirmFn := func(message string) error {
				return p.PromptForConfirmation(message)
			}
			opts = opts.WithConfirmation(confirmFn)
		}

		path, err := commonKubeconfig.WriteKubeconfig(model.Filepath, kubeconfigData, opts)
		if err != nil {
			return err
		}

		// Inform the user about the successful write operation
		p.Outputf("Wrote kubeconfig for instance %q to %q.\n", model.InstanceId, *path)

		if model.SwitchContext {
			p.Outputln("Switched context as requested.")
		}
	} else {
		p.Outputln(kubeconfigData)
	}
	return nil
}

// Marshal kubeconfig data to the specified format
func marshalKubeconfig(kubeconfigMap map[string]interface{}, format string) (string, error) {
	switch format {
	case print.JSONOutputFormat:
		kubeconfigJSON, err := json.MarshalIndent(kubeconfigMap, "", "  ")
		if err != nil {
			return "", fmt.Errorf("marshal kubeconfig to JSON: %w", err)
		}
		return string(kubeconfigJSON), nil
	case print.YAMLOutputFormat:
		kubeconfigYAML, err := yaml.MarshalWithOptions(kubeconfigMap, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return "", fmt.Errorf("marshal kubeconfig to YAML: %w", err)
		}
		return string(kubeconfigYAML), nil
	default:
		return "", fmt.Errorf("format is not JSON or YAML: %s", format)
	}
}
