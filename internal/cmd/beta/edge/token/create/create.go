// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package create

import (
	"context"
	"fmt"

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
	identifier *commonValidation.Identifier
	Expiration uint64
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
	Execute func() (*edge.Token, error)
}

// OpenApi generated code will have different types for by-instance-id and by-display-name API calls and therefore different wait handlers.
// tokenWaiter is an interface to abstract the different wait handlers so they can be used interchangeably.
type tokenWaiter interface {
	WaitWithContext(context.Context) (*edge.Token, error)
}

// A function that creates a token waiter
type tokenWaiterFactory = func(client *edge.APIClient) tokenWaiter

// waiterFactoryProvider is an interface that provides token waiters so we can inject different impl. while testing.
type waiterFactoryProvider interface {
	getTokenWaiter(ctx context.Context, model *inputModel, apiClient client.APIClient) (tokenWaiter, error)
}

// productionWaiterFactoryProvider is the real implementation used in production.
// It handles the concrete client type casting required by the SDK's wait handlers.
type productionWaiterFactoryProvider struct{}

func (p *productionWaiterFactoryProvider) getTokenWaiter(ctx context.Context, model *inputModel, apiClient client.APIClient) (tokenWaiter, error) {
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
		Short: "Creates a token for an edge instance",
		Long: fmt.Sprintf("%s\n\n%s\n%s",
			"Creates a token for a STACKIT Edge Cloud (STEC) instance.",
			fmt.Sprintf("An expiration time can be set for the token. The expiration time is set in seconds(s), minutes(m), hours(h), days(d) or months(M). Default is %d seconds.", commonKubeconfig.ExpirationSecondsDefault),
			"Note: the format for the duration is <value><unit>, e.g. 30d for 30 days. You may not combine units."),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				fmt.Sprintf(`Create a token for the edge instance with %s "xxx".`, commonInstance.InstanceIdFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud token create --%s "xxx"`, commonInstance.InstanceIdFlag)),
			examples.NewExample(
				fmt.Sprintf(`Create a token for the edge instance with %s "xxx". The token will be valid for one day.`, commonInstance.DisplayNameFlag),
				fmt.Sprintf(`$ stackit beta edge-cloud token create --%s "xxx" --expiration 1d`, commonInstance.DisplayNameFlag)),
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

			if model.Async {
				return fmt.Errorf("async mode is not supported for token create")
			}

			// Call API
			resp, err := run(ctx, model, apiClient)
			if err != nil {
				return err
			}

			// Handle output to printer
			return outputResult(params.Printer, model.OutputFormat, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(commonInstance.InstanceIdFlag, commonInstance.InstanceIdShorthand, "", commonInstance.InstanceIdUsage)
	cmd.Flags().StringP(commonInstance.DisplayNameFlag, commonInstance.DisplayNameShorthand, "", commonInstance.DisplayNameUsage)
	cmd.Flags().StringP(commonKubeconfig.ExpirationFlag, commonKubeconfig.ExpirationShorthand, "", commonKubeconfig.ExpirationUsage)

	identifierFlags := []string{commonInstance.InstanceIdFlag, commonInstance.DisplayNameFlag}
	cmd.MarkFlagsMutuallyExclusive(identifierFlags...) // InstanceId xor DisplayName
	cmd.MarkFlagsOneRequired(identifierFlags...)
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

	// Make sure to only output if the format is not none
	if globalFlags.OutputFormat == print.NoneOutputFormat {
		return nil, &cliErr.FlagValidationError{
			Flag:    globalflags.OutputFormatFlag,
			Details: fmt.Sprintf("valid formats for this command are: %s", fmt.Sprintf("%s, %s, %s", print.PrettyOutputFormat, print.JSONOutputFormat, print.YAMLOutputFormat)),
		}
	}

	// Log the parsed model if --verbosity is set to debug
	p.DebugInputModel(model)
	return &model, nil
}

// Run is the main execution function used by the command runner.
// It is decoupled from TTY output to have the ability to mock the API client during testing.
func run(ctx context.Context, model *inputModel, apiClient client.APIClient) (*edge.Token, error) {
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
	spec.Execute = func() (*edge.Token, error) {
		// Get the waiter from the provider (handles client type casting internally)
		waiter, err := waiterProvider.getTokenWaiter(ctx, model, apiClient)
		if err != nil {
			return nil, err
		}

		return waiter.WaitWithContext(ctx)
	}

	return spec, nil
}

// Returns a factory function to create the appropriate waiter based on the input model.
func getWaiterFactory(ctx context.Context, model *inputModel) (tokenWaiterFactory, error) {
	if model == nil || model.identifier == nil {
		return nil, commonErr.NewNoIdentifierError("")
	}

	// The tokenWaitHandlers don't wait for the token to be created, but for the instance to be ready to return a token.
	// Convert uint64 to int64 to match the API's type.
	var expiration = int64(model.Expiration) // #nosec G115 ValidateExpiration ensures safe bounds, conversion is safe
	switch model.identifier.Flag {
	case commonInstance.InstanceIdFlag:
		factory := func(c *edge.APIClient) tokenWaiter {
			return wait.TokenWaitHandler(ctx, c, model.ProjectId, model.Region, model.identifier.Value, &expiration)
		}
		return factory, nil
	case commonInstance.DisplayNameFlag:
		factory := func(c *edge.APIClient) tokenWaiter {
			return wait.TokenByInstanceNameWaitHandler(ctx, c, model.ProjectId, model.Region, model.identifier.Value, &expiration)
		}
		return factory, nil
	default:
		return nil, commonErr.NewInvalidIdentifierError(model.identifier.Flag)
	}
}

// Output result based on the configured output format
func outputResult(p *print.Printer, outputFormat string, token *edge.Token) error {
	if token == nil || token.Token == nil {
		// This is only to prevent nil pointer deref.
		// As long as the API behaves as defined by it's spec, instance can not be empty (HTTP 200 with an empty body)
		return fmt.Errorf("no token returned from the API")
	}
	tokenString := *token.Token

	return p.OutputResult(outputFormat, token, func() error {
		p.Outputln(tokenString)
		return nil
	})
}
