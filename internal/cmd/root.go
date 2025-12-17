package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	affinityGroups "github.com/stackitcloud/stackit-cli/internal/cmd/affinity-groups"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta"
	configCmd "github.com/stackitcloud/stackit-cli/internal/cmd/config"
	"github.com/stackitcloud/stackit-cli/internal/cmd/curl"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns"
	"github.com/stackitcloud/stackit-cli/internal/cmd/git"
	"github.com/stackitcloud/stackit-cli/internal/cmd/image"
	keypair "github.com/stackitcloud/stackit-cli/internal/cmd/key-pair"
	loadbalancer "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex"
	"github.com/stackitcloud/stackit-cli/internal/cmd/network"
	networkArea "github.com/stackitcloud/stackit-cli/internal/cmd/network-area"
	networkinterface "github.com/stackitcloud/stackit-cli/internal/cmd/network-interface"
	objectstorage "github.com/stackitcloud/stackit-cli/internal/cmd/object-storage"
	"github.com/stackitcloud/stackit-cli/internal/cmd/observability"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project"
	publicip "github.com/stackitcloud/stackit-cli/internal/cmd/public-ip"
	"github.com/stackitcloud/stackit-cli/internal/cmd/quota"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis"
	secretsmanager "github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager"
	securitygroup "github.com/stackitcloud/stackit-cli/internal/cmd/security-group"
	"github.com/stackitcloud/stackit-cli/internal/cmd/server"
	serviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/service-account"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske"
	"github.com/stackitcloud/stackit-cli/internal/cmd/volume"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(version, date string, p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "stackit",
		Short:             "Manage STACKIT resources using the command line",
		Long:              "Manage STACKIT resources using the command line.\nYour feedback is appreciated!",
		Args:              args.NoArgs,
		SilenceErrors:     true, // Error is beautified in a custom way before being printed
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			p.Cmd = cmd
			p.Verbosity = print.Level(globalflags.Parse(p, cmd).Verbosity)

			argsString := print.BuildDebugStrFromSlice(os.Args)
			p.Debug(print.DebugLevel, "arguments: %s", argsString)

			configFilePath := viper.ConfigFileUsed()
			p.Debug(print.DebugLevel, "configuration is persisted and read from: %s", configFilePath)

			profileSet, activeProfile, configMethod, err := config.GetConfiguredProfile()
			if err != nil {
				return fmt.Errorf("get configured profile: %w", err)
			}

			p.Debug(print.DebugLevel, "read configuration profile %q via %s", profileSet, configMethod)

			if activeProfile != profileSet {
				if configMethod == "" {
					p.Debug(print.DebugLevel, "no profile is configured in env var or profile file")
				} else {
					p.Debug(print.DebugLevel, "the configured profile %q does not exist: folder %q is missing", profileSet, config.GetProfileFolderPath(profileSet))
				}
				p.Debug(print.DebugLevel, "the %q profile will be used", activeProfile)

				p.Warn("configured profile %q does not exist, the %q profile configuration will be used\n", profileSet, activeProfile)
			}

			p.Debug(print.DebugLevel, "active configuration profile: %s", activeProfile)

			configKeys := viper.AllSettings()
			configKeysStr := print.BuildDebugStrFromMap(configKeys)
			p.Debug(print.DebugLevel, "configuration keys: %s", configKeysStr)

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			if flags.FlagToBoolValue(p, cmd, "version") {
				p.Outputf("STACKIT CLI\n")

				parsedDate, err := time.Parse(time.RFC3339, date)
				if err != nil {
					p.Outputf("Version: %s\n", version)
					return nil
				}
				p.Outputf("Version: %s (%s)\n", version, parsedDate.Format(time.DateOnly))
				return nil
			}

			return cmd.Help()
		},
	}
	cmd.SetOut(os.Stdout)

	err := configureFlags(cmd)
	cobra.CheckErr(err)

	addSubcommands(cmd, &types.CmdParams{
		Printer:    p,
		CliVersion: version,
	})

	// Cobra creates the help flag with "help for <command>" as the description
	// We want to override that message by capitalizing the first letter to match the other flag descriptions
	// See spf13/cobra#480
	traverseCommands(cmd, func(c *cobra.Command) {
		c.Flags().BoolP("help", "h", false, fmt.Sprintf("Help for %q", c.CommandPath()))
	})

	beautifyUsageTemplate(cmd)

	return cmd
}

func beautifyUsageTemplate(cmd *cobra.Command) {
	cobra.AddTemplateFunc("WhiteBold", print.WhiteBold)
	usageTemplate := cmd.UsageTemplate()
	usageTemplate = strings.NewReplacer(
		`Usage:`, `{{WhiteBold "USAGE"}}`,
		`Examples:`, `{{WhiteBold "EXAMPLES"}}`,
		`Aliases:`, `{{WhiteBold "ALIASES"}}`,
		`Available Commands:`, `{{WhiteBold "AVAILABLE COMMANDS"}}`,
		`Additional Commands:`, `{{WhiteBold "ADDITIONAL COMMANDS"}}`,
		`Global Flags:`, `{{WhiteBold "GLOBAL FLAGS"}}`,
		`Flags:`, `{{WhiteBold "FLAGS"}}`,
		`Additional help topics:`, `{{WhiteBold "ADDITIONAL HELP TOPICS"}}`,
	).Replace(usageTemplate)
	cmd.SetUsageTemplate(usageTemplate)
}

func configureFlags(cmd *cobra.Command) error {
	cmd.Flags().BoolP("version", "v", false, `Show "stackit" version`)

	err := globalflags.Configure(cmd.PersistentFlags())
	if err != nil {
		return fmt.Errorf("configure global flags: %w", err)
	}
	return nil
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(auth.NewCmd(params))
	cmd.AddCommand(configCmd.NewCmd(params))
	cmd.AddCommand(beta.NewCmd(params))
	cmd.AddCommand(curl.NewCmd(params))
	cmd.AddCommand(dns.NewCmd(params))
	cmd.AddCommand(loadbalancer.NewCmd(params))
	cmd.AddCommand(logme.NewCmd(params))
	cmd.AddCommand(mariadb.NewCmd(params))
	cmd.AddCommand(mongodbflex.NewCmd(params))
	cmd.AddCommand(objectstorage.NewCmd(params))
	cmd.AddCommand(observability.NewCmd(params))
	cmd.AddCommand(opensearch.NewCmd(params))
	cmd.AddCommand(organization.NewCmd(params))
	cmd.AddCommand(postgresflex.NewCmd(params))
	cmd.AddCommand(project.NewCmd(params))
	cmd.AddCommand(rabbitmq.NewCmd(params))
	cmd.AddCommand(redis.NewCmd(params))
	cmd.AddCommand(secretsmanager.NewCmd(params))
	cmd.AddCommand(serviceaccount.NewCmd(params))
	cmd.AddCommand(ske.NewCmd(params))
	cmd.AddCommand(server.NewCmd(params))
	cmd.AddCommand(networkArea.NewCmd(params))
	cmd.AddCommand(network.NewCmd(params))
	cmd.AddCommand(volume.NewCmd(params))
	cmd.AddCommand(networkinterface.NewCmd(params))
	cmd.AddCommand(publicip.NewCmd(params))
	cmd.AddCommand(securitygroup.NewCmd(params))
	cmd.AddCommand(keypair.NewCmd(params))
	cmd.AddCommand(image.NewCmd(params))
	cmd.AddCommand(quota.NewCmd(params))
	cmd.AddCommand(affinityGroups.NewCmd(params))
	cmd.AddCommand(git.NewCmd(params))
}

// traverseCommands calls f for c and all of its children.
func traverseCommands(c *cobra.Command, f func(*cobra.Command)) {
	f(c)
	for _, c := range c.Commands() {
		traverseCommands(c, f)
	}
}

func Execute(version, date string) {
	p := print.NewPrinter()
	cmd := NewRootCmd(version, date, p)

	// We need to set the printer and verbosity here because the
	// PersistentPreRun is not called when the command is wrongly called
	p.Cmd = cmd
	p.Verbosity = print.InfoLevel

	err := cmd.Execute()
	if err != nil {
		err := beautifyUnknownAndMissingCommandsError(cmd, err)
		p.Debug(print.ErrorLevel, "execute command: %v", err)
		p.Error("%s", err.Error())
		os.Exit(1)
	}
}

// Returns a more user-friendly error if the input error is due to unknown/missing subcommands (issue: https://github.com/spf13/cobra/issues/706)
//
// Otherwise, returns the input error unchanged
func beautifyUnknownAndMissingCommandsError(rootCmd *cobra.Command, cmdErr error) error {
	if !strings.HasPrefix(cmdErr.Error(), "unknown flag") {
		return cmdErr
	}

	cmd, unparsedInputs, err := rootCmd.Traverse(os.Args[1:])
	if err != nil {
		return cmdErr
	}
	if len(unparsedInputs) == 0 {
		// This shouldn't happen
		// If we're here, Cobra was able to parse everything, thus it wouldn't raise "unknown flag" errors
		return cmdErr
	}

	// If cmd itself has more subcommands, we assume it has no logic by itself (other than --help)
	// We want the error message to state that either a cmd's subcommand is missing, or that the cmd's subcommand called is wrong
	if cmd.HasSubCommands() {
		if strings.HasPrefix(unparsedInputs[0], "-") {
			return &errors.SubcommandMissingError{
				Cmd: cmd,
			}
		}

		return &errors.InputUnknownError{
			ProvidedInput: unparsedInputs[0],
			Cmd:           cmd,
		}
	}

	// If we're here, cmd doesn't have subcommands command
	// If Cobra raised "unknown flag" errors, then it was while parsing cmd's flags
	// To be more user-friendly, we add a usage tip
	err = cmd.ParseFlags(unparsedInputs)
	if err != nil {
		return errors.AppendUsageTip(err, cmd)
	}

	// This shouldn't happen
	// If we're here, Cobra was able to parse cmd's flags, thus it wouldn't raise "unknown flag" errors
	return &errors.InputUnknownError{
		ProvidedInput: unparsedInputs[0],
		Cmd:           cmd,
	}
}
