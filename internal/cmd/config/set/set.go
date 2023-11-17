package set

import (
	"fmt"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	sessionTimeLimitFlag         = "session-time-limit"
	projectIdFlag                = "project-id"
	dnsCustomEndpointFlag        = "dns-custom-endpoint"
	postgreSQLCustomEndpointFlag = "postgresql-custom-endpoint"
)

type flagModel struct {
	SessionTimeLimit *string
}

var Cmd = &cobra.Command{
	Use:     "set",
	Short:   "Set CLI configuration options",
	Long:    "Set CLI configuration options",
	Example: `$ stackit config set --project-id xxx`,
	RunE: func(cmd *cobra.Command, args []string) error {
		model, err := parseFlags(cmd)
		if err != nil {
			return err
		}

		if model.SessionTimeLimit != nil {
			cmd.Println("Authenticate again to apply changes to session time limit")
			viper.Set(config.SessionTimeLimitKey, *model.SessionTimeLimit)
		}

		err = viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("write new config to file: %w", err)
		}
		return nil
	},
}

func init() {
	configureFlags(Cmd)
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(sessionTimeLimitFlag, "", "Maximum time before authentication is required again. Can't be larger than 24h. Examples: 3h, 5h30m40s")
	cmd.Flags().String(dnsCustomEndpointFlag, "", "DNS custom endpoint")
	cmd.Flags().String(postgreSQLCustomEndpointFlag, "", "PostgreSQL custom endpoint")

	err := viper.BindPFlag(config.DNSCustomEndpointKey, cmd.Flags().Lookup(dnsCustomEndpointFlag))
	cobra.CheckErr(err)
	err = viper.BindPFlag(config.PostgreSQLCustomEndpointKey, cmd.Flags().Lookup(postgreSQLCustomEndpointFlag))
	cobra.CheckErr(err)
}

func parseFlags(cmd *cobra.Command) (*flagModel, error) {
	sessionTimeLimit, err := parseSessionTimeLimit(cmd)
	if err != nil {
		return nil, fmt.Errorf("parse --%s: %w", sessionTimeLimitFlag, err)
	}

	return &flagModel{
		SessionTimeLimit: sessionTimeLimit,
	}, nil
}

func parseSessionTimeLimit(cmd *cobra.Command) (*string, error) {
	sessionTimeLimit := utils.FlagToStringPointer(cmd, sessionTimeLimitFlag)
	if sessionTimeLimit == nil {
		return nil, nil
	}

	// time.ParseDuration doesn't recognize unit "d", for simplicity we allow the value "1d"
	if *sessionTimeLimit == "1d" {
		*sessionTimeLimit = "24h"
	}

	duration, err := time.ParseDuration(*sessionTimeLimit)
	if err != nil {
		return nil, fmt.Errorf("parse value \"%s\": %w", *sessionTimeLimit, err)
	}
	if duration <= 0 {
		return nil, fmt.Errorf("value must be positive")
	}
	if duration > time.Duration(24)*time.Hour {
		return nil, fmt.Errorf("value can't be larger than 24h")
	}

	return sessionTimeLimit, nil
}
