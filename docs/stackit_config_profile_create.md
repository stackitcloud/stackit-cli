## stackit config profile create

Creates a CLI configuration profile

### Synopsis

Creates a CLI configuration profile based on the currently active profile and sets it as active.
The profile name can be provided via the STACKIT_CLI_PROFILE environment variable or as an argument in this command.
The environment variable takes precedence over the argument.
If you do not want to set the profile as active, use the --no-set flag.
If you want to create the new profile with the initial default configurations, use the --empty flag.

```
stackit config profile create PROFILE [flags]
```

### Examples

```
  Create a new configuration profile "my-profile" with the current configuration, setting it as the active profile
  $ stackit config profile create my-profile

  Create a new configuration profile "my-profile" with a default initial configuration and don't set it as the active profile
  $ stackit config profile create my-profile --empty --no-set
```

### Options

```
      --empty    Create the profile with the initial default configurations
  -h, --help     Help for "stackit config profile create"
      --no-set   Do not set the profile as the active profile
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --region string          Target region for region-specific requests
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit config profile](./stackit_config_profile.md)	 - Manage the CLI configuration profiles

