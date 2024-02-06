package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/cmd/auth"
	"github.com/stackitcloud/stackit-cli/internal/cmd/config"
	"github.com/stackitcloud/stackit-cli/internal/cmd/curl"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns"
	"github.com/stackitcloud/stackit-cli/internal/cmd/mongodbflex"
	"github.com/stackitcloud/stackit-cli/internal/cmd/opensearch"
	"github.com/stackitcloud/stackit-cli/internal/cmd/organization"
	"github.com/stackitcloud/stackit-cli/internal/cmd/postgresflex"
	"github.com/stackitcloud/stackit-cli/internal/cmd/project"
	serviceaccount "github.com/stackitcloud/stackit-cli/internal/cmd/service-account"
	"github.com/stackitcloud/stackit-cli/internal/cmd/ske"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"

	"github.com/spf13/cobra"
)

func NewRootCmd(version, date string) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "stackit",
		Short:             "Manage STACKIT resources using the command line",
		Long:              "Manage STACKIT resources using the command line.\nThis CLI is in a BETA state.\nMore services and functionality will be supported soon. Your feedback is appreciated!",
		Args:              args.NoArgs,
		SilenceErrors:     true, // Error is beautified in a custom way before being printed
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.FlagToBoolValue(cmd, "version") {
				cmd.Printf("STACKIT CLI (BETA)\n")

				parsedDate, err := time.Parse(time.RFC3339, date)
				if err != nil {
					cmd.Printf("Version: %s\n", version)
					return nil
				}
				cmd.Printf("Version: %s (%s)\n", version, parsedDate.Format(time.DateOnly))
				return nil
			}

			return cmd.Help()
		},
	}
	cmd.SetOut(os.Stdout)

	err := configureFlags(cmd)
	cobra.CheckErr(err)

	addSubcommands(cmd)

	// Cobra creates the help flag with "help for <command>" as the description
	// We want to override that message by capitalizing the first letter to match the other flag descriptions
	// See spf13/cobra#480
	traverseCommands(cmd, func(c *cobra.Command) {
		c.Flags().BoolP("help", "h", false, fmt.Sprintf("Help for %q", c.CommandPath()))
	})

	return cmd
}

func configureFlags(cmd *cobra.Command) error {
	cmd.Flags().BoolP("version", "v", false, `Show "stackit" version`)

	err := globalflags.Configure(cmd.PersistentFlags())
	if err != nil {
		return fmt.Errorf("configure global flags: %w", err)
	}
	return nil
}

func addSubcommands(cmd *cobra.Command) {
	cmd.AddCommand(auth.NewCmd())
	cmd.AddCommand(config.NewCmd())
	cmd.AddCommand(curl.NewCmd())
	cmd.AddCommand(dns.NewCmd())
	cmd.AddCommand(mongodbflex.NewCmd())
	cmd.AddCommand(opensearch.NewCmd())
	cmd.AddCommand(organization.NewCmd())
	cmd.AddCommand(postgresflex.NewCmd())
	cmd.AddCommand(project.NewCmd())
	cmd.AddCommand(serviceaccount.NewCmd())
	cmd.AddCommand(ske.NewCmd())
}

// traverseCommands calls f for c and all of its children.
func traverseCommands(c *cobra.Command, f func(*cobra.Command)) {
	f(c)
	for _, c := range c.Commands() {
		traverseCommands(c, f)
	}
}

func Execute(version, date string) {
	cmd := NewRootCmd(version, date)
	err := cmd.Execute()
	if err != nil {
		err := beautifyUnknownAndMissingCommandsError(cmd, err)
		cmd.PrintErrln(cmd.ErrPrefix(), err.Error())
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
