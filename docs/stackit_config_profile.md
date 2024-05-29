## stackit config profile

Manage the CLI configuration profiles

### Synopsis

Manage the CLI configuration profiles.
The profile to be used can be managed via the "STACKIT_CLI_PROFILE" environment variable or using the "stackit config profile set PROFILE" and "stackit config profile unset" commands.
The environment variable takes precedence over what is set via the commands.
When no profile is set, the default profile is used.

```
stackit config profile [flags]
```

### Options

```
  -h, --help   Help for "stackit config profile"
```

### Options inherited from parent commands

```
  -y, --assume-yes             If set, skips all confirmation prompts
      --async                  If set, runs the command asynchronously
  -o, --output-format string   Output format, one of ["json" "pretty" "none" "yaml"]
  -p, --project-id string      Project ID
      --verbosity string       Verbosity of the CLI, one of ["debug" "info" "warning" "error"] (default "info")
```

### SEE ALSO

* [stackit config](./stackit_config.md)	 - Provides functionality for CLI configuration options
* [stackit config profile create](./stackit_config_profile_create.md)	 - Creates a CLI configuration profile
* [stackit config profile list](./stackit_config_profile_list.md)	 - Lists all CLI configuration profiles
* [stackit config profile set](./stackit_config_profile_set.md)	 - Set a CLI configuration profile
* [stackit config profile set](./stackit_config_profile_set.md)	 - Delete a CLI configuration profile
* [stackit config profile unset](./stackit_config_profile_unset.md)	 - Unset the current active CLI configuration profile

