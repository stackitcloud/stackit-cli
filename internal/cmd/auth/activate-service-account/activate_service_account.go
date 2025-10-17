package activateserviceaccount

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
	sdkAuth "github.com/stackitcloud/stackit-sdk-go/core/auth"
	sdkConfig "github.com/stackitcloud/stackit-sdk-go/core/config"
)

const (
	serviceAccountTokenFlag   = "service-account-token"
	serviceAccountKeyPathFlag = "service-account-key-path"
	privateKeyPathFlag        = "private-key-path"
	onlyPrintAccessTokenFlag  = "only-print-access-token" // #nosec G101
)

type inputModel struct {
	ServiceAccountToken   string
	ServiceAccountKeyPath string
	PrivateKeyPath        string
	OnlyPrintAccessToken  bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activate-service-account",
		Short: "Authenticates using a service account",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Authenticates to the CLI using service account credentials.",
			"Subsequent commands will be authenticated using the service account credentials provided.",
			"For more details on how to configure your service account, check our Authentication guide at https://github.com/stackitcloud/stackit-cli/blob/main/AUTHENTICATION.md.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Activate service account authentication in the STACKIT CLI using a service account key which includes the private key`,
				"$ stackit auth activate-service-account --service-account-key-path path/to/service_account_key.json"),
			examples.NewExample(
				`Activate service account authentication in the STACKIT CLI using the service account key and explicitly providing the private key in a PEM encoded file, which will take precedence over the one in the service account key`,
				"$ stackit auth activate-service-account --service-account-key-path path/to/service_account_key.json --private-key-path path/to/private.key"),
			examples.NewExample(
				`Activate service account authentication in the STACKIT CLI using the service account token`,
				"$ stackit auth activate-service-account --service-account-token my-service-account-token"),
			examples.NewExample(
				`Only print the corresponding access token by using the service account token. This access token can be stored as environment variable (STACKIT_ACCESS_TOKEN) in order to be used for all subsequent commands.`,
				"$ stackit auth activate-service-account --service-account-token my-service-account-token --only-print-access-token",
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			tokenCustomEndpoint := viper.GetString(config.TokenCustomEndpointKey)
			if !model.OnlyPrintAccessToken {
				if err := storeCustomEndpoint(tokenCustomEndpoint); err != nil {
					return err
				}
			}

			cfg := &sdkConfig.Configuration{
				Token:                 model.ServiceAccountToken,
				ServiceAccountKeyPath: model.ServiceAccountKeyPath,
				PrivateKeyPath:        model.PrivateKeyPath,
				TokenCustomUrl:        tokenCustomEndpoint,
			}

			// Setup authentication based on the provided credentials and the environment
			// Initializes the authentication flow
			rt, err := sdkAuth.SetupAuth(cfg)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "setup auth: %v", err)
				return &cliErr.ActivateServiceAccountError{}
			}

			// Authenticates the service account and stores credentials
			email, accessToken, err := auth.AuthenticateServiceAccount(params.Printer, rt, model.OnlyPrintAccessToken)
			if err != nil {
				var activateServiceAccountError *cliErr.ActivateServiceAccountError
				if !errors.As(err, &activateServiceAccountError) {
					return fmt.Errorf("authenticate service account: %w", err)
				}
				return err
			}

			if model.OnlyPrintAccessToken {
				// Only output is the access token
				params.Printer.Outputf("%s\n", accessToken)
			} else {
				params.Printer.Outputf("You have been successfully authenticated to the STACKIT CLI!\nService account email: %s\n", email)
			}
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(serviceAccountTokenFlag, "", "Service account long-lived access token")
	cmd.Flags().String(serviceAccountKeyPathFlag, "", "Service account key path")
	cmd.Flags().String(privateKeyPathFlag, "", "RSA private key path. It takes precedence over the private key included in the service account key, if present")
	cmd.Flags().Bool(onlyPrintAccessTokenFlag, false, "If this is set to true the credentials are not stored in either the keyring or a file")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	model := inputModel{
		ServiceAccountToken:   flags.FlagToStringValue(p, cmd, serviceAccountTokenFlag),
		ServiceAccountKeyPath: flags.FlagToStringValue(p, cmd, serviceAccountKeyPathFlag),
		PrivateKeyPath:        flags.FlagToStringValue(p, cmd, privateKeyPathFlag),
		OnlyPrintAccessToken:  flags.FlagToBoolValue(p, cmd, onlyPrintAccessTokenFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func storeCustomEndpoint(tokenCustomEndpoint string) error {
	return auth.SetAuthField(auth.TOKEN_CUSTOM_ENDPOINT, tokenCustomEndpoint)
}
