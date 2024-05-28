package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/cmd/argus"
	"github.com/stackitcloud/stackit-cli/internal/cmd/auth"
	configCmd "github.com/stackitcloud/stackit-cli/internal/cmd/config"
	"github.com/stackitcloud/stackit-cli/internal/cmd/curl"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns"
	loadbalancer "github.com/stackitcloud/stackit-cli/internal/cmd/load-balancer"
	"github.com/stackitcloud/stackit-cli/internal/cmd/logme"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mariadb"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex"
	objectstorage "github.com/stackitcloud/stackit-cli/internal/cmd/object-storage"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project"
	"github.com/stackitcloud/stackit-cli/internal/cmd/rabbitmq"
	"github.com/stackitcloud/stackit-cli/internal/cmd/redis"
	secretsmanager "github.com/stackitcloud/stackit-cli/internal/cmd/secrets-manager"
	serviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/service-account"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(version, date string, p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "stackit",
		Short:             "Manage STACKIT resources using the command line",
		Long:              "Manage STACKIT resources using the command line.\nThis CLI is in a BETA state.\nMore services and functionality will be supported soon. Your feedback is appreciated!",
		Args:              args.NoArgs,
		SilenceErrors:     true, // Error is beautified in a custom way before being printed
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			p.Cmd = cmd
			p.Verbosity = print.Level(globalflags.Parse(p, cmd).Verbosity)

			argsString := print.BuildDebugStrFromSlice(os.Args)
			p.Debug(print.DebugLevel, "arguments: %s", argsString)

			configFilePath := viper.ConfigFileUsed()
			p.Debug(print.DebugLevel, "configuration is persisted and read from: %s", configFilePath)

			activeProfile, err := config.GetProfile()
			if err != nil {
				return fmt.Errorf("get profile: %w", err)
			}

			profileSet, err := config.GetConfiguredProfile()
			if err != nil {
				return fmt.Errorf("get profile raw: %w", err)
			}

			if activeProfile != profileSet {
				p.Debug(print.DebugLevel, "the active profile %q does not exist, following folder is missing: %q", profileSet, config.GetProfileFolderPath(profileSet))
				p.Warn("active profile %q does not exist, the %q profile configuration will be used\n", profileSet, activeProfile)
			}
			p.Debug(print.DebugLevel, "active configuration profile: %s", activeProfile)

			configKeys := viper.AllSettings()
			configKeysStr := print.BuildDebugStrFromMap(configKeys)
			p.Debug(print.DebugLevel, "configuration keys: %s", configKeysStr)

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.FlagToBoolValue(p, cmd, "version") {
				p.Outputf("STACKIT CLI (BETA)\n")

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

	addSubcommands(cmd, p)

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
	cobra.AddTemplateFunc("WhiteBold", color.New(color.FgHiWhite, color.Bold).SprintFunc())
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

func addSubcommands(cmd *cobra.Command, p *print.Printer) {
	cmd.AddCommand(argus.NewCmd(p))
	cmd.AddCommand(auth.NewCmd(p))
	cmd.AddCommand(configCmd.NewCmd(p))
	cmd.AddCommand(curl.NewCmd(p))
	cmd.AddCommand(dns.NewCmd(p))
	cmd.AddCommand(loadbalancer.NewCmd(p))
	cmd.AddCommand(logme.NewCmd(p))
	cmd.AddCommand(mariadb.NewCmd(p))
	cmd.AddCommand(mongodbflex.NewCmd(p))
	cmd.AddCommand(objectstorage.NewCmd(p))
	cmd.AddCommand(opensearch.NewCmd(p))
	cmd.AddCommand(organization.NewCmd(p))
	cmd.AddCommand(postgresflex.NewCmd(p))
	cmd.AddCommand(project.NewCmd(p))
	cmd.AddCommand(rabbitmq.NewCmd(p))
	cmd.AddCommand(redis.NewCmd(p))
	cmd.AddCommand(secretsmanager.NewCmd(p))
	cmd.AddCommand(serviceaccount.NewCmd(p))
	cmd.AddCommand(ske.NewCmd(p))
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
		p.Error(err.Error())
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
