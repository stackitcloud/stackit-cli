package activateserviceaccount

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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

type flagModel struct {
	ServiceAccountToken   string
	ServiceAccountKeyPath string
	PrivateKeyPath        string
	TokenCustomEndpoint   string
	JwksCustomEndpoint    string
}

var Cmd = &cobra.Command{
	Use:     "activate-service-account",
	Short:   "Activate service account authentication",
	Long:    "Activate authentication using service account credentials.\nFor more details on how to configure your service account, check the Authentication section on our documentation (LINK HERE README)",
	Example: `$ stackit auth activate-service-account --service-account-key-path path/to/service_account_key.json --private-key-path path/to/private_key.pem`,
	RunE: func(cmd *cobra.Command, args []string) error {
		model := parseFlags(cmd)

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
			return fmt.Errorf("set up authentication: %w", err)
		}

		// Authenticates the service account and stores credentials
		email, err := auth.AuthenticateServiceAccount(rt)
		if err != nil {
			return fmt.Errorf("authenticate service account: %w", err)
		}

		cmd.Printf("You have been successfully authenticated to the STACKIT CLI!\nService account email: %s\n", email)

		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(serviceAccountTokenFlag, "", "Service account long-lived access token")
	cmd.Flags().String(serviceAccountKeyPathFlag, "", "Service account key path")
	cmd.Flags().String(privateKeyPathFlag, "", "RSA private key path")
	cmd.Flags().String(tokenCustomEndpointFlag, "", "Custom endpoint for the token API, which is used to request access tokens when the service-account authentication is activated")
	cmd.Flags().String(jwksCustomEndpointFlag, "", "Custom endpoint for the jwks API, which is used to get the json web key sets (jwks) to validate tokens when the service-account authentication is activated")
}

func parseFlags(cmd *cobra.Command) *flagModel {
	return &flagModel{
		ServiceAccountToken:   utils.FlagToStringValue(cmd, serviceAccountTokenFlag),
		ServiceAccountKeyPath: utils.FlagToStringValue(cmd, serviceAccountKeyPathFlag),
		PrivateKeyPath:        utils.FlagToStringValue(cmd, privateKeyPathFlag),
		TokenCustomEndpoint:   utils.FlagToStringValue(cmd, tokenCustomEndpointFlag),
		JwksCustomEndpoint:    utils.FlagToStringValue(cmd, jwksCustomEndpointFlag),
	}
}

func storeFlags(model *flagModel) error {
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
