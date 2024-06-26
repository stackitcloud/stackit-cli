package activateserviceaccount

import (
	"errors"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
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
	tokenCustomEndpointFlag   = "token-custom-endpoint"
	jwksCustomEndpointFlag    = "jwks-custom-endpoint"
)

type inputModel struct {
	ServiceAccountToken   string
	ServiceAccountKeyPath string
	PrivateKeyPath        string
	TokenCustomEndpoint   string
	JwksCustomEndpoint    string
}

func NewCmd(p *print.Printer) *cobra.Command {
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
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			model := parseInput(p, cmd)

			err := storeFlags(model)
			if err != nil {
				return err
			}

			cfg := &sdkConfig.Configuration{
				Token:                 model.ServiceAccountToken,
				ServiceAccountKeyPath: model.ServiceAccountKeyPath,
				PrivateKeyPath:        model.PrivateKeyPath,
				TokenCustomUrl:        model.TokenCustomEndpoint,
				JWKSCustomUrl:         model.JwksCustomEndpoint,
			}

			// Setup authentication based on the provided credentials and the environment
			// Initializes the authentication flow
			rt, err := sdkAuth.SetupAuth(cfg)
			if err != nil {
				p.Debug(print.ErrorLevel, "setup auth: %v", err)
				return &cliErr.ActivateServiceAccountError{}
			}

			// Authenticates the service account and stores credentials
			email, err := auth.AuthenticateServiceAccount(p, rt)
			if err != nil {
				var activateServiceAccountError *cliErr.ActivateServiceAccountError
				if !errors.As(err, &activateServiceAccountError) {
					return fmt.Errorf("authenticate service account: %w", err)
				}
				return err
			}

			p.Info("You have been successfully authenticated to the STACKIT CLI!\nService account email: %s\n", email)

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
	cmd.Flags().String(tokenCustomEndpointFlag, "", "Custom endpoint for the token API, which is used to request access tokens when the service-account authentication is activated")
	cmd.Flags().String(jwksCustomEndpointFlag, "", "Custom endpoint for the jwks API, which is used to get the json web key sets (jwks) to validate tokens when the service-account authentication is activated")
}

func parseInput(p *print.Printer, cmd *cobra.Command) *inputModel {
	model := inputModel{
		ServiceAccountToken:   flags.FlagToStringValue(p, cmd, serviceAccountTokenFlag),
		ServiceAccountKeyPath: flags.FlagToStringValue(p, cmd, serviceAccountKeyPathFlag),
		PrivateKeyPath:        flags.FlagToStringValue(p, cmd, privateKeyPathFlag),
		TokenCustomEndpoint:   flags.FlagToStringValue(p, cmd, tokenCustomEndpointFlag),
		JwksCustomEndpoint:    flags.FlagToStringValue(p, cmd, jwksCustomEndpointFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model
}

func storeFlags(model *inputModel) error {
	err := auth.SetAuthField(auth.TOKEN_CUSTOM_ENDPOINT, model.TokenCustomEndpoint)
	if err != nil {
		return fmt.Errorf("set %s: %w", auth.TOKEN_CUSTOM_ENDPOINT, err)
	}
	err = auth.SetAuthField(auth.JWKS_CUSTOM_ENDPOINT, model.JwksCustomEndpoint)
	if err != nil {
		return fmt.Errorf("set %s: %w", auth.JWKS_CUSTOM_ENDPOINT, err)
	}
	return nil
}
