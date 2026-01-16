// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/client"
	commonErr "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/error"
	commonInstance "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/instance"
	commonKubeconfig "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/kubeconfig"
	commonValidation "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/validation"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-sdk-go/core/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/edge"
	"github.com/stackitcloud/stackit-sdk-go/services/edge/wait"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	identifier     *commonValidation.Identifier
	DisableWriting bool
	Filepath       *string
	Overwrite      bool
	Expiration     uint64
	SwitchContext  bool
}

// createRequestSpec captures the details of the request for testing.
type createRequestSpec struct {
	// Exported fields allow tests to inspect the request inputs
	ProjectID    string
	Region       string
	InstanceId   string
	InstanceName string
	Expiration   int64

	// Execute is a closure that wraps the actual SDK call
	Execute func() (*edge.Kubeconfig, error)
}

// OpenApi generated code will have different types for by-instance-id and by-display-name API calls and therefore different wait handlers.
// KubeconfigWaiter is an interface to abstract the different wait handlers so they can be used interchangeably.
type kubeconfigWaiter interface {
	WaitWithContext(context.Context) (*edge.Kubeconfig, error)
}

// A function that creates a kubeconfig waiter
type kubeconfigWaiterFactory = func(client *edge.APIClient) kubeconfigWaiter

// waiterFactoryProvider is an interface that provides kubeconfig waiters so we can inject different impl. while testing.
type waiterFactoryProvider interface {
	getKubeconfigWaiter(ctx context.Context, model *inputModel, apiClient client.APIClient) (kubeconfigWaiter, error)
}

// productionWaiterFactoryProvider is the real implementation used in production.
// It handles the concrete client type casting required by the SDK's wait handlers.
type productionWaiterFactoryProvider struct{}

func (p *productionWaiterFactoryProvider) getKubeconfigWaiter(ctx context.Context, model *inputModel, apiClient client.APIClient) (kubeconfigWaiter, error) {
	waiterFactory, err := getWaiterFactory(ctx, model)
	if err != nil {
		return nil, err
	}
	// The waiter handler needs a concrete client type. We can safely cast here as the real implementation will always match.
	edgeClient, ok := apiClient.(*edge.APIClient)
	if !ok {
		return nil, cliErr.NewBuildRequestError("failed to configure API client", nil)
	}
	return waiterFactory(edgeClient), nil
}

// waiterProvider is the package-level variable used to get the waiter.
// It is initialized with the production implementation but can be overridden in tests.
var waiterProvider waiterFactoryProvider = &productionWaiterFactoryProvider{}

// Command constructor
// Instance id and displayname are likely to be refactored in future. For the time being we decided to use flags
// instead of args to provide the instance-id xor displayname to uniquely identify an instance. The displayname
// is guaranteed to be unique within a given project as of today. The chosen flag over args approach ensures we
// won't need a breaking change of the CLI when we refactor the commands to take the identifier as arg at some point.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates or updates a local kubeconfig file of an edge instance",
		Long: fmt.Sprintf("%s\n\n%s\n%s\n%s\n%s",
			"Creates or updates a local kubeconfig file of a STACKIT Edge Cloud (STEC) instance. If the config exists in the kubeconfig file, the information will be updated.",
			"By default, the kubeconfig information of the edge instance is merged into the current kubeconfig file which is determined by Kubernetes client logic. If the kubeconfig file doesn't exist, a new one will be created.",
			fmt.Sprintf("You can override this behavior by specifying a custom filepath with the --%s flag or disable writing with the --%s flag.", commonKubeconfig.FilepathFlag, commonKubeconfig.DisableWritingFlag),
			fmt.Sprintf("An expiration time can be set for the kubeconfig. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is %d seconds.", commonKubeconfig.ExpirationSecondsDefault),
			"Note: the format for the duration is <value><unit>, e.g. 30d for 30 days. You may not combine units."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				fmt.Sprintf(`Create or update a kubeconfig for the edge instance with %s "xxx". If the config exists in the kubeconfig file, the information will be updated.`, commonInstance.InstanceIdFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud kubeconfig create --%s "xxx"`, commonInstance.InstanceIdFlag)),
			examples.NewExample(
				fmt.Sprintf(`Create or update a kubeconfig for the edge instance with %s "xxx" in a custom filepath.`, commonInstance.DisplayNameFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud kubeconfig create --%s "xxx" --filepath "yyy"`, commonInstance.DisplayNameFlag)),
			examples.NewExample(
				fmt.Sprintf(`Get a kubeconfig for the edge instance with %s "xxx" without writing it to a file and format the output as json.`, commonInstance.DisplayNameFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud kubeconfig create --%s "xxx" --disable-writing --output-format json`, commonInstance.DisplayNameFlag)),
			examples.NewExample(
				fmt.Sprintf(`Create a kubeconfig for the edge instance with %s "xxx". This will replace your current kubeconfig file.`, commonInstance.InstanceIdFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud kubeconfig create --%s "xxx" --overwrite`, commonInstance.InstanceIdFlag)),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()

			// Parse user input (arguments and/or flags)
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			// Prompt for confirmation is handled in outputResult

			if model.Async {
				return fmt.Errorf("async mode is not supported for kubeconfig create")
			}

			// Call API via waiter (which handles both the API call and waiting)
			kubeconfig, err := run(ctx, model, apiClient)
			if err != nil {
				return err
			}

			// Handle file operations or output to printer
			return outputResult(params.Printer, model.OutputFormat, model, kubeconfig)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(commonInstance.InstanceIdFlag, commonInstance.InstanceIdShorthand, "", commonInstance.InstanceIdUsage)
	cmd.Flags().StringP(commonInstance.DisplayNameFlag, commonInstance.DisplayNameShorthand, "", commonInstance.DisplayNameUsage)
	cmd.Flags().Bool(commonKubeconfig.DisableWritingFlag, false, commonKubeconfig.DisableWritingUsage)
	cmd.Flags().StringP(commonKubeconfig.FilepathFlag, commonKubeconfig.FilepathShorthand, "", commonKubeconfig.FilepathUsage)
	cmd.Flags().StringP(commonKubeconfig.ExpirationFlag, commonKubeconfig.ExpirationShorthand, "", commonKubeconfig.ExpirationUsage)
	cmd.Flags().Bool(commonKubeconfig.OverwriteFlag, false, commonKubeconfig.OverwriteUsage)
	cmd.Flags().Bool(commonKubeconfig.SwitchContextFlag, false, commonKubeconfig.SwitchContextUsage)

	identifierFlags := []string{commonInstance.InstanceIdFlag, commonInstance.DisplayNameFlag}
	cmd.MarkFlagsMutuallyExclusive(identifierFlags...) // InstanceId xor DisplayName
	cmd.MarkFlagsOneRequired(identifierFlags...)
	cmd.MarkFlagsMutuallyExclusive(commonKubeconfig.DisableWritingFlag, commonKubeconfig.FilepathFlag)  // DisableWriting xor Filepath
	cmd.MarkFlagsMutuallyExclusive(commonKubeconfig.DisableWritingFlag, commonKubeconfig.OverwriteFlag) // DisableWriting xor Overwrite
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
		Filepath:        flags.FlagToStringPointer(p, cmd, commonKubeconfig.FilepathFlag),
		Overwrite:       flags.FlagToBoolValue(p, cmd, commonKubeconfig.OverwriteFlag),
		SwitchContext:   flags.FlagToBoolValue(p, cmd, commonKubeconfig.SwitchContextFlag),
	}

	// Parse and validate user input then add it to the model
	id, err := commonValidation.GetValidatedInstanceIdentifier(p, cmd)
	if err != nil {
		return nil, err
	}
	model.identifier = id

	// Parse and validate kubeconfig expiration time
	if expString := flags.FlagToStringPointer(p, cmd, commonKubeconfig.ExpirationFlag); expString != nil {
		expTime, err := utils.ConvertToSeconds(*expString)
		if err != nil {
			return nil, &cliErr.FlagValidationError{
				Flag:    commonKubeconfig.ExpirationFlag,
				Details: err.Error(),
			}
		}
		if err := commonKubeconfig.ValidateExpiration(&expTime); err != nil {
			return nil, &cliErr.FlagValidationError{
				Flag:    commonKubeconfig.ExpirationFlag,
				Details: err.Error(),
			}
		}
		model.Expiration = expTime
	} else {
		// Default expiration is 1 hour
		defaultExp := uint64(commonKubeconfig.ExpirationSecondsDefault)
		model.Expiration = defaultExp
	}

	disableWriting := flags.FlagToBoolValue(p, cmd, commonKubeconfig.DisableWritingFlag)
	model.DisableWriting = disableWriting
	// Make sure to only output if the format is explicitly set
	if disableWriting {
		if globalFlags.OutputFormat == "" || globalFlags.OutputFormat == print.NoneOutputFormat {
			return nil, &cliErr.FlagValidationError{
				Flag:    commonKubeconfig.DisableWritingFlag,
				Details: fmt.Sprintf("must be used with --%s", globalflags.OutputFormatFlag),
			}
		}
		if globalFlags.OutputFormat != print.JSONOutputFormat && globalFlags.OutputFormat != print.YAMLOutputFormat {
			return nil, &cliErr.FlagValidationError{
				Flag:    globalflags.OutputFormatFlag,
				Details: fmt.Sprintf("valid output formats for this command are: %s", fmt.Sprintf("%s, %s", print.JSONOutputFormat, print.YAMLOutputFormat)),
			}
		}
	}

	// Log the parsed model if --verbosity is set to debug
	p.DebugInputModel(model)
	return &model, nil
}

// Run is the main execution function used by the command runner.
// It is decoupled from TTY output to have the ability to mock the API client during testing.
func run(ctx context.Context, model *inputModel, apiClient client.APIClient) (*edge.Kubeconfig, error) {
	spec, err := buildRequest(ctx, model, apiClient)
	if err != nil {
		return nil, err
	}

	resp, err := spec.Execute()
	if err != nil {
		return nil, cliErr.NewRequestFailedError(err)
	}

	return resp, nil
}

// buildRequest constructs the spec that can be tested.
func buildRequest(ctx context.Context, model *inputModel, apiClient client.APIClient) (*createRequestSpec, error) {
	if model == nil || model.identifier == nil {
		return nil, commonErr.NewNoIdentifierError("")
	}

	spec := &createRequestSpec{
		ProjectID:  model.ProjectId,
		Region:     model.Region,
		Expiration: int64(model.Expiration), // #nosec G115 ValidateExpiration ensures safe bounds, conversion is safe
	}

	switch model.identifier.Flag {
	case commonInstance.InstanceIdFlag:
		spec.InstanceId = model.identifier.Value
	case commonInstance.DisplayNameFlag:
		spec.InstanceName = model.identifier.Value
	default:
		return nil, fmt.Errorf("%w: %w", cliErr.NewBuildRequestError("invalid identifier flag", nil), commonErr.NewInvalidIdentifierError(model.identifier.Flag))
	}

	// Closure used to decouple the actual SDK call for easier testing
	spec.Execute = func() (*edge.Kubeconfig, error) {
		// Get the waiter from the provider (handles client type casting internally)
		waiter, err := waiterProvider.getKubeconfigWaiter(ctx, model, apiClient)
		if err != nil {
			return nil, err
		}

		return waiter.WaitWithContext(ctx)
	}

	return spec, nil
}

// Returns a factory function to create the appropriate waiter based on the input model.
func getWaiterFactory(ctx context.Context, model *inputModel) (kubeconfigWaiterFactory, error) {
	if model == nil || model.identifier == nil {
		return nil, commonErr.NewNoIdentifierError("")
	}

	// The KubeconfigWaitHandlers don't wait for the kubeconfig to be created, but for the instance to be ready to return a kubeconfig.
	// Convert uint64 to int64 to match the API's type.
	var expiration = int64(model.Expiration) // #nosec G115 ValidateExpiration ensures safe bounds, conversion is safe
	switch model.identifier.Flag {
	case commonInstance.InstanceIdFlag:
		factory := func(c *edge.APIClient) kubeconfigWaiter {
			return wait.KubeconfigWaitHandler(ctx, c, model.ProjectId, model.Region, model.identifier.Value, &expiration)
		}
		return factory, nil
	case commonInstance.DisplayNameFlag:
		factory := func(c *edge.APIClient) kubeconfigWaiter {
			return wait.KubeconfigByInstanceNameWaitHandler(ctx, c, model.ProjectId, model.Region, model.identifier.Value, &expiration)
		}
		return factory, nil
	default:
		return nil, commonErr.NewInvalidIdentifierError(model.identifier.Flag)
	}
}

// Output result based on the configured output format
func outputResult(p *print.Printer, outputFormat string, model *inputModel, kubeconfig *edge.Kubeconfig) error {
	// Ensure kubeconfig data is present
	if kubeconfig == nil || kubeconfig.Kubeconfig == nil {
		return fmt.Errorf("no kubeconfig returned from the API")
	}
	kubeconfigMap := *kubeconfig.Kubeconfig

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
	kubeconfigData, err := marshalKubeconfig(kubeconfigMap, format)
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
		p.Outputf("Wrote kubeconfig for instance %q to %q.\n", model.identifier.Value, *path)

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
		return "", fmt.Errorf("%w: %s", commonErr.NewNoIdentifierError(""), format)
	}
}
