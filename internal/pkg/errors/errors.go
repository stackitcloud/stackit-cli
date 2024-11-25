package errors

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const (
	MISSING_PROJECT_ID = `the project ID is not currently set.
	
It can be set on the command level by re-running your command with the --project-id flag.
	
You can configure it for all commands by running:
	
  $ stackit config set --project-id xxx
	
or you can also set it through the environment variable [STACKIT_PROJECT_ID]`

	EMPTY_UPDATE = `please specify at least one field to update.
	
Get details on the available flags by re-running your command with the --help flag.`

	REQUIRED_MUTUALLY_EXCLUSIVE_FLAGS = `the following flags are mutually exclusive and at least one of them is required: %s`

	FAILED_AUTH = `you are not authenticated.

You can authenticate as a user by running:
  $ stackit auth login

or use a service account by running:
  $ stackit auth activate-service-account`

	FAILED_SERVICE_ACCOUNT_ACTIVATION = `could not setup authentication based on the provided service account credentials. 
Please double check if they are correctly configured.

For more details run:
  $ stackit auth activate-service-account -h`

	SET_INEXISTENT_PROFILE = `the configuration profile %[1]q you are trying to set doesn't exist.

To create it, run:
  $ stackit config profile create %[1]q`

	DELETE_INEXISTENT_PROFILE = `the configuration profile %q does not exist.

To list all profiles, run:
  $ stackit config profile list`

	DELETE_DEFAULT_PROFILE = `the default configuration profile %q cannot be deleted.`

	OBSERVABILITY_INVALID_INPUT_PLAN = `the instance plan was not correctly provided. 

Either provide the plan ID:
  $ %[1]s --plan-id <PLAN ID> [flags]

or provide plan name:
  $ %[1]s --plan-name <PLAN NAME> [flags]

For more details on the available plans, run:
  $ stackit %[2]s plans`

	OBSERVABILITY_INVALID_PLAN = `the provided instance plan is not valid.
	
  %s
  
  For more details on the available plans, run:
	$ stackit %s plans`

	DSA_INVALID_INPUT_PLAN = `the instance plan was not correctly provided. 

Either provide the plan ID:
  $ %[1]s --plan-id <PLAN ID> [flags]

or provide plan name and version:
  $ %[1]s --plan-name <PLAN NAME> --version <VERSION> [flags]

For more details on the available plans, run:
  $ stackit %[2]s plans`

	DSA_INVALID_PLAN = `the provided instance plan is not valid.
	
%s

For more details on the available plans, run:
  $ stackit %s plans`

	DATABASE_INVALID_INPUT_FLAVOR = `the instance flavor was not correctly provided. 

Either provide flavor ID by:
  $ %[1]s --flavor-id <FLAVOR ID> [flags]

or provide CPU and RAM:
  $ %[1]s --cpu <CPU> --ram <RAM> [flags]

For more details on the available flavors, run:
  $ stackit %[2]s options --flavors`

	DATABASE_INVALID_FLAVOR = `the provided instance flavor is not valid.
	
%s

For more details on the available flavors, run:
  $ stackit %s options --flavors`

	DATABASE_INVALID_STORAGE = `invalid instance storage.
	
%[1]s

For more details on the available storages for the configured flavor (%[3]s), run:
  $ stackit %[2]s options --storages --flavor-id %[3]s`

	FLAG_VALIDATION = `the provided flag --%s is invalid: %s`

	ARG_VALIDATION = `the provided argument "%s" is invalid: %s`

	ARG_UNKNOWN = `unknown argument %q`

	ARG_MISSING = `missing argument %q`

	SINGLE_ARG_EXPECTED = `expected 1 argument %q, %d were provided`

	SINGLE_OPTIONAL_ARG_EXPECTED = `expected no more than 1 argument %q, %d were provided`

	SUBCOMMAND_UNKNOWN = `unknown subcommand %q`

	SUBCOMMAND_MISSING = `missing subcommand`

	INVALID_PROFILE_NAME = `the profile name %q is invalid.
	
The profile name can only contain lowercase letters, numbers, and "-" and cannot be empty or "default". It can't start with a "-".`

	USAGE_TIP = `For usage help, run:
  $ %s --help`

	SERVICE_DISABLED = `This service isn't enabled for the current project.

To enable it, run:
  $ stackit %s enable`

	IAAS_SERVER_MISSING_VOLUME_SIZE = `Boot volume size must be provided when "source_type" is "image".`

	IAAS_SERVER_MISSING_IMAGE_OR_VOLUME_FLAGS = `Either Image ID or boot volume flags must be provided.`
)

type ServerCreateMissingFlagsError struct {
	Cmd *cobra.Command
}

func (e *ServerCreateMissingFlagsError) Error() string {
	return IAAS_SERVER_MISSING_IMAGE_OR_VOLUME_FLAGS
}

type ServerCreateError struct {
	Cmd *cobra.Command
}

func (e *ServerCreateError) Error() string {
	return IAAS_SERVER_MISSING_VOLUME_SIZE
}

type ProjectIdError struct{}

func (e *ProjectIdError) Error() string {
	return MISSING_PROJECT_ID
}

type EmptyUpdateError struct{}

func (e *EmptyUpdateError) Error() string {
	return EMPTY_UPDATE
}

type AuthError struct{}

func (e *AuthError) Error() string {
	return FAILED_AUTH
}

type ActivateServiceAccountError struct{}

func (e *ActivateServiceAccountError) Error() string {
	return FAILED_SERVICE_ACCOUNT_ACTIVATION
}

type SetInexistentProfile struct {
	Profile string
}

func (e *SetInexistentProfile) Error() string {
	return fmt.Sprintf(SET_INEXISTENT_PROFILE, e.Profile)
}

type DeleteInexistentProfile struct {
	Profile string
}

func (e *DeleteInexistentProfile) Error() string {
	return fmt.Sprintf(DELETE_INEXISTENT_PROFILE, e.Profile)
}

type DeleteDefaultProfile struct {
	DefaultProfile string
}

func (e *DeleteDefaultProfile) Error() string {
	return fmt.Sprintf(DELETE_DEFAULT_PROFILE, e.DefaultProfile)
}

type ObservabilityInputPlanError struct {
	Cmd  *cobra.Command
	Args []string
}

func (e *ObservabilityInputPlanError) Error() string {
	fullCommandPath := e.Cmd.CommandPath()
	if len(e.Args) > 0 {
		fullCommandPath = fmt.Sprintf("%s %s", fullCommandPath, strings.Join(e.Args, " "))
	}
	// Assumes a structure of the form "stackit <service> <resource> <operation>"
	service := e.Cmd.Parent().Parent().Use

	return fmt.Sprintf(OBSERVABILITY_INVALID_INPUT_PLAN, fullCommandPath, service)
}

type ObservabilityInvalidPlanError struct {
	Service string
	Details string
}

func (e *ObservabilityInvalidPlanError) Error() string {
	return fmt.Sprintf(OBSERVABILITY_INVALID_PLAN, e.Details, e.Service)
}

type DSAInputPlanError struct {
	Cmd  *cobra.Command
	Args []string
}

func (e *DSAInputPlanError) Error() string {
	fullCommandPath := e.Cmd.CommandPath()
	if len(e.Args) > 0 {
		fullCommandPath = fmt.Sprintf("%s %s", fullCommandPath, strings.Join(e.Args, " "))
	}
	// Assumes a structure of the form "stackit <service> <resource> <operation>"
	service := e.Cmd.Parent().Parent().Use

	return fmt.Sprintf(DSA_INVALID_INPUT_PLAN, fullCommandPath, service)
}

type DSAInvalidPlanError struct {
	Service string
	Details string
}

func (e *DSAInvalidPlanError) Error() string {
	return fmt.Sprintf(DSA_INVALID_PLAN, e.Details, e.Service)
}

type DatabaseInputFlavorError struct {
	Service string
	Cmd     *cobra.Command
	Args    []string
}

func (e *DatabaseInputFlavorError) Error() string {
	fullCommandPath := e.Cmd.CommandPath()
	if len(e.Args) > 0 {
		fullCommandPath = fmt.Sprintf("%s %s", fullCommandPath, strings.Join(e.Args, " "))
	}

	if e.Service == "" {
		// Assumes a structure of the form "stackit <service> <resource> <operation>"
		e.Service = e.Cmd.Parent().Parent().Use
	}

	return fmt.Sprintf(DATABASE_INVALID_INPUT_FLAVOR, fullCommandPath, e.Service)
}

type DatabaseInvalidFlavorError struct {
	Service string
	Details string
}

func (e *DatabaseInvalidFlavorError) Error() string {
	return fmt.Sprintf(DATABASE_INVALID_FLAVOR, e.Details, e.Service)
}

type DatabaseInvalidStorageError struct {
	Service  string
	Details  string
	FlavorId string
}

func (e *DatabaseInvalidStorageError) Error() string {
	return fmt.Sprintf(DATABASE_INVALID_STORAGE, e.Details, e.Service, e.FlavorId)
}

type FlagValidationError struct {
	Flag    string
	Details string
}

func (e *FlagValidationError) Error() string {
	return fmt.Sprintf(FLAG_VALIDATION, e.Flag, e.Details)
}

type RequiredMutuallyExclusiveFlagsError struct {
	Flags []string
}

func (e *RequiredMutuallyExclusiveFlagsError) Error() string {
	return fmt.Sprintf(REQUIRED_MUTUALLY_EXCLUSIVE_FLAGS, strings.Join(e.Flags, ", "))
}

type ArgValidationError struct {
	Arg     string
	Details string
}

func (e *ArgValidationError) Error() string {
	return fmt.Sprintf(ARG_VALIDATION, e.Arg, e.Details)
}

type SingleArgExpectedError struct {
	Cmd      *cobra.Command
	Expected string
	Count    int
}

func (e *SingleArgExpectedError) Error() string {
	var err error
	if e.Count > 1 {
		err = fmt.Errorf(SINGLE_ARG_EXPECTED, e.Expected, e.Count)
	} else {
		err = fmt.Errorf(ARG_MISSING, e.Expected)
	}
	return AppendUsageTip(err, e.Cmd).Error()
}

type SingleOptionalArgExpectedError struct {
	Cmd      *cobra.Command
	Expected string
	Count    int
}

func (e *SingleOptionalArgExpectedError) Error() string {
	err := fmt.Errorf(SINGLE_OPTIONAL_ARG_EXPECTED, e.Expected, e.Count)
	return AppendUsageTip(err, e.Cmd).Error()
}

// Used when an unexpected non-flag input (either arg or subcommand) is found
type InputUnknownError struct {
	ProvidedInput string
	Cmd           *cobra.Command
}

func (e *InputUnknownError) Error() string {
	// To decide whether the unexpected input is an arg or a subcommand, we assume that only leaf commands (ie, don't have subcomamnds) take args
	var err error
	if !e.Cmd.HasSubCommands() {
		err = fmt.Errorf(ARG_UNKNOWN, e.ProvidedInput)
	} else {
		err = fmt.Errorf(SUBCOMMAND_UNKNOWN, e.ProvidedInput)
	}
	return AppendUsageTip(err, e.Cmd).Error()
}

type SubcommandMissingError struct {
	Cmd *cobra.Command
}

func (e *SubcommandMissingError) Error() string {
	err := fmt.Errorf(SUBCOMMAND_MISSING)
	return AppendUsageTip(err, e.Cmd).Error()
}

// Returns a wrapped error whose message adds a tip on how to check out --help for the command
func AppendUsageTip(err error, cmd *cobra.Command) error {
	tip := fmt.Sprintf(USAGE_TIP, cmd.CommandPath())
	return fmt.Errorf("%w.\n\n%s", err, tip)
}

type InvalidProfileNameError struct {
	Profile string
}

func (e *InvalidProfileNameError) Error() string {
	return fmt.Sprintf(INVALID_PROFILE_NAME, e.Profile)
}

type ServiceDisabledError struct {
	Service string
}

func (e *ServiceDisabledError) Error() string {
	return fmt.Sprintf(SERVICE_DISABLED, e.Service)
}
